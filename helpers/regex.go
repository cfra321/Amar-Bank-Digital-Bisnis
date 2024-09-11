package helpers

import (
	"errors"
	"regexp"
)

// ValidateUsername checks if the username is between 6-20 characters and contains only alphabetic characters.
func ValidateUsername(username string) error {
	if len(username) < 6 || len(username) > 20 {
		return errors.New("username must be between 6-20 characters")
	}
	if match, _ := regexp.MatchString("^[a-zA-Z0-9]{6,20}$", username); !match {
		return errors.New("username must be contain only alphanumeric characters")
	}
	return nil
}

// ValidatePassword checks if the password is at least 6 characters long and alphanumeric.
func ValidatePassword(password string) error {
	// Panjang password minimal 6 karakter
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	// Memeriksa apakah password mengandung setidaknya satu huruf
	if match, _ := regexp.MatchString("[a-zA-Z]", password); !match {
		return errors.New("password must contain at least one letter")
	}

	// Memeriksa apakah password mengandung setidaknya satu angka
	if match, _ := regexp.MatchString("[0-9]", password); !match {
		return errors.New("password must contain at least one number")
	}

	// Memeriksa apakah password mengandung setidaknya satu karakter spesial dari daftar
	if match, _ := regexp.MatchString("[!@#$%^&*\\-_]", password); !match {
		return errors.New("password must contain at least one special character (!@#$%^&*-_)")
	}

	return nil
}

// ValidateEmail checks if the email domain is either @gmail.com or @yahoo.co.id.
func ValidateEmail(email string) error {
	if match, _ := regexp.MatchString("^[a-zA-Z0-9._%+-]+@(gmail\\.com|yahoo\\.co\\.id)$", email); !match {
		return errors.New("email must be a valid @gmail.com or @yahoo.co.id address")
	}
	return nil
}

// ValidatePhoneNumber checks if the phone number is numeric and between 10-13 characters long.
func ValidatePhoneNumber(phone string) error {
	if len(phone) < 10 || len(phone) > 16 {
		return errors.New("phone number must be between 10-16 characters")
	}
	if match, _ := regexp.MatchString("^[0-9]+$", phone); !match {
		return errors.New("phone number must be numeric")
	}
	return nil
}
