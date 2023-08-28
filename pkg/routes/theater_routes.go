package routes

import (
	"net/http"

	"forum.golangbridge.org/cinemaGo/pkg/controllers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetUpTheaterRouters(db *gorm.DB, route *mux.Router) {
	route.HandleFunc("/theaters", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateTheater(w, r, db)
	}).Methods("POST")

	route.HandleFunc("/theaters", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllTheaters(w, r, db)
	}).Methods("GET")

	route.HandleFunc("/theaters/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetTheater(w, r, db)
	}).Methods("GET")

	route.HandleFunc("/theaters/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteTheater(w, r, db)

	}).Methods("DELETE")

	route.HandleFunc("/theaters/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateTheater(w, r, db)
	}).Methods("PUT")

}
