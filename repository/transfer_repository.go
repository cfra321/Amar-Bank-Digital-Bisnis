package repository

import (
	"database/sql"

	"github.com/cfra321/Amar-Bank-Digital-Bisnis/structs"
)

// IsUserActive untuk mengecek apakah pengguna aktif
func IsUserActive(db *sql.DB, userID int) bool {
	var isActive bool
	err := db.QueryRow("SELECT is_active FROM users WHERE id = $1", userID).Scan(&isActive)
	if err != nil || !isActive {
		return false
	}
	return true
}

// UpdateAccountBalances untuk mengurangi saldo pengirim dan menambah saldo penerima
func UpdateAccountBalances(db *sql.DB, senderID, receiverID int, amount, fee float64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Kurangi saldo pengirim
	_, err = tx.Exec("UPDATE accounts SET balance = balance - $1 - $2 WHERE id = $3", amount, fee, senderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Tambah saldo penerima
	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, receiverID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaksi
	return tx.Commit()
}

// CreateTransaction untuk menyimpan data transfer ke dalam tabel transfers
func CreateTransaction(db *sql.DB, transfer *structs.Transfer) error {
	query := `INSERT INTO transfers (sender_account_id, receiver_account_id, amount, transfer_type, fee, status, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := db.QueryRow(query, transfer.SenderAccountID, transfer.ReceiverAccountID, transfer.Amount, transfer.TransferType, transfer.Fee, transfer.Status, transfer.CreatedAt).Scan(&transfer.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetSenderAccountId(db *sql.DB, senderAccountID int) (structs.Account, error) {
	var account structs.Account
	query := `SELECT id, user_id, account_number, balance, account_type, created_at FROM accounts WHERE id = $1`
	err := db.QueryRow(query, senderAccountID).Scan(&account.ID, &account.UserID, &account.AccountNumber, &account.Balance, &account.AccountType, &account.CreatedAt)
	if err != nil {
		return account, err
	}
	return account, nil
}
