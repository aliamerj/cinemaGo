package modules

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Title           string
	Genre           string
	ReleaseDate     string
	DurationMinutes int
}
