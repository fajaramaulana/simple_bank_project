package main

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/runner"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/setup"
	"github.com/fajaramaulana/simple_bank_project/util"
)

func main() {
	// Load configuration from file
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
	}

	// Check if essential configuration values are set
	if config.DBUser == "" || config.DBPassword == "" || config.DBName == "" {
		log.Fatal().Msg("Environment variables are not properly loaded")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Establish a database connection
	conn := setup.DbConnection(config)

	// Create a database store
	store := setup.GetDbStore(config, conn)

	redisClient := setup.RedisConnection(config)

	setup.InitializeDBMigrationsAndSeeder(config, conn)

	// Start the gateway server in a separate goroutine
	go runGatewayServer(config, store, redisClient)

	go runner.SendVerificationEmails(context.Background(), redisClient, config)

	// Start the gRPC server
	rungRPCServer(config, store, redisClient)
}

// func runGinServer(config util.Config, conn *sql.DB) {
// 	setuphttp.InitializeAndStartAppHTTPApi(config, conn)
// }

// rungRPCServer starts the gRPC server using the provided configuration and database store.
func rungRPCServer(config util.Config, store db.Store, redisClient *redis.Client) {
	setup.InitializeAndStartAppGRPCApi(config, store, redisClient)
}

// runGatewayServer starts the gateway server using the provided configuration and database store.
func runGatewayServer(config util.Config, store db.Store, redisClient *redis.Client) {
	setup.InitializeAndStartGatewayServer(config, store, redisClient)
}
