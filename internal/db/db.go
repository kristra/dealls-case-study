package db

import (
	"dealls-case-study/internal/models"
	"fmt"
	"log"
	"os"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	var err error

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = RunMigrations(database)

	if err != nil {
		panic("Failed to auto-migrate database!")
	}

	DB = database

	log.Println("Database initialized!")
}

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "202506100930",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Attendance{}, &models.Overtime{}, &models.Payroll{}, &models.Payslip{}, &models.Reimbursement{}, &models.Role{}, &models.User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.Attendance{}, &models.Overtime{}, &models.Payroll{}, &models.Payslip{}, &models.Reimbursement{}, &models.Role{}, &models.User{})
			},
		},
	})

	return m.Migrate()
}
