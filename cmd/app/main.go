package main

import (
	"log"

	"dealls-case-study/internal/db"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("App started!")
}
