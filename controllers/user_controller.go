package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"go-restapi-crud/models"
	"go-restapi-crud/utils"
)

// CreateUser creates a new user
func CreateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Check if user already exists
		var existingUser models.User
		result := db.Where("email = ?", user.Email).First(&existingUser)
		if result.RowsAffected > 0 {
			utils.RespondWithError(w, http.StatusConflict, "User with this email already exists")
			return
		}

		// Use goroutine for database operation
		var wg sync.WaitGroup
		var dbErr error
		var createdUser models.User

		wg.Add(1)
		go func() {
			defer wg.Done()
			result := db.Create(&user)
			if result.Error != nil {
				dbErr = result.Error
				return
			}
			createdUser = user
		}()
		wg.Wait()

		if dbErr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, createdUser.ToUserResponse())
	}
}

// GetAllUsers returns all users
func GetAllUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []models.User
		var userResponses []models.UserResponse

		// Use goroutine for database operation
		var wg sync.WaitGroup
		var dbErr error

		wg.Add(1)
		go func() {
			defer wg.Done()
			result := db.Find(&users)
			if result.Error != nil {
				dbErr = result.Error
				return
			}

			userResponses = make([]models.UserResponse, len(users))
			for i, user := range users {
				userResponses[i] = user.ToUserResponse()
			}
		}()
		wg.Wait()

		if dbErr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving users")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, userResponses)
	}
}

// GetUser returns a single user by ID
func GetUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		var user models.User
		var wg sync.WaitGroup
		var dbErr error

		wg.Add(1)
		go func() {
			defer wg.Done()
			result := db.First(&user, id)
			if result.Error != nil {
				dbErr = result.Error
				return
			}
		}()
		wg.Wait()

		if dbErr != nil {
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, user.ToUserResponse())
	}
}

// UpdateUser updates a user
func UpdateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		var updateData map[string]interface{}
		err = json.NewDecoder(r.Body).Decode(&updateData)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Remove password from update data if it's empty
		if password, ok := updateData["password"]; ok {
			if password == "" {
				delete(updateData, "password")
			}
		}

		var user models.User
		var wg sync.WaitGroup
		var dbErr error

		wg.Add(1)
		go func() {
			defer wg.Done()
			// First check if user exists
			result := db.First(&user, id)
			if result.Error != nil {
				dbErr = result.Error
				return
			}

			// If email is being updated, check for uniqueness
			if email, ok := updateData["email"]; ok {
				var existingUser models.User
				result = db.Where("email = ? AND id != ?", email, id).First(&existingUser)
				if result.RowsAffected > 0 {
					dbErr = gorm.ErrDuplicatedKey
					return
				}
			}

			// Update user
			result = db.Model(&user).Updates(updateData)
			if result.Error != nil {
				dbErr = result.Error
				return
			}

			// Get updated user
			db.First(&user, id)
		}()
		wg.Wait()

		if dbErr != nil {
			if dbErr == gorm.ErrDuplicatedKey {
				utils.RespondWithError(w, http.StatusConflict, "Email already in use")
				return
			}
			if dbErr == gorm.ErrRecordNotFound {
				utils.RespondWithError(w, http.StatusNotFound, "User not found")
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, "Error updating user")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, user.ToUserResponse())
	}
}

// DeleteUser deletes a user
func DeleteUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		var user models.User
		var wg sync.WaitGroup
		var dbErr error

		wg.Add(1)
		go func() {
			defer wg.Done()
			// First check if user exists
			result := db.First(&user, id)
			if result.Error != nil {
				dbErr = result.Error
				return
			}

			// Delete user
			result = db.Delete(&user)
			if result.Error != nil {
				dbErr = result.Error
				return
			}
		}()
		wg.Wait()

		if dbErr != nil {
			if dbErr == gorm.ErrRecordNotFound {
				utils.RespondWithError(w, http.StatusNotFound, "User not found")
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting user")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
	}
}
