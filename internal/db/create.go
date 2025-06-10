package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func CreateDatabase() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		host, port, user, password, sslmode,
	)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	var exists bool
	err = sqlDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbname).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Database %s already exists", dbname)
		return nil
	}

	_, err = sqlDB.Exec("CREATE DATABASE " + dbname)
	if err != nil {
		return err
	}
	log.Printf("Database %s created successfully", dbname)
	return nil
}
