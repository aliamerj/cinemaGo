package routes

import (
	"net/http"

	"forum.golangbridge.org/cinemaGo/pkg/controllers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupMoviesRoute(db *gorm.DB, route *mux.Router) {
	route.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateMovie(w, r, db)
	}).Methods("POST")
	route.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllMovies(w, r, db)
	}).Methods("GET")

}
