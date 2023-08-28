package modules

import (
	"time"

	"gorm.io/gorm"
)

// Define Movie struct to hold movie details.
type Showtime struct {
	gorm.Model
	MovieID   uint
	TheaterID uint
	StartTime time.Time
}
