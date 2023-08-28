package routes

import (
	"net/http"

	"forum.golangbridge.org/cinemaGo/pkg/controllers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetUpShowtimeRoutes(db *gorm.DB, route *mux.Router) {
	route.HandleFunc("/showtimes", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateShowtime(w, r, db)
	}).Methods("POST")

	route.HandleFunc("/showtimes", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllShowtimes(w, r, db)
	}).Methods("GET")

	route.HandleFunc("/showtimes/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetShowtime(w, r, db)
	}).Methods("GET")

	route.HandleFunc("/showtimes/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteShowtime(w, r, db)

	}).Methods("DELETE")

	route.HandleFunc("/showtimes/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateShowtime(w, r, db)
	}).Methods("PUT")

}
