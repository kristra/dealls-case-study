package main

import (
	"log"

	_ "dealls-case-study/docs"
	"dealls-case-study/internal/db"
	_ "dealls-case-study/internal/dto"

	"dealls-case-study/internal/route"

	"github.com/joho/godotenv"
)

// @title        Payroll System API
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
