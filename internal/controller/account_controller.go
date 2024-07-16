package controller

import (
	"log"
	"net/http"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/request"
	"github.com/fajaramaulana/simple_bank_project/internal/service"
	"github.com/gin-gonic/gin"
)

type AccountController struct {
	accountService *service.AccountService
}

func NewAccountController(accountService *service.AccountService) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}

func (a *AccountController) CreateAccount(ctx *gin.Context) {
	var request request.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		message, data := helper.GlobalCheckingErrorBindJson(err.Error())
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

	account, err := a.accountService.CreateAccount(ctx.Request.Context(), &request)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		helper.ReturnJSONError(ctx, http.StatusInternalServerError, "Internal server error", nil, nil)
		return
	}

	// checking account already exist
	if account.Email == "" {
		log.Println("Error: Account already exists")
		helper.ReturnJSONError(ctx, http.StatusConflict, "Account already exists", nil, nil)
		return
	}
	helper.ReturnJSON(ctx, http.StatusCreated, "Account created", account)
}

func (a *AccountController) GetAccount(ctx *gin.Context) {
	// uuid := ctx.Param("uuid")
	var req request.GetAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		message, data := helper.GlobalCheckingErrorBindJson(err.Error())
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

	account, err := a.accountService.GetAccountByUUID(ctx.Request.Context(), uuidAcc)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		helper.ReturnJSONError(ctx, http.StatusInternalServerError, "Internal server error", nil, nil)
		return
	}

	if account.Email == "" {
		log.Println("Error: Account not found")
		helper.ReturnJSONError(ctx, http.StatusNotFound, "Account not found", nil, nil)
		return
	}

	helper.ReturnJSON(ctx, http.StatusOK, "Account found", account)
}

func (a *AccountController) GetAccounts(ctx *gin.Context) {
	var req request.ListAccountRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		message, data := helper.GlobalCheckingErrorBindJson(err.Error())
		log.Printf("Error: %s", message)
		helper.ReturnJSONError(ctx, http.StatusBadRequest, message, nil, data)
		return
	}

	param := db.ListAccountsParams{
		Limit:  req.Limit,
		Offset: (req.Page - 1) * req.Limit,
	}

	accounts, totalData, err := a.accountService.ListAccount(ctx.Request.Context(), param)

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

func (a *AccountController) UpdateAccount(ctx *gin.Context) {
	var req request.GetAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		message, data := helper.GlobalCheckingErrorBindJson(err.Error())
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
		message, data := helper.GlobalCheckingErrorBindJson(err.Error())
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
		Status:      request.Status,
	}

	account, err := a.accountService.UpdateAccount(ctx.Request.Context(), arg)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		helper.ReturnJSONError(ctx, http.StatusInternalServerError, "Internal server error", nil, nil)
		return
	}

	if account.Email == "" {
		log.Println("Error: Account not found")
		helper.ReturnJSONError(ctx, http.StatusNotFound, "Account not found", nil, nil)
		return
	}

	helper.ReturnJSON(ctx, http.StatusOK, "Account updated", account)
}
