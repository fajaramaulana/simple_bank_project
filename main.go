package main

import (
	"log"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/setup"
	"github.com/fajaramaulana/simple_bank_project/util"
)

func main() {
	// Load configuration from file
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	// Check if essential configuration values are set
	if config.DBUser == "" || config.DBPassword == "" || config.DBName == "" {
		log.Fatal("Environment variables are not properly loaded")
	}

	// Establish a database connection
	conn := setup.DbConnection(config)

	// Create a database store
	store := setup.GetDbStore(config, conn)

	// Start the gateway server in a separate goroutine
	go runGatewayServer(config, store)

	// Start the gRPC server
	rungRPCServer(config, store)
}

// func runGinServer(config util.Config, conn *sql.DB) {
// 	setuphttp.InitializeAndStartAppHTTPApi(config, conn)
// }

// rungRPCServer starts the gRPC server using the provided configuration and database store.
func rungRPCServer(config util.Config, store db.Store) {
	setup.InitializeAndStartAppGRPCApi(config, store)
}

// runGatewayServer starts the gateway server using the provided configuration and database store.
func runGatewayServer(config util.Config, store db.Store) {
	setup.InitializeAndStartGatewayServer(config, store)
}
