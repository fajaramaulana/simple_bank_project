package server

import (
	"fmt"
	"net"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/handler/token"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/logger"
	"github.com/fajaramaulana/simple_bank_project/pb"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("Failed to create token maker")
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
		log.Fatal().Err(err).Msg("Failed to listen")
	}
	log.Info().Msgf("Start gRPC server at port: %s", port)

	grpcLogger := grpc.UnaryInterceptor(logger.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, s)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
}
