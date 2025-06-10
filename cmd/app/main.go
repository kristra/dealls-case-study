package main

import (
	"log"

	"dealls-case-study/internal/db"
	_ "dealls-case-study/internal/dto"

	"dealls-case-study/internal/route"

	"github.com/joho/godotenv"
)

// @title        Payslip System API
// @version      1.0
// @description  Documentation for Payroll and Payslip management.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	_ = godotenv.Load()

	db.InitDB()
	route.SetupRoutes()

	log.Println("App started!")
}
