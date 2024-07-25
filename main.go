package main

import (
	"log"

	"github.com/fajaramaulana/simple_bank_project/internal/setup"
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

	setup.InitializeAndStartApp(config)
}
