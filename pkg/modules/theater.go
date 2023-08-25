package modules

import "gorm.io/gorm"

// Define Theater struct to hold director details.
type Theaters struct {
	gorm.Model
	Name     string
	Location string
}
