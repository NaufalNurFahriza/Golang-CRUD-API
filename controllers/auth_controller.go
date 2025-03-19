package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"golang_MySQL/models"
	"golang_MySQL/utils"
)

// Register handles user registration
func Register(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request models.RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Check if user already exists
		var existingUser models.User
		result := db.Where("email = ?", request.Email).First(&existingUser)
		if result.RowsAffected > 0 {
			utils.RespondWithError(w, http.StatusConflict, "User with this email already exists")
			return
		}

		// Create new user
		user := models.User{
			Name:     request.Name,
			Email:    request.Email,
			Password: request.Password,
		}

		result = db.Create(&user)
		if result.Error != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
			return
		}

		// Create response register
		response := map[string]interface{}{
			"message": "Register successful",
			"user":    user.ToUserResponse(),
		}

		utils.RespondWithJSON(w, http.StatusCreated, response)
	}
}

// Login handles user login
func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request models.LoginRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Find user by email
		var user models.User
		result := db.Where("email = ?", request.Email).First(&user)
		if result.Error != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		// Check password
		if !user.CheckPassword(request.Password) {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		// Generate token
		token, err := generateToken(user.ID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
			return
		}

		// Create response login
		response := map[string]interface{}{
			"message": "Login successful",
			"token":   token,
			"user":    user.ToUserResponse(),
		}

		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}

// generateToken generates a new JWT token for a user
func generateToken(userID uint) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Get secret key from environment or use a default for development
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-secret-key" // Default secret key for development
	}

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
