package database

import (
	"database/sql"
	"fmt"
	"log"
	"planify/internal/domain/config"

	_ "github.com/lib/pq" // <-- вот это важно
)

var DB *sql.DB

func Init() error {

	dbconfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config:%w", err)
	}
	configStr := dbconfig.GetConnectionString()

	DB, err = sql.Open("postgres", configStr)
	if err != nil {
		return fmt.Errorf("failed to connect database:%w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database:%w", err)
	}
	log.Println("Successfully connected to database")
	return nil

}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
