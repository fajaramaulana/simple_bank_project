package grpcapi

import (
	"fmt"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/token"
	"github.com/fajaramaulana/simple_bank_project/pb"
	"github.com/fajaramaulana/simple_bank_project/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	db         db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, db db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		db:         db,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
