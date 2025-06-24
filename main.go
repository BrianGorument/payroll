package main

import (
	"log"

	"payroll/config"
	"payroll/database"
	"payroll/logger"
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
	db, err := database.InitDBMysql()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Gin Router
	router := gin.Default()

	// Register User Routes
	users.RegisterRoutes(router, db, log)

	// Start Server
	port := ":8888"
	log.Infof("Starting server on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
