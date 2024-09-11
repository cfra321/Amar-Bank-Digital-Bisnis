package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cfra321/Amar-Bank-Digital-Bisnis/database"
	"github.com/cfra321/Amar-Bank-Digital-Bisnis/repository"
	"github.com/cfra321/Amar-Bank-Digital-Bisnis/structs"
	"github.com/gin-gonic/gin"
)

func TransferOverBooking(c *gin.Context) {
	var transfer structs.Transfer
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi agar amount tidak boleh minus
	if transfer.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than zero"})
		return
	}

	// Validasi akun pengirim dan penerima
	senderAccount, err := repository.GetAccountByID(database.DbConnection, transfer.SenderAccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sender account not found"})
		return
	}
	receiverAccount, err := repository.GetAccountByID(database.DbConnection, transfer.ReceiverAccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Receiver account not found"})
		return
	}

	// Cek status akun aktif
	if !repository.IsUserActive(database.DbConnection, senderAccount.UserID) || !repository.IsUserActive(database.DbConnection, receiverAccount.UserID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One of the accounts is inactive"})
		return
	}

	// Cek apakah saldo cukup
	if transfer.TransferType == "bifast" {
		if senderAccount.Balance < transfer.Amount+2500 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance for bifast transfer"})
			return
		}
	} else {
		if senderAccount.Balance < transfer.Amount {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
			return
		}
	}

	// Cek apakah receiver account adalah external dan transfer type bukan bifast
	if receiverAccount.AccountType == "external" && transfer.TransferType != "bifast" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer to external account without bifast"})
		return
	}

	// Atur biaya transfer jika menggunakan bifast
	if transfer.TransferType == "bifast" {
		transfer.Fee = 2500.00 // contoh biaya bifast
	} else {
		transfer.Fee = 0
	}

	// Update saldo pengirim dan penerima
	err = repository.UpdateAccountBalances(database.DbConnection, senderAccount.ID, receiverAccount.ID, transfer.Amount, transfer.Fee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process transfer"})
		return
	}

	// Simpan transfer dengan userID yang melakukan transfer
	transfer.Status = "completed"
	transfer.CreatedAt = time.Now()

	if err := repository.CreateTransaction(database.DbConnection, &transfer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Simpan log transaksi
	logMessage := fmt.Sprintf("Transfer ID %d: %s of %.2f completed.", transfer.ID, transfer.TransferType, transfer.Amount)
	transactionLog := structs.TransactionLog{
		TransferID: transfer.ID,
		LogMessage: logMessage,
		CreatedAt:  time.Now(),
	}

	if err := repository.CreateTransactionLog(database.DbConnection, &transactionLog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transfer})
}
