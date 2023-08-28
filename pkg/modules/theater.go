package modules

import "gorm.io/gorm"

// Define Theater struct to hold director details.
type Theater struct {
	gorm.Model
	Name     string
	Location string
}
