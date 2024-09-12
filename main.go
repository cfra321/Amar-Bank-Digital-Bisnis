package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/cfra321/Amar-Bank-Digital-Bisnis/controllers"
	"github.com/cfra321/Amar-Bank-Digital-Bisnis/database"
	"github.com/cfra321/Amar-Bank-Digital-Bisnis/middleware"
	"github.com/cfra321/Amar-Bank-Digital-Bisnis/seeder"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	DB  *sql.DB
	err error
)

func main() {
	// Load environment variables from .env file
	err = godotenv.Load("config/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Database connection
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		// os.Getenv("DB_HOST"),
		// os.Getenv("PGHOST"),
		// os.Getenv("DB_USER"),
		// os.Getenv("DB_PASSWORD"),
		// os.Getenv("DB_NAME"),
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
	)

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer DB.Close()

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Could not ping the database: %v", err)
	}

	// Run migrations
	database.DBMigrate(DB)

	// Seed the database with initial data
	seeder.SeedUsers(DB)

	fmt.Println("Successfully connected to the database!")

	// Setup Gin router
	router := gin.Default()

	// Use CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Public routes
	router.POST("/api/register", controllers.RegisterUser)
	router.POST("/api/auth/login", controllers.Login)

	// Protected routes with JWT middleware
	authorized := router.Group("/")
	authorized.Use(middleware.JWTAuthMiddleware())
	{
		// User management routes
		authorized.PUT("/api/users/:id/status", controllers.UpdateUserStatus)
		authorized.GET("/api/users", controllers.GetAllUsers)
		authorized.GET("/api/users/:id", controllers.GetUserByID)
		authorized.DELETE("/api/users/:id", controllers.DeleteUser)

		// Account management routes
		authorized.POST("/api/accounts", controllers.CreateAccount)
		authorized.GET("/api/accounts", controllers.GetAllAccounts)
		authorized.GET("/api/:id/accounts", controllers.GetAccountsByID)
		authorized.GET("/api/account-number/:account_number", controllers.GetAccountByAccountNumber)

		// Transaction routes
		authorized.POST("/api/transactions/transfer", controllers.TransferOverBooking)
		authorized.GET("/api/transaction_logs", controllers.GetAllTransactionLogs)
	}

	// Start the server
	router.Run(":" + os.Getenv("PORT"))
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// }
	// router.Run(":" + port)
}
