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

func checkingEnv() {
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

func dbConnection() *sql.DB {
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
	checkingEnv()
	// / Get environment variables
	conn := dbConnection()
	store := db.NewStore(conn)

	// account
	accountService := service.NewAccountService(store)
	accountController := controller.NewAccountController(accountService)

	// transfer

	transferService := service.NewTransactionService(store)
	transferController := controller.NewTransactionController(transferService)

	server := router.NewRouter(accountController, transferController)

	server.SetupRouter()

	PORT := os.Getenv("PORT")
	server.StartServer(PORT)
}
