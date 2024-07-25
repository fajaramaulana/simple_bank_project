package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/fajaramaulana/simple_bank_project/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// checking if already have .env file
	config, err := util.LoadConfig("../../.")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	// Get environment variables
	dbUser := config.DBUser
	dbPassword := config.DBPassword
	dbHost := config.DBHost
	dbPort := config.DBPort
	dbName := config.DBName
	dbSSLMode := config.DBSSLMode

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
