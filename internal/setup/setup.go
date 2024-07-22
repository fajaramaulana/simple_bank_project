package setup

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/router"
	"github.com/fajaramaulana/simple_bank_project/internal/service"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func CheckingEnv() {
	// checking if already have .env file
	envPath := filepath.Join("/home/fajar/go_app/simplebankproject", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		// check if ENV already set
		if os.Getenv("DB_USER") == "" {
			log.Fatal("Error loading .env file")
		}
	}
}

func DbConnection() *sql.DB {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Create the connection string
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	conn, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	return conn
}

func InitializeAndStartApp() {
	CheckingEnv()
	// / Get environment variables
	conn := DbConnection()
	store := db.NewStore(conn)

	// configToken
	configToken := map[string]string{
		"token_secret":          os.Getenv("TOKEN_SYMMETRIC_KEY"),
		"access_token_duration": os.Getenv("ACCESS_TOKEN_DURATION"),
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

	server.SetupRouter()

	PORT := os.Getenv("PORT")
	server.StartServer(PORT)
}

func InitializeAndStartAppTest() *router.Router {
	CheckingEnv()
	// / Get environment variables
	conn := DbConnection()
	store := db.NewStore(conn)

	// configToken
	configToken := map[string]string{
		"token_secret":          os.Getenv("TOKEN_SYMMETRIC_KEY"),
		"access_token_duration": os.Getenv("ACCESS_TOKEN_DURATION"),
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

	return server
}
