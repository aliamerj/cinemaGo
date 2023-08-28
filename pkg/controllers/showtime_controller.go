package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"forum.golangbridge.org/cinemaGo/pkg/modules"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func CreateShowtime(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	var showtime modules.Showtime

	// Decode JSON payload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&showtime); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the Showtime
	if !isValidShowtime(showtime) {
		http.Error(w, "Invalid showtime data", http.StatusBadRequest)
		return
	}

	// Using Goroutines for concurrent database checks.
	// Be cautious with this approach if your database driver/connection pool can't handle it.
	var wg sync.WaitGroup
	var movie modules.Movie
	var theater modules.Theater
	var movieError, theaterError error

	wg.Add(2)

	go func() {
		defer wg.Done()
		resMovie := db.Find(&movie, showtime.MovieID)
		if resMovie.Error != nil || resMovie.RowsAffected == 0 {
			movieError = errors.New("Movie Not Found")
		}
	}()

	go func() {
		defer wg.Done()
		resTheater := db.Find(&theater, showtime.TheaterID)
		if resTheater.Error != nil || resTheater.RowsAffected == 0 {
			theaterError = errors.New("Theater Not Found")
		}
	}()

	wg.Wait()

	if movieError != nil || theaterError != nil {
		http.Error(w, "Movie or Theater Not Found", http.StatusBadRequest)
		return
	}

	// Create the Showtime
	result := db.Create(&showtime)
	if result.Error != nil {
		log.Printf("Error while creating showtime: %v", result.Error)
		http.Error(w, "Could not create showtime", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&showtime); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func GetAllShowtimes(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	w.Header().Set("Content-Type", "application/json")
	var showtime []modules.Showtime

	// Perform the database query; populate the theater slice
	result := db.Find(&showtime)

	// Check for errors
	if result.Error != nil {
		// It might be useful to log the error as well, for debugging purposes.
		log.Printf("Error while getting movies: %v", result.Error)
		http.Error(w, "Could not retrieve movies", http.StatusInternalServerError)
		return
	}

	// Encode the slice into JSON and write it to the response
	if err := json.NewEncoder(w).Encode(showtime); err != nil {
		log.Printf("Error while encoding movies to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func GetShowtime(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Set content type for the response
	w.Header().Set("Content-Type", "application/json")

	// Initialize a showtime variable
	var showtime modules.Showtime

	// Retrieve 'id' from the URL
	id := mux.Vars(r)["id"]

	// Fetch the showtime from the database
	result := db.Find(&showtime, id)

	// Check for errors during the database query
	if result.Error != nil {
		// Log the error for debugging purposes
		log.Printf("Error while getting showtime: %v", result.Error)
		http.Error(w, "Could not retrieve showtime", http.StatusInternalServerError)
		return
	}

	// Check if the record exists
	if result.RowsAffected == 0 {
		http.Error(w, "showtime not found", http.StatusNotFound)
		return
	}

	// Send the showtime details in JSON format
	if err := json.NewEncoder(w).Encode(showtime); err != nil {
		log.Printf("Error while encoding showtime to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func DeleteShowtime(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Initialize a showtime variable
	var showtime modules.Showtime

	// Retrieve 'id' from the URL
	id := mux.Vars(r)["id"]

	// Delete the showtime from the database
	result := db.Delete(&showtime, id)

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

func UpdateShowtime(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	// Retrieve 'id' from the URL
	id := mux.Vars(r)["id"]

	// Initialize a showtime variable
	var showtime modules.Showtime

	// Decode the incoming JSON into a map
	updateData := make(map[string]interface{})

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the showtime in the database
	result := db.Model(&showtime).Where("id = ?", id).Updates(updateData)

	// Check for errors and whether any rows were updated
	if result.Error != nil {
		log.Printf("Error while updating the showtime: %v", result.Error)
		http.Error(w, "Could not update the showtime", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "showtime not found", http.StatusNotFound)
		return
	}

	// Return the updated showtime
	response := map[string]interface{}{"message": "showtime info get updated successfully", "data": &showtime}
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		log.Printf("Error while encoding showtime to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
func isValidShowtime(showtime modules.Showtime) bool {
	// Checking if MovieID and TheaterID are not zero
	if showtime.MovieID == 0 || showtime.TheaterID == 0 {
		return false
	}

	// Checking if the StartTime is in the future
	if showtime.StartTime.Before(time.Now()) {
		return false
	}

	return true
}
