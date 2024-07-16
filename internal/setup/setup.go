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

func InitializeAndStartApp() {
	// checking if already have .env file
	envPath := filepath.Join("/home/fajar/go_app/simplebankproject", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		// check if ENV already set
		if os.Getenv("DB_USER") == "" {
			log.Fatal("Error loading .env file")
		}
	}

	// / Get environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	PORT := os.Getenv("PORT")

	// Create the connection string
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	conn, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	store := db.NewStore(conn)

	// service
	accountService := service.NewAccountService(store)

	// controller
	accountController := controller.NewAccountController(accountService)

	server := router.NewRouter(accountController)

	server.SetupRouter()

	server.StartServer(PORT)
}
