package main

import (
	"log"

	"planify/internal/app/routers"
	"planify/internal/domain/infrastructure/database"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	r := routers.SetupPubRouter()

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
