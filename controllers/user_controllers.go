package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/cfra321/Amar-Bank-Digital-Bisnis/database"
	"github.com/cfra321/Amar-Bank-Digital-Bisnis/helpers"
	"github.com/cfra321/Amar-Bank-Digital-Bisnis/repository"
	"github.com/cfra321/Amar-Bank-Digital-Bisnis/structs"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

// Login handles user login and returns a JWT token
func Login(c *gin.Context) {

	var loginData structs.LoginData

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user details from repository
	user, err := repository.GetUserByUsername(database.DbConnection, loginData.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User account is not active"})
		return
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// UpdateUserStatus updates the is_active status of a user
func UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")
	var statusUpdate struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update status via repository
	err := repository.UpdateUserStatus(database.DbConnection, id, statusUpdate.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})
}

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var registerData struct {
		Username    string `json:"username" binding:"required"`
		Email       string `json:"email" binding:"required"`
		Password    string `json:"password" binding:"required"`
		PhoneNumber string `json:"phonenumber" binding:"required"`
	}

	// Validate input from request
	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate username
	if err := helpers.ValidateUsername(registerData.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate password
	if err := helpers.ValidatePassword(registerData.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate email
	if err := helpers.ValidateEmail(registerData.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate phone number
	if err := helpers.ValidatePhoneNumber(registerData.PhoneNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username or email already exists
	exists, err := repository.CheckUserExists(database.DbConnection, registerData.Username, registerData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	// Create new user
	newUser := structs.User{
		Username:    registerData.Username,
		Email:       registerData.Email,
		Password:    string(hashedPassword),
		PhoneNumber: registerData.PhoneNumber,
		IsActive:    false,
		CreatedBy:   "system",
	}

	// Save new user to database
	err = repository.RegisterUser(database.DbConnection, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// GenerateToken creates a JWT token
func GenerateToken(username string) (string, error) {
	claims := &helpers.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	return token.SignedString(secretKey)
}

func GetAllUsers(c *gin.Context) {
	users, err := repository.GetAllUser(database.DbConnection)
	if err != nil {
		// Log the error for debugging purposes
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := repository.GetUserByID(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := repository.DeleteUser(database.DbConnection, id)
	if err != nil {
		if err.Error() == "id not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User ID not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
