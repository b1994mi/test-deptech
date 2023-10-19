package sqlmodel

import (
	"time"
)

type Admin struct {
	ID          int       `json:"id"`
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	Password    string    `json:"password,omitempty"`
}
