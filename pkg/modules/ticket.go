package modules

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model
	ShowtimeID uint
	UserID     uint
	SeatNumber int
}
