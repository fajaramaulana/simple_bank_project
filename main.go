package main

import (
	"log"

	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/setup"
	"github.com/fajaramaulana/simple_bank_project/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	// Check if essential configuration values are set
	if config.DBUser == "" || config.DBPassword == "" || config.DBName == "" {
		log.Fatal("Environment variables are not properly loaded")
	}

	// runGinServer(config)
	rungRPCServer(config)
}

// func runGinServer(config util.Config) {
// 	setup.InitializeAndStartApp(config)
// }

func rungRPCServer(config util.Config) {
	setup.InitializeAndStartAppGRPCApi(config)
}
