package models

import "time"

type Message struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Text      string    `json:"text"`
	Date      time.Time `json:"date"`
}
