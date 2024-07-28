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
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/fajaramaulana/simple_bank_project/doc/statik"
	_ "github.com/lib/pq"
)

// DbConnection establishes a connection to the database using the provided configuration.
// It returns a pointer to the sql.DB object representing the database connection.
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

// GetDbStore returns a new instance of db.Store using the provided configuration and database connection.
func GetDbStore(config util.Config, conn *sql.DB) db.Store {
	store := db.NewStore(conn)
	return store
}

// InitializeAndStartAppGRPCApi initializes and starts the gRPC API server for the application.
// It takes a configuration object and a database store as parameters.
// It creates instances of the required services and controllers,
// and then creates a gRPC server with the provided store, controllers, and configuration.
// Finally, it starts the gRPC server on the specified port from the configuration.
func InitializeAndStartAppGRPCApi(config util.Config, store db.Store) {

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

// InitializeAndStartGatewayServer initializes and starts the gRPC gateway server.
// It takes a `config` parameter of type `util.Config` which represents the server configuration,
// and a `store` parameter of type `db.Store` which represents the database store.
// The function creates instances of the required services and controllers,
// and then creates a gRPC server using the provided store, auth controller, user controller, and config.
// It registers the gRPC server with the gateway server and starts serving HTTP requests.
// The function also serves the Swagger UI using the provided statik file system.
// It listens on the configured port and logs the server startup message.
// If any error occurs during the initialization or serving, the function logs the error and exits.
func InitializeAndStartGatewayServer(config util.Config, store db.Store) {

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
