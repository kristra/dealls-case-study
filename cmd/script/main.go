package main

import (
	"log"
	"os"

	"dealls-case-study/internal/db"
	"dealls-case-study/internal/seed"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	opt := os.Args[1]

	switch opt {
	case "create":
		{
			err := db.CreateDatabase()
			if err != nil {
				log.Fatalf("failed creating db: %v", err)
			}
		}

	case "drop":
		{
			err := db.DropDatabase()
			if err != nil {
				log.Fatalf("failed dropping db: %v", err)
			}
		}

	case "seed":
		{
			db.InitDB()
			err := seed.Run(db.DB)
			if err != nil {
				log.Fatalf("failed seeding data: %v", err)
			}
		}

	case "init":
		{
			err := db.DropDatabase()
			if err != nil {
				log.Fatalf("failed dropping db: %v", err)
			}
			err = db.CreateDatabase()
			if err != nil {
				log.Fatalf("failed creating db: %v", err)
			}
			db.InitDB()
			err = seed.Run(db.DB)
			if err != nil {
				log.Fatalf("failed seeding data: %v", err)
			}
		}

	// helper function to easily test salary calculations
	case "seed44":
		{
			db.InitDB()
			seed.Seed44(db.DB)
		}

	case "reset":
		{
			err := db.DropDatabase()
			if err != nil {
				log.Fatalf("failed dropping db: %v", err)
			}
			err = db.CreateDatabase()
			if err != nil {
				log.Fatalf("failed creating db: %v", err)
			}
			db.InitDB()
			err = seed.Run(db.DB)
			if err != nil {
				log.Fatalf("failed seeding data: %v", err)
			}
			seed.Seed44(db.DB)
		}
	}
}
