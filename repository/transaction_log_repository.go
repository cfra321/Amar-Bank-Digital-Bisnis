package repository

import (
	"database/sql"

	"github.com/cfra321/Amar-Bank-Digital-Bisnis/structs"
)

// CreateTransactionLog untuk menyimpan log transaksi ke dalam tabel transaction_logs
func CreateTransactionLog(db *sql.DB, log *structs.TransactionLog) error {
	query := `INSERT INTO transaction_logs (transfer_id, log_message, created_at) 
			  VALUES ($1, $2, $3)`
	_, err := db.Exec(query, log.TransferID, log.LogMessage, log.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetAllTransactionLogs untuk mendapatkan semua log transaksi dari tabel transaction_logs
func GetAllTransactionLogs(db *sql.DB) ([]structs.TransactionLog, error) {
	query := `SELECT id, transfer_id, log_message, created_at FROM transaction_logs`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactionLogs []structs.TransactionLog
	for rows.Next() {
		var log structs.TransactionLog
		if err := rows.Scan(&log.ID, &log.TransferID, &log.LogMessage, &log.CreatedAt); err != nil {
			return nil, err
		}
		transactionLogs = append(transactionLogs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactionLogs, nil
}
