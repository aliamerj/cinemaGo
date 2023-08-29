package routes

import (
	"net/http"

	"forum.golangbridge.org/cinemaGo/api/middleware"
	"forum.golangbridge.org/cinemaGo/pkg/controllers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetUpTicketRoutes(db *gorm.DB, route *mux.Router) {
	s := route.PathPrefix("/").Subrouter()
	s.Use(middleware.AuthenticateMiddleware)

	s.HandleFunc("/booking", func(w http.ResponseWriter, r *http.Request) {
		controllers.BookingNewTicket(w, r, db)
	}).Methods("POST")
	s.HandleFunc("/tickets", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetMyTickets(w, r, db)
	}).Methods("GET")

}
