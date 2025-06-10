package main

import (
	"log"

	"dealls-case-study/internal/db"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db.InitDB()

	log.Println("App started!")
}
