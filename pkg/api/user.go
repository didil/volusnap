package api

import "time"

// User model
type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string `gorm:"type:varchar(100);unique_index"`
	Password  string
}
