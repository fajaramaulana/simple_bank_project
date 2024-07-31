package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/request"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/response"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/token"
)

type TransactionService struct {
	db db.Store
}

func NewTransactionService(db db.Store) *TransactionService {
	return &TransactionService{
		db: db,
	}
}

func (a *TransactionService) CreateTransferTrans(ctx context.Context, req *request.CreateTransferRequest, authPayload *token.Payload) (response.SuccessTransactionResponse, error) {
	fromAccountUUID, err := helper.ConvertStringToUUID(req.FromAccountUUID)
	if err != nil {
		return response.SuccessTransactionResponse{}, fmt.Errorf("%w: from account uuid not valid", err)
	}
	dataFromAccount, err := a.db.GetAccountByUUID(ctx, fromAccountUUID)

	if err != nil {
		return response.SuccessTransactionResponse{}, fmt.Errorf("%w: from account not found", err)
	}

	fmt.Printf("%# v\n", dataFromAccount.UserUuid.String())
	fmt.Printf("%# v\n", authPayload.UserUUID.String())
	// check if from account belongs to user
	if dataFromAccount.UserUuid != authPayload.UserUUID {
		return response.SuccessTransactionResponse{}, errors.New("unauthorized")
	}
	// end check if from account exists

	// check if to account exists
	// convert string to uuid
	toAccountUUID, err := helper.ConvertStringToUUID(req.ToAccountUUID)
	if err != nil {
		return response.SuccessTransactionResponse{}, fmt.Errorf("%w: to account uuid not valid", err)
	}

	dataToAccount, err := a.db.GetAccountByUUID(ctx, toAccountUUID)

	if err != nil {
		return response.SuccessTransactionResponse{}, fmt.Errorf("%w: to account not found", err)
	}
	// end check if to account exists

	// convert existing balance to float64
	fromAccountBalance, err := strconv.ParseFloat(dataFromAccount.Balance.Int.String(), 64)
	if err != nil {
		return response.SuccessTransactionResponse{}, err
	}

	// convert transfer amount from int64 to float64
	transferAmount := float64(req.Amount)

	// checking if balance enough
	if fromAccountBalance < transferAmount {
		return response.SuccessTransactionResponse{}, errors.New("balance not enough")
	}

	// checking if currency same
	if dataFromAccount.Currency != req.Currency {
		return response.SuccessTransactionResponse{}, errors.New("currency from account not same")
	}

	if dataToAccount.Currency != req.Currency {
		return response.SuccessTransactionResponse{}, errors.New("currency to account not same")
	}

	if dataFromAccount.Currency != dataToAccount.Currency {
		return response.SuccessTransactionResponse{}, errors.New("both account currency not same")
	}

	// create transaction
	arg := db.TransferTxParam{
		FromAccountID: dataFromAccount.ID,
		ToAccountID:   dataToAccount.ID,
		Amount:        int64(transferAmount),
		Type:          "transfer",
	}

	res, err := a.db.TransferTx(ctx, arg)

	if err != nil {
		return response.SuccessTransactionResponse{}, err
	}

	return response.SuccessTransactionResponse{
		TransactionUUID: res.Transaction.TransactionUuid.String(),
		FromAccountUUID: res.FromAccount.AccountUuid.String(),
		ToAccountUUID:   res.ToAccount.AccountUuid.String(),
		Amount:          res.Transaction.Amount.Int.String(),
		Currency:        dataFromAccount.Currency,
		LastedBalance:   res.FromAccount.Balance.Int.String(),
		Type:            "transfer",
	}, nil

}
