package modules

import "gorm.io/gorm"

type Seat struct {
	gorm.Model
	UserID     uint
	ShowtimeID uint
	SeatNumber int
}
