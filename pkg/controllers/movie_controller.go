package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"forum.golangbridge.org/cinemaGo/pkg/modules"
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
