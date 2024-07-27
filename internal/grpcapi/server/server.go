package server

import (
	"fmt"
	"log"
	"net"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/token"
	"github.com/fajaramaulana/simple_bank_project/pb"
	"github.com/fajaramaulana/simple_bank_project/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	// grpcServer     *grpc.Server
	config         util.Config
	authController *controller.AuthController
	userController *controller.UserController
	tokenMaker     token.Maker
}

func NewServer(store db.Store, authController *controller.AuthController, userController *controller.UserController, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		// grpcServer:     grpcServer,
		config:         config,
		authController: authController,
		userController: userController,
		tokenMaker:     tokenMaker,
	}
	return server, nil
}

// Start runs the gRPC server on the specified port.
func (s *Server) Start(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Starting gRPC server on %s", port)

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, s)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
