package models

import (
	"encoding/json"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type SendMessageEvent struct {
	Message string `json:"message"`
}

// NewMessageEvent is returned when responding to send_message
type NewMessageEvent struct {
	SendMessageEvent
	From string    `json:"from"`
	Sent time.Time `json:"sent"`
}
