package controller

import (
	"log"
	"net/http"

	"github.com/fajaramaulana/simple_bank_project/internal/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/request"
	"github.com/fajaramaulana/simple_bank_project/internal/service"
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

	transfer, err := tf.transactionService.CreateTransferTrans(ctx.Request.Context(), &req)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		helper.ReturnJSONError(ctx, http.StatusInternalServerError, "Internal server error", nil, nil)
		return
	}

	helper.ReturnJSON(ctx, http.StatusCreated, "success transaction", transfer)
}
