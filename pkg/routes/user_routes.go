package routes

import (
	"net/http"

	"forum.golangbridge.org/cinemaGo/pkg/controllers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetUpUserRoutes(db *gorm.DB, route *mux.Router) {
	route.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		controllers.Register(w, r, db)

	}).Methods("POST")
	route.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controllers.Login(w, r, db)
	}).Methods("POST")

}
