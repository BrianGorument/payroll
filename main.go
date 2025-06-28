package main

import (
	"log"

	"payroll/config"
	"payroll/database"
	"payroll/logger"
	"payroll/shared/utils"
	"payroll/src/attendances"
	"payroll/src/overtimes"
	"payroll/src/payroll_periods"
	"payroll/src/users"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Config
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Logger
	log := logger.NewLogger()

	// Initialize Database
	db, err := database.InitDBPostgres()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Gin Router
	router := gin.Default()
	
    router.Use(utils.AuditLogMiddleware(db))

	// Register User Routes
	users.RegisterRoutes(router, db, log)
	payroll_periods.RegisterRoutes(router, db, log)
	attendances.RegisterRoutes(router, db, log)
	overtimes.RegisterRoutes(router, db, log)

	// Start Server
	port := ":8888"
	log.Infof("Starting server on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
