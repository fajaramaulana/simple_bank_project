package setup

import (
	"database/sql"
	"fmt"
	"log"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/service"
	"github.com/fajaramaulana/simple_bank_project/util"

	_ "github.com/lib/pq"
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
	if conn == nil {
		log.Fatal("Failed to connect to database")
	} else {
		log.Println("Database connection successful")
	}

	store := db.NewStore(conn)
	log.Println("Store initialized")

	authService := service.NewAuthService(store, config)
	authController := controller.NewAuthController(authService)

	userService := service.NewUserService(store, config)
	userController := controller.NewUserController(userService)

	server, err := grpcapi.NewServer(store, authController, userController, config)
	if err != nil {
		log.Fatal("Cannot create gRPC server: ", err)
	}
	log.Println("gRPC server created")

	server.Start(config.GRPCPort)
	log.Println("gRPC server started")
}

// func InitializeAndStartGatewayServer(config util.Config) {
// 	conn := DbConnection(config)
// 	if conn == nil {
// 		log.Fatal("Failed to connect to database")
// 	} else {
// 		log.Println("Database connection successful")
// 	}

// 	store := db.NewStore(conn)
// 	log.Println("Store initialized")

// 	authService := service.NewAuthService(store, config)
// 	authController := controller.NewAuthController(authService)

// 	userService := service.NewUserService(store, config)
// 	userController := controller.NewUserController(userService)

// 	server, err := grpcapi.NewServer(store, authController, userController, config)
// 	if err != nil {
// 		log.Fatal("Cannot create gRPC server: ", err)
// 	}

// 	grpcMux := runtime.NewServeMux()
// }
