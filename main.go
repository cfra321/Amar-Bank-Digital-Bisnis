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
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
	)

	DB, err = sql.Open("postgres", psqlInfo)
	defer DB.Close()
	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	database.DBMigrate(DB)

	psqlSeeder := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
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
		router.POST("/api/register", controllers.RegisterUser)

		// Route untuk login
		router.POST("/api/auth/login", controllers.Login)

		authorized := router.Group("/")
		authorized.Use(middleware.JWTAuthMiddleware())
		{
			// Endpoint untuk update status pengguna
			authorized.PUT("/api/users/:id/status", controllers.UpdateUserStatus)
			authorized.GET("/api/users", controllers.GetAllUsers)
			authorized.GET("/api/users/:id", controllers.GetUserByID)
			authorized.DELETE("/api/users/:id", controllers.DeleteUser)

			authorized.POST("/api/accounts", controllers.CreateAccount)
			authorized.GET("/api/accounts", controllers.GetAllAccounts)
			authorized.GET("/api/:id/accounts", controllers.GetAccountsByID)
			authorized.GET("/api/account-number/:account_number", controllers.GetAccountByAccountNumber)

			authorized.POST("/api/transactions/transfer", controllers.TransferOverBooking)

			authorized.GET("/api/transaction_logs", controllers.GetAllTransactionLogs)
		}

	}

	router.Run(":" + os.Getenv("PORT"))
	// router.Run(":8080")
}
