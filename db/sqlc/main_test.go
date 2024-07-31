package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

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

	// testDB, err = sql.Open("postgres", connStr)
	connPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
