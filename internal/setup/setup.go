package setup

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/router"
	"github.com/fajaramaulana/simple_bank_project/internal/service"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func CheckingEnv(config util.Config) util.Config {
	return config
}

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

func InitializeAndStartApp(config util.Config) {
	// / Get environment variables
	conn := DbConnection(config)

	// checking table user is empty or not and return count
	var count int
	row := conn.QueryRow("SELECT COUNT(*) FROM users")
	err := row.Scan(&count)
	if err != nil {
		log.Fatal("Cannot check table users: ", err)
	}

	if count == 0 {
		defaultPass, err := util.MakePasswordBcrypt("Passw0rd!")
		if err != nil {
			log.Fatal("Cannot create password: ", err)
		}

		// insert data to table user
		queries := []string{
			"INSERT INTO users (username, full_name, email, hashed_password, role) VALUES ('admin', 'administrator', 'admin@simplebank.org', '" + defaultPass + "', 'admin')",
			"INSERT INTO users (username, full_name, email, hashed_password, role) VALUES ('user', 'customer', 'fajar1@gmail.com', '" + defaultPass + "', 'customer')",
		}

		for _, query := range queries {
			_, err := conn.Exec(query)
			if err != nil {
				log.Fatal("Cannot insert data to DB: ", err)
			}
		}
	}

	// Create a new store
	store := db.NewStore(conn)

	// configToken
	configToken := map[string]string{
		"token_secret":           config.TokenSymmetricKey,
		"access_token_duration":  config.AccessTokenDuration.String(),
		"refresh_token_duration": config.RefreshTokenDuration.String(),
	}
	// account
	accountService := service.NewAccountService(store)
	accountController := controller.NewAccountController(accountService)

	// transfer
	transferService := service.NewTransactionService(store)
	transferController := controller.NewTransactionController(transferService)

	// user
	userService := service.NewUserService(store)
	userController := controller.NewUserController(userService)

	// auth
	authService := service.NewAuthService(store, configToken)
	authController := controller.NewAuthController(authService)

	server, err := router.NewRouter(accountController, transferController, userController, authController, configToken)
	if err != nil {
		log.Fatal("Cannot create router: ", err)
	}

	PORT := config.Port
	server.StartServer(PORT)
}

func InitializeAndStartAppTest(t *testing.T, store db.Store) *router.Router {
	// Check environment variables

	// Config token
	configToken := map[string]string{
		"token_secret":           util.RandomString(32),
		"access_token_duration":  time.Minute.String(),
		"refresh_token_duration": (15 * time.Minute).String(),
	}

	// Initialize services and controllers
	accountService := service.NewAccountService(store)
	accountController := controller.NewAccountController(accountService)

	transferService := service.NewTransactionService(store)
	transferController := controller.NewTransactionController(transferService)

	userService := service.NewUserService(store)
	userController := controller.NewUserController(userService)

	authService := service.NewAuthService(store, configToken)
	authController := controller.NewAuthController(authService)

	// Create router
	server, err := router.NewRouter(accountController, transferController, userController, authController, configToken)
	require.NoError(t, err)

	return server
}
