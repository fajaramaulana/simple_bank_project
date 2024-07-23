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

	if config.Port == "" {
		log.Fatal("env not load properly")
	}
	setup.InitializeAndStartApp(config)
}
