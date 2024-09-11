package repository

import (
	"database/sql"
	"errors"

	"github.com/cfra321/Amar-Bank-Digital-Bisnis/structs"
)

// GetUserByUsername fetches a user by their username
func GetUserByUsername(db *sql.DB, username string) (structs.User, error) {
	var user structs.User
	query := "SELECT id, username, password, phonenumber, is_active FROM users WHERE username = $1"

	err := db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.PhoneNumber,
		&user.IsActive,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

// UpdateUserStatus updates the is_active status of a user by ID
func UpdateUserStatus(db *sql.DB, id string, isActive bool) error {
	query := "UPDATE users SET is_active = $1 WHERE id = $2"
	_, err := db.Exec(query, isActive, id)
	return err
}

// CheckUserExists checks if a user with the same username or email already exists
func CheckUserExists(db *sql.DB, username, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1 OR email = $2)`

	err := db.QueryRow(query, username, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// RegisterUser registers a new user in the database
func RegisterUser(db *sql.DB, user structs.User) error {
	query := `
		INSERT INTO users (username, email, password, phonenumber, is_active, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())`

	_, err := db.Exec(query, user.Username, user.Email, user.Password, user.PhoneNumber, user.IsActive, user.CreatedBy)
	return err
}

// GetUsers retrieves all users

// GetUsers retrieves all users from the database
func GetAllUser(db *sql.DB) ([]structs.User, error) {
	var users []structs.User
	query := "SELECT id, username, email, phonenumber, is_active, created_at, created_by, modified_at, modified_by, activated_by FROM users"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user structs.User
		var modifiedBy sql.NullString
		var activatedBy sql.NullString

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PhoneNumber,
			&user.IsActive,
			&user.CreatedAt,
			&user.CreatedBy,
			&user.ModifiedAt,
			&modifiedBy,
			&activatedBy,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if modifiedBy.Valid {
			user.ModifiedBy = &modifiedBy.String
		} else {
			user.ModifiedBy = nil
		}
		user.ActivatedBy = &activatedBy.String

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByID(db *sql.DB, id string) (structs.User, error) {
	var user structs.User
	query := "SELECT id, username, email, phonenumber, is_active, created_at, created_by, modified_at, modified_by, activated_by FROM users WHERE id = $1"

	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&user.IsActive,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.ModifiedAt,
		&user.ModifiedBy,
		&user.ActivatedBy,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

// DeleteUser deletes a user from the database
func DeleteUser(db *sql.DB, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Jika tidak ada baris yang terpengaruh, artinya ID tidak ditemukan
	if rowsAffected == 0 {
		return errors.New("id not found")
	}

	return nil
}
