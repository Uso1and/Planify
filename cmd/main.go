package main

import (
	"log"
	"planify/internal/app/routers"
	"planify/internal/domain/config"
	"planify/internal/domain/infrastructure/database"
)

func main() {
	if _, err := config.LoadConfig(); err != nil {
		log.Fatalf("failed to config:%v", err)
	}
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database :%v", err)
	}

	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Printf("Error closing database connect :%v", err)
		}
	}()

	r := routers.SetupPubRouter()

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
