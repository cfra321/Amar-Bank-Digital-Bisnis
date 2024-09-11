package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/cfra321/Amar-Bank-Digital-Bisnis/structs"
)

// Generate a random 5-character ID
func generateAccountID() int {
	rand.Seed(time.Now().UnixNano())
	// Generate a random integer between 10000 and 99999
	min := 10000
	max := 99999
	return rand.Intn(max-min+1) + min
}

// generateAccountNumber generates an account number with format "1022" + 4 random digits
func generateAccountNumber() string {
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	randomDigits := fmt.Sprintf("%04d", seededRand.Intn(10000)) // Generate 4 random digits
	return "1022" + randomDigits
}

// formatBalance formats the balance to have 2 decimal places
func formatBalance(balance float64) float64 {
	return float64(int(balance*100)) / 100 // Format balance to 2 decimal places
}

// validateAccountType validates that the account_type is either 'internal' or 'external'
func validateAccountType(accountType string) error {
	if accountType != "internal" && accountType != "external" {
		return errors.New("account_type must be 'internal' or 'external'")
	}
	return nil
}

// CreateAccount inserts a new account into the database
func CreateAccount(db *sql.DB, account structs.Account) error {
	// Generate account ID, account number, and format balance
	account.ID = generateAccountID() // Generate account ID
	account.AccountNumber = generateAccountNumber()
	account.Balance = formatBalance(account.Balance)

	// Validate account_type
	if err := validateAccountType(account.AccountType); err != nil {
		return err
	}

	// Insert account into the database with created_at auto-generated as current timestamp
	query := `
		INSERT INTO accounts (id, user_id, account_number, balance, account_type, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, account.ID, account.UserID, account.AccountNumber, account.Balance, account.AccountType, time.Now())
	return err
}

// CheckAccountByUserID checks if an account with the same user_id already exists
func CheckAccountByUserID(db *sql.DB, userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM accounts WHERE user_id = $1)`
	err := db.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CheckUserIDExists checks if a user_id exists in the users table
func CheckUserIDExists(db *sql.DB, userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
	err := db.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetAllAccounts retrieves all accounts from the database
func GetAllAccounts(db *sql.DB) ([]structs.Account, error) {
	var accounts []structs.Account
	query := `SELECT id, user_id, account_number, balance, account_type, created_at FROM accounts`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account structs.Account
		err := rows.Scan(&account.ID, &account.UserID, &account.AccountNumber, &account.Balance, &account.AccountType, &account.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// GetAccountByID retrieves an account by its ID
func GetAccountByID(db *sql.DB, accountID int) (structs.Account, error) {
	var account structs.Account
	query := `SELECT id, user_id, account_number, balance, account_type, created_at FROM accounts WHERE id = $1`
	err := db.QueryRow(query, accountID).Scan(&account.ID, &account.UserID, &account.AccountNumber, &account.Balance, &account.AccountType, &account.CreatedAt)
	if err != nil {
		return account, err
	}
	return account, nil
}

func GetAccountByAccountNumber(db *sql.DB, accountNumber string) (structs.Account, error) {
	var account structs.Account
	query := "SELECT id, user_id, account_number, balance, account_type, created_at FROM accounts WHERE account_number = $1"

	err := db.QueryRow(query, accountNumber).Scan(&account.ID, &account.UserID, &account.AccountNumber, &account.Balance, &account.AccountType, &account.CreatedAt)
	if err != nil {
		return account, err
	}

	return account, nil
}
