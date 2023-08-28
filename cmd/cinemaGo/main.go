package main

import (
	"log"
	"net/http"

	"forum.golangbridge.org/cinemaGo/internal/database"
	"forum.golangbridge.org/cinemaGo/pkg/modules"
	"forum.golangbridge.org/cinemaGo/pkg/routes"
	"github.com/gorilla/mux"
)

func main() {
	db, err := database.Initialization()
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)

	}
	// Automatically migrate multiple schemas
	err = db.AutoMigrate(&modules.Movie{}, &modules.Theater{}, &modules.User{}, &modules.Ticket{}, &modules.Showtime{})
	if err != nil {
		log.Fatalf("Could not migrate database: %v", err)
	}
	r := mux.NewRouter()
	routes.SetUpMovieRoutes(db, r)
	routes.SetUpTheaterRouters(db, r)
	routes.SetUpShowtimeRoutes(db, r)

	const port = ":8080"
	// Log the port number where the server will start.
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal("Failed to start server: ", err)
	}

}
