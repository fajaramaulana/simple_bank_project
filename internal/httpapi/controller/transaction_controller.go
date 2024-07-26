package controller

import (
	"log"
	"net/http"

	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/middleware"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/request"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/token"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/service"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	transactionService *service.TransactionService
}

func NewTransactionController(transactionService *service.TransactionService) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
	}
}

func (tf *TransactionController) CreateTransfer(ctx *gin.Context) {
	var req request.CreateTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		message, data := helper.GlobalCheckingErrorBindJson(err.Error(), req)
		log.Printf("Error: %s", message)
		helper.ReturnJSONError(ctx, http.StatusBadRequest, message, nil, data)
		return
	}

	res := helper.DoValidation(&req)

	if len(res) > 0 {
		log.Println("Error: Validation error")
		helper.ReturnJSONError(ctx, http.StatusBadRequest, "Validation error", nil, res)
		return
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

	transfer, err := tf.transactionService.CreateTransferTrans(ctx.Request.Context(), &req, authPayload)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		if err.Error() == "sql: no rows in result set: from account not found" {
			helper.ReturnJSONError(ctx, http.StatusNotFound, "from account not found", nil, nil)
			return
		}

		if err.Error() == "sql: no rows in result set: to account not found" {
			helper.ReturnJSONError(ctx, http.StatusNotFound, "to account not found", nil, nil)
			return
		}

		if err.Error() == "unauthorized" {
			helper.ReturnJSONError(ctx, http.StatusUnauthorized, "unauthorized", nil, nil)
			return
		}

		if err.Error() == "balance not enough" {
			helper.ReturnJSONError(ctx, http.StatusBadRequest, "balance not enough", nil, nil)
			return
		}

		helper.ReturnJSONError(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	helper.ReturnJSON(ctx, http.StatusCreated, "success transaction", transfer)
}
