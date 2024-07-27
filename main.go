package main

import (
	"log"

	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/setup"
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

	go runGatewayServer(config)
	rungRPCServer(config)
}

// func runGinServer(config util.Config) {
// 	setuphttp.InitializeAndStartAppHTTPApi(config)
// }

func rungRPCServer(config util.Config) {
	setup.InitializeAndStartAppGRPCApi(config)
}

func runGatewayServer(config util.Config) {
	setup.InitializeAndStartGatewayServer(config)
}
