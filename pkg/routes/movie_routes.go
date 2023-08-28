package routes

import (
	"net/http"

	"forum.golangbridge.org/cinemaGo/pkg/controllers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetUpMovieRoutes(db *gorm.DB, route *mux.Router) {
	route.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateMovie(w, r, db)
	}).Methods("POST")

	route.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllMovies(w, r, db)
	}).Methods("GET")

	route.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetMovie(w, r, db)
	}).Methods("GET")

	route.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteMovie(w, r, db)

	}).Methods("DELETE")

	route.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateMovie(w, r, db)
	}).Methods("PUT")

}
