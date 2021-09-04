package main

import (
	"github.com/joho/godotenv"
	"log"
	"mlauth/pkg/api"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Failed to load dotenv:", err.Error())
	}

	err = api.Run()
	if err != nil {
		log.Fatalln("Failed to start API server:", err.Error())
	}
}
