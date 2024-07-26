package controller

import (
	"log"
	"net/http"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/middleware"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/request"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/token"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/service"
	"github.com/gin-gonic/gin"
)

// AccountController handles HTTP requests related to accounts.
type AccountController struct {
	accountService *service.AccountService
}

// NewAccountController creates a new instance of the AccountController struct.
// It takes an accountService parameter of type *service.AccountService and returns a pointer to the AccountController.
func NewAccountController(accountService *service.AccountService) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}

// CreateAccount creates a new account based on the provided JSON request.
// It binds the JSON request to the CreateAccountRequest struct and performs validation.
// If the request is invalid, it returns a JSON error response with the validation errors.
// If the request is valid, it creates a new account using the accountService and returns a JSON response with the created account.
// If there is an error during the account creation process, it returns a JSON error response with the appropriate error message.
// If the account already exists, it returns a JSON error response indicating that an account with the same currency already exists.
func (a *AccountController) CreateAccount(ctx *gin.Context) {
	var request request.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		message, data := helper.GlobalCheckingErrorBindJson(err.Error(), request)
		log.Printf("Error Bind: %s", message)
		helper.ReturnJSONError(ctx, http.StatusBadRequest, message, nil, data)
		return
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)
	res := helper.DoValidation(&request)

	if len(res) > 0 {
		log.Println("Error: Validation error")
		helper.ReturnJSONError(ctx, http.StatusBadRequest, "Validation error", nil, res)
		return
	}

	account, err := a.accountService.CreateAccount(ctx.Request.Context(), &request, authPayload)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		if err.Error() == "sql: no rows in result set" {
			helper.ReturnJSONError(ctx, http.StatusNotFound, "Data Not found", nil, nil)
		} else {
			helper.ReturnJSONError(ctx, http.StatusInternalServerError, "Internal server error", nil, err.Error())
		}
		return
	}

	// checking account already exist
	if account.Owner == "" {
		log.Println("Error: Account already exists")
		helper.ReturnJSONError(ctx, http.StatusConflict, "Account with same currency already exists", nil, nil)
		return
	}
	helper.ReturnJSON(ctx, http.StatusCreated, "Account created", account)
}

// GetAccount retrieves an account based on the provided UUID from the request context.
// It handles the binding of the URI parameters, converts the UUID to the appropriate format,
// and performs authorization checks before retrieving the account.
// If the account is found, it returns a JSON response with the account details.
// If any errors occur during the process, appropriate error responses are returned.
func (a *AccountController) GetAccount(ctx *gin.Context) {
	var req request.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		message, data := helper.GlobalCheckingErrorBindJson(err.Error(), req)
		log.Printf("Error: %s", message)
		helper.ReturnJSONError(ctx, http.StatusBadRequest, message, nil, data)
		return
	}

	// covert uuid to uuid.UUID
	uuidAcc, err := helper.ConvertStringToUUID(req.UUIDAcc)
	if err != nil {
		log.Println("Error: Invalid UUID")
		helper.ReturnJSONError(ctx, http.StatusBadRequest, "Invalid Parameter", nil, nil)
		return
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

	account, err := a.accountService.GetAccountByUUID(ctx.Request.Context(), uuidAcc, authPayload)
	if err != nil {
		if err.Error() == "unauthorized" {
			log.Println("Error: Unauthorized")
			helper.ReturnJSONError(ctx, http.StatusUnauthorized, "Unauthorized", nil, nil)
			return
		} else if err.Error() == "sql: no rows in result set" {
			log.Println("Error: Data not found")
			helper.ReturnJSONError(ctx, http.StatusNotFound, "Data not found", nil, nil)
			return
		}
		log.Printf("Error: %s", err.Error())
		helper.ReturnJSONError(ctx, http.StatusInternalServerError, "Internal server error", nil, nil)
		return
	}

	if account.Owner == "" {
		log.Println("Error: Account not found")
		helper.ReturnJSONError(ctx, http.StatusNotFound, "Account not found", nil, nil)
		return
	}

	helper.ReturnJSON(ctx, http.StatusOK, "Account found", account)
}

