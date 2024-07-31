package setup

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/logger"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/seed"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/server"
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/service"
	"github.com/fajaramaulana/simple_bank_project/pb"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/fajaramaulana/simple_bank_project/doc/statik"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

// DbConnection establishes a connection to the database using the provided configuration.
// It returns a pointer to the sql.DB object representing the database connection.
func DbConnection(config util.Config) *pgxpool.Pool {
	dbUser := config.DBUser
	dbPassword := config.DBPassword
	dbHost := config.DBHost
	dbPort := config.DBPort
	dbName := config.DBName
	dbSSLMode := config.DBSSLMode

	// Create the connection string
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	conn, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot open a connection to the database")
	}

	return conn
}

// RedisConnection establishes a connection to Redis using the provided configuration.
// It returns a pointer to a redis.Client that can be used to interact with Redis.
func RedisConnection(config util.Config) *redis.Client {
	redissb, err := strconv.Atoi(config.RedisDB)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot convert Redis DB to integer")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: "",
		DB:       redissb,
	})
	return client
}

// GetDbStore returns a new instance of db.Store using the provided configuration and database connection.
func GetDbStore(config util.Config, conn *pgxpool.Pool) db.Store {
	store := db.NewStore(conn)
	return store
}

// InitializeAndStartAppGRPCApi initializes and starts the gRPC API server for the application.
// It takes a configuration object and a database store as parameters.
// It creates instances of the required services and controllers,
// and then creates a gRPC server with the provided store, controllers, and configuration.
// Finally, it starts the gRPC server on the specified port from the configuration.
func InitializeAndStartAppGRPCApi(config util.Config, store db.Store, redisClient *redis.Client) {

	authService := service.NewAuthService(store, config)
	authController := controller.NewAuthController(authService)

	userService := service.NewUserService(store, config, redisClient)
	userController := controller.NewUserController(userService)

	server, err := server.NewServer(store, authController, userController, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create gRPC server")
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
func InitializeAndStartGatewayServer(config util.Config, store db.Store, redisClient *redis.Client) {

	authService := service.NewAuthService(store, config)
	authController := controller.NewAuthController(authService)

	userService := service.NewUserService(store, config, redisClient)
	userController := controller.NewUserController(userService)

	server, err := server.NewServer(store, authController, userController, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create gRPC server")
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot register gRPC server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create statik file system")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", ":"+config.PortGatewayGrpc)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot listen to the port")
	}

	log.Printf("Starting gRPC gateway server on %s", config.PortGatewayGrpc)
	handler := logger.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start gRPC gateway server")
	}
	log.Info().Msg("gRPC gateway server started")
}

// InitializeDBMigrationsAndSeeder initializes the database migrations and seeder.
// It takes a `config` object of type `util.Config` and a `conn` object of type `*sql.DB`.
// The `config` object contains the database configuration details such as DBUser, DBPassword, DBHost, DBPort, DBName, and DBSSLMode.
// The `conn` object is the database connection object.
// This function performs the following steps:
// 1. Creates a database migration using the provided configuration details.
// 2. Applies the migration to the database.
// 3. Logs any errors that occur during migration.
// 4. Logs a success message if the migration is successful.
// 5. Creates a seeder object using the provided database connection.
// 6. Seeds the database with initial data using the seeder object.
func InitializeDBMigrationsAndSeeder(config util.Config, conn *pgxpool.Pool) {
	// migration
	dbConf := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName, config.DBSSLMode)
	m, err := migrate.New(config.DBSource, dbConf)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create migration")
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("Cannot migrate the database")
	}
	if err != migrate.ErrNoChange {
		log.Info().Msg("Database migration successful")
	}

	seed := seed.NewSeeder(conn)
	seed.Seed()
}
