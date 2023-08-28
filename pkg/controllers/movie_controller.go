package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"forum.golangbridge.org/cinemaGo/pkg/modules"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// CreateMovie is an exported HTTP handler function for creating a new Movie.
// It expects a JSON object in the HTTP request body that corresponds to the Movie structure.
// Upon successful decoding, it echoes the same JSON object back to the client.
func CreateMovie(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Always close the request body, even if an error occurs.
	defer r.Body.Close()

	// Declare a variable of type Movie from the 'modules' package to hold the incoming request data.
	var movie modules.Movie

	// Create a new JSON decoder and decode the incoming request body into the 'movie' variable.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&movie); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// Validate the movie before creating it
	if !isValidMovie(movie) {
		http.Error(w, "Invalid movie data", http.StatusBadRequest)
		return
	}

	// Here you would typically add code to process the movie object, for example, store it in a database.
	result := db.Create(&movie)
	// Check for errors
	if result.Error != nil {
		// It might be useful to log the error as well, for debugging purposes.
		log.Printf("Error while creating movie: %v", result.Error)
		http.Error(w, "Could not create movies", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type for the response to indicate we are returning JSON.
	w.Header().Set("Content-Type", "application/json")

	// Create a new JSON encoder and write the 'movie' variable to the response.
	// Note that we are now checking for errors when encoding.
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(movie); err != nil {
		// If encoding fails, set the status code to 'Internal Server Error' (500) and return.
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func GetAllMovies(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	w.Header().Set("Content-Type", "application/json")
	var movies []modules.Movie

	// Perform the database query; populate the movies slice
	result := db.Find(&movies)

	// Check for errors
	if result.Error != nil {
		// It might be useful to log the error as well, for debugging purposes.
		log.Printf("Error while getting movies: %v", result.Error)
		http.Error(w, "Could not retrieve movies", http.StatusInternalServerError)
		return
	}

	// Encode the slice into JSON and write it to the response
	if err := json.NewEncoder(w).Encode(movies); err != nil {
		log.Printf("Error while encoding movies to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func GetMovie(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Set content type for the response
	w.Header().Set("Content-Type", "application/json")

	// Initialize a Movie variable
	var movie modules.Movie

	// Retrieve 'id' from the URL
	id := mux.Vars(r)["id"]

	// Fetch the movie from the database
	result := db.Find(&movie, id)

	// Check for errors during the database query
	if result.Error != nil {
		// Log the error for debugging purposes
		log.Printf("Error while getting movie: %v", result.Error)
		http.Error(w, "Could not retrieve movie", http.StatusInternalServerError)
		return
	}

	// Check if the record exists
	if result.RowsAffected == 0 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	// Send the movie details in JSON format
	if err := json.NewEncoder(w).Encode(movie); err != nil {
		log.Printf("Error while encoding movie to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func DeleteMovie(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Initialize a Movie variable
	var movie modules.Movie

	// Retrieve 'id' from the URL
	id := mux.Vars(r)["id"]

	// Delete the movie from the database
	result := db.Delete(&movie, id)

	// Check for errors during the database operation
	if result.Error != nil {
		// Log the error for debugging
		log.Printf("Error while deleting the movie: %v", result.Error)
		http.Error(w, "Could not delete the movie", http.StatusInternalServerError)
		return
	}

	// Check if a record was actually deleted
	if result.RowsAffected == 0 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	// Send a success message
	response := map[string]string{"message": "The movie has been deleted"}

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		log.Printf("Error while encoding response to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func UpdateMovie(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	// Retrieve 'id' from the URL
	id := mux.Vars(r)["id"]

	// Initialize a Movie variable
	var movie modules.Movie

	// Decode the incoming JSON into a map
	updateData := make(map[string]interface{})

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the movie in the database
	result := db.Model(&movie).Where("id = ?", id).Updates(updateData)

	// Check for errors and whether any rows were updated
	if result.Error != nil {
		log.Printf("Error while updating the movie: %v", result.Error)
		http.Error(w, "Could not update the movie", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	// Return the updated movie
	response := map[string]interface{}{"message": "Movie info get updated successfully", "data": &movie}
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		log.Printf("Error while encoding movie to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
func isValidMovie(movie modules.Movie) bool {
	// Your validation logic here
	// For example, check that required fields are present:
	if movie.Title == "" || movie.Genre == "" || movie.DurationMinutes == 0 {
		return false
	}
	// More validation logic
	return true
}
