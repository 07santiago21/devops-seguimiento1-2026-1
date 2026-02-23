package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/07santiago21/devops-seguimiento1-2026-1/internal/database"
)

func main() {

	_ = godotenv.Load()

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	_ = db.Debug()
}
