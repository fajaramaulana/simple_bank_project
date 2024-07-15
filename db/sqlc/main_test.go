package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// checking if already have .env file
	envPath := filepath.Join("/home/fajar/go_app/simplebankproject", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Printf("%# v\n", os.Getenv("DB_USER"))
		log.Fatalf("Error loading .env file")
	}

	// Get environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Create the connection string
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	testDB, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
