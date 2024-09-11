package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/cfra321/Amar-Bank-Digital-Bisnis/database"
	"github.com/cfra321/Amar-Bank-Digital-Bisnis/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/cfra321/Amar-Bank-Digital-Bisnis/controllers"
	"github.com/cfra321/Amar-Bank-Digital-Bisnis/seeder"
)

var (
	DB  *sql.DB
	err error
)

func main() {
	err = godotenv.Load("config/.env")
	if err != nil {
		panic("Error loading .env file")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	DB, err = sql.Open("postgres", psqlInfo)
	defer DB.Close()
	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	database.DBMigrate(DB)

	psqlSeeder := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	DB, err = sql.Open("postgres", psqlSeeder)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer DB.Close()

	// Seed the database
	seeder.SeedUsers(DB)
	fmt.Println("Successfully connected!")

	router := gin.Default()

	// Use the CORSMiddleware
	router.Use(middleware.CORSMiddleware())
	{

		// Route untuk register
		router.POST("/register", controllers.RegisterUser)

		// Route untuk login
		router.POST("/login", controllers.Login)

		authorized := router.Group("/")
		authorized.Use(middleware.JWTAuthMiddleware())
		{
			// Endpoint untuk update status pengguna
			authorized.PUT("/users/:id/status", controllers.UpdateUserStatus)
			authorized.GET("/users", controllers.GetAllUsers)
			authorized.GET("/users/:id", controllers.GetUserByID)
			authorized.DELETE("/users/:id", controllers.DeleteUser)

			authorized.POST("/users/accounts", controllers.CreateAccount)
			authorized.GET("/users/:id/accounts", controllers.GetAccountsByID)
			authorized.GET("/users/accounts", controllers.GetAllAccounts)
			authorized.GET("/users/account-number/:account_number", controllers.GetAccountByAccountNumber)

			authorized.POST("/users/transfer", controllers.TransferOverBooking)

			authorized.GET("/users/transaction_logs", controllers.GetAllTransactionLogs)
		}

	}

	router.Run(":8080")
}
