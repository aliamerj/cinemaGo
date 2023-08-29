package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"forum.golangbridge.org/cinemaGo/pkg/modules"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// isValidPassword validates if the password meets your requirements
func isValidPassword(password string) bool {
	// Add your validation logic here. For example:
	return len(password) >= 8
}

func Register(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Set the Content-Type for the response to indicate we are returning JSON.
	w.Header().Set("Content-Type", "application/json")

	// Close the request body
	defer r.Body.Close()

	// Decode incoming request to newUser
	var newUser modules.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate the password
	if !isValidPassword(newUser.Password) {
		http.Error(w, "Password does not meet the requirements", http.StatusBadRequest)
		return
	}

	// Hash the password before saving the user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)

	// Check if email exists and save the newUser
	if err := db.Create(&newUser).Error; err != nil {
		// Log the error internally
		log.Printf("Could not register user: %v", err)

		// Check for duplicate email
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			http.Error(w, "Email already registered", http.StatusConflict)
		} else {
			http.Error(w, "Could not register user", http.StatusInternalServerError)
		}
		return
	}

	// Create response struct, maybe without sensitive information
	response := map[string]string{
		"message": "User successfully registered",
		"userId":  fmt.Sprintf("%d", newUser.ID),
	}

	// Encode the response and write it to w
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func Login(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Close the request body
	defer r.Body.Close()

	// Decode the incoming request
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user modules.User
	// Find the user
	if err := db.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Check the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Create a response
	response := map[string]string{
		"message": "Logged in successfully",
		"email":   user.Email,
	}

	// Set the Authorization header for the response
	w.Header().Set("Authorization", "Bearer "+token)

	// Respond
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func generateJWT(user modules.User) (string, error) {
	// Define the expiration time of the token
	// Here, the token will expire in one hour
	expireToken := time.Now().Add(time.Hour * 1).Unix()

	// Create the Claims
	claims := jwt.StandardClaims{
		ExpiresAt: expireToken,
		Issuer:    "CinemaGo",
		Subject:   fmt.Sprintf("%d", user.ID),
	}

	// Create the token using your structure that implements jwt.Claims
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte("mySecretKey"))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
