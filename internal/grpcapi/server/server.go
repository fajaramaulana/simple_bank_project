package grpcapi

import (
	"fmt"

	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/token"
	"github.com/fajaramaulana/simple_bank_project/pb"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	account     *controller.AccountController
	transaction *controller.TransactionController
	user        *controller.UserController
	auth        *controller.AuthController
	TokenMaker  token.Maker
}

func NewServer(account *controller.AccountController, transaction *controller.TransactionController, user *controller.UserController, auth *controller.AuthController, configToken map[string]string) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(configToken["token_secret"])
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		account:     account,
		transaction: transaction,
		user:        user,
		auth:        auth,
		TokenMaker:  tokenMaker,
	}

	return server, nil
}
