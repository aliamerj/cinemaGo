package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"forum.golangbridge.org/cinemaGo/pkg/modules"
	"gorm.io/gorm"
)

func BookingNewTicket(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	w.Header().Set("Content-Type", "application/json")

	userIdStr := r.Context().Value("userID").(string)
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	var ticketRequest struct {
		SeatNumber int  `json:"seatNumber"`
		ShowID     uint `json:"showId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&ticketRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Transaction rolled back: %v", r)
		}
	}()

	if tx.Error != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	seat := modules.Seat{
		UserID:     uint(userId),
		ShowtimeID: ticketRequest.ShowID,
		SeatNumber: ticketRequest.SeatNumber,
	}

	if err := bookSeat(tx, &seat); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newTicket := modules.Ticket{
		UserID:     uint(userId),
		ShowtimeID: seat.ShowtimeID,
		SeatID:     seat.ID,
	}

	if err := createTicket(tx, &newTicket); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx.Commit()

	response := map[string]interface{}{
		"message": "Seat has been booked successfully",
		"data":    &newTicket,
	}

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		log.Printf("Error while encoding to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func GetMyTickets(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Set the Content-Type for the response to indicate we are returning JSON.
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("userID").(string)

	var tickets []modules.Ticket

	result := db.Find(&tickets).Where("user_id = ?", userID)
	// Check for errors
	if result.Error != nil {
		// It might be useful to log the error as well, for debugging purposes.
		log.Printf("Error while getting tickets: %v", result.Error)
		http.Error(w, "Could not retrieve tickets", http.StatusInternalServerError)
		return
	}

	// Encode the slice into JSON and write it to the response
	if err := json.NewEncoder(w).Encode(tickets); err != nil {
		log.Printf("Error while encoding movies to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}

func bookSeat(tx *gorm.DB, seat *modules.Seat) error {
	var showtime modules.Showtime
	if seat.SeatNumber < 1 || seat.SeatNumber > 100 {
		return fmt.Errorf("There is no seat with this number")

	}

	resShow := tx.Find(&showtime, seat.ShowtimeID)
	if resShow.Error != nil || resShow.RowsAffected == 0 {
		return errors.New("No show Found")
	}

	if tx.Where(&modules.Seat{SeatNumber: seat.SeatNumber, ShowtimeID: seat.ShowtimeID}).First(&seat).RowsAffected != 0 {
		return errors.New("Seat already booked")
	}

	if err := tx.Create(&seat).Error; err != nil {
		return fmt.Errorf("Failed to book seat: %v", err)
	}
	return nil
}

func createTicket(tx *gorm.DB, ticket *modules.Ticket) error {
	if err := tx.Create(&ticket).Error; err != nil {
		return fmt.Errorf("Unable to book ticket: %v", err)
	}
	return nil
}
