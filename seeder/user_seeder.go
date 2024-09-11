package seeder

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// SeedUsers seeds the database with initial user data
func SeedUsers(db *sql.DB) {
	// Data pengguna yang akan di-seed
	users := []struct {
		Username    string
		Password    string
		Email       string
		PhoneNumber string
		IsActive    bool
		CreatedAt   string
		CreatedBy   string
		ActivatedBy string
		ActivatedAt string
	}{
		{
			Username:    "adminbank",
			Password:    "admin123", // ini akan di-hash
			Email:       "admin@amarbank.co.id",
			PhoneNumber: "081234567890", // Include phone number here
			IsActive:    true,
			CreatedAt:   "2024-08-03T14:55:00Z",
			CreatedBy:   "system",
			ActivatedBy: "admin",
			ActivatedAt: "2024-08-03T15:00:00Z",
		},
	}

	for _, user := range users {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Tidak dapat meng-hash password: %v", err)
		}

		// Periksa apakah pengguna sudah ada
		var existingID int
		err = db.QueryRow(`SELECT id FROM users WHERE username = $1`, user.Username).Scan(&existingID)
		if err != nil && err != sql.ErrNoRows {
			log.Fatalf("Error memeriksa apakah pengguna ada: %v", err)
		}

		// Jika pengguna tidak ada, lakukan insert
		if err == sql.ErrNoRows {
			_, err = db.Exec(`
				INSERT INTO users (username, password, email, phonenumber, is_active, created_at, created_by, activated_by, activated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`, user.Username, hashedPassword, user.Email, user.PhoneNumber, user.IsActive, user.CreatedAt, user.CreatedBy, user.ActivatedBy, user.ActivatedAt)

			if err != nil {
				log.Fatalf("Tidak dapat memasukkan pengguna: %v", err)
			}
		}
	}

	log.Println("Pengguna berhasil di-seed.")
}
