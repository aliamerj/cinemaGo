package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"forum.golangbridge.org/cinemaGo/pkg/modules"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func CreateTheater(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	var theater modules.Theater

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&theater); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// Validate the theater before creating it
	if !isValidTheater(theater) {
		http.Error(w, "Invalid theater data", http.StatusBadRequest)
		return
	}

	result := db.Create(&theater)
	if result.Error != nil {
		// It might be useful to log the error as well, for debugging purposes.
		log.Printf("Error while creating theater: %v", result.Error)
		http.Error(w, "Could not create theater", http.StatusInternalServerError)
		return
	}

	// Create a new JSON encoder and write the 'theater' variable to the response.
	// Note that we are now checking for errors when encoding.
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&theater); err != nil {
		// If encoding fails, set the status code to 'Internal Server Error' (500) and return.
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}

}

func GetAllTheaters(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	w.Header().Set("Content-Type", "application/json")
	var theaters []modules.Theater

	// Perform the database query; populate the theater slice
	result := db.Find(&theaters)

	// Check for errors
	if result.Error != nil {
		// It might be useful to log the error as well, for debugging purposes.
		log.Printf("Error while getting movies: %v", result.Error)
		http.Error(w, "Could not retrieve movies", http.StatusInternalServerError)
		return
	}

	// Encode the slice into JSON and write it to the response
	if err := json.NewEncoder(w).Encode(theaters); err != nil {
		log.Printf("Error while encoding movies to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func GetTheater(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Set content type for the response
	w.Header().Set("Content-Type", "application/json")

	// Initialize a theater variable
	var theater modules.Theater

	// Retrieve 'id' from the URL
	id := mux.Vars(r)["id"]

	// Fetch the theater from the database
	result := db.Find(&theater, id)

	// Check for errors during the database query
	if result.Error != nil {
		// Log the error for debugging purposes
		log.Printf("Error while getting theater: %v", result.Error)
		http.Error(w, "Could not retrieve theater", http.StatusInternalServerError)
		return
	}

	// Check if the record exists
	if result.RowsAffected == 0 {
		http.Error(w, "theater not found", http.StatusNotFound)
		return
	}

	// Send the theater details in JSON format
	if err := json.NewEncoder(w).Encode(theater); err != nil {
		log.Printf("Error while encoding theater to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func DeleteTheater(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Initialize a theater variable
	var theater modules.Theater

	// Retrieve 'id' from the URL
	id := mux.Vars(r)["id"]

	// Delete the theater from the database
	result := db.Delete(&theater, id)

	// Check for errors during the database operation
	if result.Error != nil {
		// Log the error for debugging
		log.Printf("Error while deleting the theater: %v", result.Error)
		http.Error(w, "Could not delete the theater", http.StatusInternalServerError)
		return
	}

	// Check if a record was actually deleted
	if result.RowsAffected == 0 {
		http.Error(w, "theater not found", http.StatusNotFound)
		return
	}

	// Send a success message
	response := map[string]string{"message": "The theater has been deleted"}

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		log.Printf("Error while encoding response to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func UpdateTheater(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	// Retrieve 'id' from the URL
	id := mux.Vars(r)["id"]

	// Initialize a Theater variable
	var theater modules.Theater

	// Decode the incoming JSON into a map
	updateData := make(map[string]interface{})

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the theater in the database
	result := db.Model(&theater).Where("id = ?", id).Updates(updateData)

	// Check for errors and whether any rows were updated
	if result.Error != nil {
		log.Printf("Error while updating the theater: %v", result.Error)
		http.Error(w, "Could not update the theater", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "theater not found", http.StatusNotFound)
		return
	}

	// Return the updated theater
	response := map[string]interface{}{"message": "theater info get updated successfully", "data": &theater}
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		log.Printf("Error while encoding theater to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
func isValidTheater(theater modules.Theater) bool {
	// Your validation logic here
	// For example, check that required fields are present:
	if theater.Name == "" || theater.Location == "" {
		return false
	}
	// More validation logic
	return true
}
