package setup

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/server"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/service"
	"github.com/fajaramaulana/simple_bank_project/pb"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"

	_ "github.com/fajaramaulana/simple_bank_project/doc/statik"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/encoding/protojson"
)

func DbConnection(config util.Config) *sql.DB {
	dbUser := config.DBUser
	dbPassword := config.DBPassword
	dbHost := config.DBHost
	dbPort := config.DBPort
	dbName := config.DBName
	dbSSLMode := config.DBSSLMode
	// Create the connection string
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	conn, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	return conn
}

func InitializeAndStartAppGRPCApi(config util.Config) {
	conn := DbConnection(config)
	db.NewStore(conn)
	store := db.NewStore(conn)

	authService := service.NewAuthService(store, config)
	authController := controller.NewAuthController(authService)

	userService := service.NewUserService(store, config)
	userController := controller.NewUserController(userService)

	server, err := server.NewServer(store, authController, userController, config)
	if err != nil {
		log.Fatal("Cannot create gRPC server: ", err)
	}

	server.Start(config.GRPCPort)
}

func InitializeAndStartGatewayServer(config util.Config) {
	conn := DbConnection(config)

	store := db.NewStore(conn)

	authService := service.NewAuthService(store, config)
	authController := controller.NewAuthController(authService)

	userService := service.NewUserService(store, config)
	userController := controller.NewUserController(userService)

	server, err := server.NewServer(store, authController, userController, config)
	if err != nil {
		log.Fatal("Cannot create gRPC server: ", err)
	}

	jsonOpt := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOpt)
	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("Failed to register gRPC server: ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("Cannot create statik file system: ", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", ":"+config.PortGatewayGrpc)
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	log.Printf("Starting gRPC gateway server on %s", config.PortGatewayGrpc)

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("Failed to serve: ", err)
	}
	log.Println("gRPC gateway server started")
}
