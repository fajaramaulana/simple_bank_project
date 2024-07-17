package service

import (
	"context"
	"errors"
	"strconv"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/request"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/response"
)

type TransactionService struct {
	db db.Store
}

func NewTransactionService(db db.Store) *TransactionService {
	return &TransactionService{
		db: db,
	}
}

func (a *TransactionService) CreateTransferTrans(ctx context.Context, req *request.CreateTransferRequest) (response.SuccessTransactionResponse, error) {
	// check if from account exists
	// convert string to uuid
	fromAccountUUID, err := helper.ConvertStringToUUID(req.FromAccountUUID)
	if err != nil {
		return response.SuccessTransactionResponse{}, err
	}
	dataFromAccount, err := a.db.GetAccountByUUID(ctx, fromAccountUUID)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = errors.New("from account not found")
		}
		return response.SuccessTransactionResponse{}, err
	}
	// end check if from account exists

	// check if to account exists
	// convert string to uuid
	toAccountUUID, err := helper.ConvertStringToUUID(req.ToAccountUUID)
	if err != nil {
		return response.SuccessTransactionResponse{}, err
	}

	dataToAccount, err := a.db.GetAccountByUUID(ctx, toAccountUUID)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = errors.New("to account not found")
		}
		return response.SuccessTransactionResponse{}, err
	}
	// end check if to account exists

	// convert existing balance to float64
	fromAccountBalance, err := strconv.ParseFloat(dataFromAccount.Balance, 64)
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
		return response.SuccessTransactionResponse{}, errors.New("currency not same")
	}

	if dataToAccount.Currency != req.Currency {
		return response.SuccessTransactionResponse{}, errors.New("currency not same")
	}

	if dataFromAccount.Currency != dataToAccount.Currency {
		return response.SuccessTransactionResponse{}, errors.New("currency not same")
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
		Amount:          res.Transaction.Amount,
		Currency:        dataFromAccount.Currency,
		LastedBalance:   res.FromAccount.Balance,
		Type:            "transfer",
	}, nil

}
