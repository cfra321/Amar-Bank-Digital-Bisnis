package structs

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID          int        `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Password    string     `json:"password"`
	PhoneNumber string     `json:"phonenumber"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   string     `json:"created_by"`
	ModifiedAt  *time.Time `json:"modified_at,omitempty"`
	ModifiedBy  *string    `json:"modified_by,omitempty"`
	ActivatedBy *string    `json:"activated_by,omitempty"`
}

type LoginData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Account struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	AccountNumber string    `json:"account_number"`
	Balance       float64   `json:"balance"`
	AccountType   string    `json:"account_type"` // 'internal' or 'external'
	CreatedAt     time.Time `json:"created_at"`
}

// Transfer represents a transfer transaction in the system
type Transfer struct {
	ID                int       `json:"id"`
	SenderAccountID   int       `json:"sender_account_id"`
	ReceiverAccountID int       `json:"receiver_account_id"`
	Amount            float64   `json:"amount"`
	TransferType      string    `json:"transfer_type"` // 'overbook' or 'bifast'
	Fee               float64   `json:"fee"`
	Status            string    `json:"status"` // 'pending', 'completed', or 'failed'
	CreatedAt         time.Time `json:"created_at"`
	CompletedAt       time.Time `json:"completed_at,omitempty"`
}

// TransactionLog struct untuk tabel transaction_logs
type TransactionLog struct {
	ID         int       `json:"id"`
	TransferID int       `json:"transfer_id"`
	LogMessage string    `json:"log_message"`
	CreatedAt  time.Time `json:"created_at"`
}
