package controllers

import (
	"net/http"
	"strconv"

	"github.com/cfra321/Amar-Bank-Digital-Bisnis/database"
	"github.com/cfra321/Amar-Bank-Digital-Bisnis/repository"
	"github.com/cfra321/Amar-Bank-Digital-Bisnis/structs"
	"github.com/gin-gonic/gin"
)

func CreateAccount(c *gin.Context) {
	var account structs.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := account.UserID

	userExists, err := repository.CheckUserIDExists(database.DbConnection, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}
	if !userExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	// Check if an account with the same user_id already exists
	accountExists, err := repository.CheckAccountByUserID(database.DbConnection, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}
	if accountExists {
		c.JSON(http.StatusConflict, gin.H{"error": "Account with this user_id already exists"})
		return
	}

	err = repository.CreateAccount(database.DbConnection, account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Account created successfully"})
}

func GetAllAccounts(c *gin.Context) {
	accounts, err := repository.GetAllAccounts(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func GetAccountsByID(c *gin.Context) {
	accountID := c.Param("id")
	id, err := strconv.Atoi(accountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := repository.GetAccountByID(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

// GetAccountByAccountNumber retrieves account details by account number
func GetAccountByAccountNumber(c *gin.Context) {
	accountNumber := c.Param("account_number")
	account, err := repository.GetAccountByAccountNumber(database.DbConnection, accountNumber)
	if err != nil {
		if err.Error() == "account not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch account details", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}
