package main

import (
	"log"

	"dealls-case-study/internal/db"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	err := db.CreateDatabase()
	if err != nil {
		log.Fatalf("failed creating db: %v", err)
	}
}
