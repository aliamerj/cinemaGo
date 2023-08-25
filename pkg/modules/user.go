package modules

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string
	Email    string
}