// GetAccounts retrieves a list of accounts based on the provided query parameters.
// It binds the query parameters from the request context, validates them, and then
// calls the accountService to fetch the accounts from the database. The function
// returns the list of accounts along with the total number of data, or an error
// if there was a problem with the request or the database.
//
// Parameters:
// - ctx: The gin.Context object representing the HTTP request and response.
//
// Returns:
// - None
//
// Example usage:
//
//	router.GET("/accounts", accountController.GetAccounts)
func (a *AccountController) GetAccounts(ctx *gin.Context) {
	var req request.ListAccountRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		message, data := helper.GlobalCheckingErrorBindJson(err.Error(), req)
		log.Printf("Error: %s", message)
		helper.ReturnJSONError(ctx, http.StatusBadRequest, message, nil, data)
		return
	}

	param := db.ListAccountsParams{
		Limit:  req.Limit,
		Offset: (req.Page - 1) * req.Limit,
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

	accounts, totalData, err := a.accountService.ListAccount(ctx.Request.Context(), param, authPayload)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		helper.ReturnJSONError(ctx, http.StatusInternalServerError, "Internal server error", nil, nil)
		return
	}

	if len(accounts) == 0 {
		log.Println("Error: Account not found")
		helper.ReturnJSONError(ctx, http.StatusNotFound, "Account not found", nil, nil)
		return
	}

	helper.ReturnJSONWithMetaPage(ctx, http.StatusOK, "Account found", accounts, int(totalData), len(accounts), int(req.Page), int(req.Limit))
}

// UpdateAccount updates the account information based on the provided request.
// It first binds the URI parameters and checks for any binding errors.
// If there are binding errors, it returns a JSON error response with the error details.
// Then, it converts the UUID string to a UUID type and checks for any conversion errors.
// If there are conversion errors, it returns a JSON error response with the error details.
// Next, it binds the JSON request body and checks for any binding errors.
// If there are binding errors, it returns a JSON error response with the error details.
// After that, it performs validation on the request data and checks for any validation errors.
// If there are validation errors, it returns a JSON error response with the error details.
// Finally, it updates the account in the database and returns a JSON response with the updated account information.
func (a *AccountController) UpdateAccount(ctx *gin.Context) {
	var req request.GetAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		message, data := helper.GlobalCheckingErrorBindJson(err.Error(), req)
		log.Printf("Error: %s", message)
		helper.ReturnJSONError(ctx, http.StatusBadRequest, message, nil, data)
		return
	}

	// covert uuid to uuid.UUID
	uuidAcc, err := helper.ConvertStringToUUID(req.UUIDAcc)
	if err != nil {
		log.Println("Error: Invalid UUID")
		helper.ReturnJSONError(ctx, http.StatusBadRequest, "Invalid Parameter", nil, nil)
		return
	}

	var request request.UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		message, data := helper.GlobalCheckingErrorBindJson(err.Error(), req)
		log.Printf("Error: %s", message)
		helper.ReturnJSONError(ctx, http.StatusBadRequest, message, nil, data)
		return
	}

	res := helper.DoValidation(&request)

	if len(res) > 0 {
		log.Println("Error: Validation error")
		helper.ReturnJSONError(ctx, http.StatusBadRequest, "Validation error", nil, res)
		return
	}

	arg := db.UpdateProfileAccountParams{
		AccountUuid: uuidAcc,
		Owner:       request.Owner,
		Currency:    request.Currency,
		Status:      int16(request.Status),
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)
	account, err := a.accountService.UpdateAccount(ctx.Request.Context(), arg, authPayload)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Println("Error: Data not found")
			helper.ReturnJSONError(ctx, http.StatusNotFound, "Data not found", nil, nil)
			return
		}

		if err.Error() == "unauthorized" {
			log.Println("Error: Unauthorized")
			helper.ReturnJSONError(ctx, http.StatusUnauthorized, "Unauthorized", nil, nil)
			return
		}

		log.Printf("Error: %s", err.Error())
		helper.ReturnJSONError(ctx, http.StatusInternalServerError, "Internal server error", nil, nil)
		return
	}

	if account.AccountUUID.String() == "" {
		log.Println("Error: Account not found")
		helper.ReturnJSONError(ctx, http.StatusNotFound, "Account not found", nil, nil)
		return
	}

	helper.ReturnJSON(ctx, http.StatusOK, "Account updated", account)
}
