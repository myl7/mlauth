package main

import (
	"log"
	"mlauth/pkg/api"
)

func main() {
	r := api.Route()
	err := r.Run()
	if err != nil {
		log.Fatalln("Failed to start API server:", err.Error())
	}
}
