package main

import (
	_ "github.com/joho/godotenv/autoload"
	"log"
	"mlauth/pkg/api"
)

func main() {
	err := api.Run()
	if err != nil {
		log.Fatalln("Failed to start API server:", err.Error())
	}
}
