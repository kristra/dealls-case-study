package db

import (
	"context"
	"dealls-case-study/internal/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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

func InitTestDB() (*gorm.DB, func(), error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		// panic(err)
		return nil, nil, err
	}

	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("host=%s port=%s user=postgres password=password dbname=testdb sslmode=disable", host, port.Port())

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// panic("failed to connect to database")
		return nil, nil, err
	}

	db.AutoMigrate(&models.Attendance{}, &models.Overtime{}, &models.Payroll{}, &models.Payslip{}, &models.Reimbursement{}, &models.Role{}, &models.User{})

	DB = db

	roles := []models.Role{
		{Name: "Admin", CreatedBy: 999},
		{Name: "Employee", CreatedBy: 999},
	}

	// seed roles
	for _, role := range roles {
		err := db.FirstOrCreate(&role, models.Role{Name: role.Name}).Error
		if err != nil {
			log.Printf("Failed to seed role %s: %v", role.Name, err)
			return nil, nil, err
		}
	}

	// Tear down the container after tests
	// defer postgresC.Terminate(ctx)
	cleanup := func() {
		postgresC.Terminate(ctx)
	}

	return db, cleanup, nil
}
