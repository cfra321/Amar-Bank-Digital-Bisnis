package controllers

import (
	"net/http"

	"github.com/cfra321/Amar-Bank-Digital-Bisnis/database"
	"github.com/cfra321/Amar-Bank-Digital-Bisnis/repository"
	"github.com/gin-gonic/gin"
)

// GetAllTransactionLogs untuk mendapatkan semua log transaksi
func GetAllTransactionLogs(c *gin.Context) {
	transactionLogs, err := repository.GetAllTransactionLogs(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactionLogs)
}
