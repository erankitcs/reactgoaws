package chatws

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, client *Client) error

const (
	EventTypeSendMessage = "send_message"
	EventTypeNewMessage  = "new_message"
	EventTypeUserLeft    = "user_left"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

// NewMessageEvent is returned when responding to send_message
type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

// UserLeftEvent is returned when a client leaving or disconnected
type UserLeftEvent struct {
	User string    `json:"user"`
	Sent time.Time `json:"sent"`
}

func SendMessageHandler(event Event, client *Client) error {
	// Decode the event payload
	var sendMessageEvent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &sendMessageEvent); err != nil {
		return fmt.Errorf("bad event payload: %w", err)
	}

	var newMessageEvent NewMessageEvent
	newMessageEvent.Sent = time.Now()
	newMessageEvent.SendMessageEvent = sendMessageEvent

	//Marshal the new message event
	newMessageEventPayload, err := json.Marshal(newMessageEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal new message event: %w", err)
	}

	// Create a new event with the new message event payload
	newEvent := Event{
		Type:    EventTypeNewMessage,
		Payload: newMessageEventPayload,
	}

	// Broadcast the new event to all clients connected for same movied Id
	for c := range client.manager.clients {
		if client.movieID != c.movieID || client.id == c.id {
			continue
		}
		c.egress <- newEvent
	}
	return nil
}

// Will be called by client removal event
func UserLeftHandler(client *Client) error {
	// Decode the event payload
	var userLeftEvent UserLeftEvent
	userLeftEvent.User = client.id
	userLeftEvent.Sent = time.Now()

	//Marshal the user left event
	userLeftEventPayload, err := json.Marshal(userLeftEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal user left event: %w", err)
	}

	// Create a new event with the new message event payload
	newEvent := Event{
		Type:    EventTypeUserLeft,
		Payload: userLeftEventPayload,
	}

	// Broadcast the new event to all clients connected for same movied Id
	for c := range client.manager.clients {
		if client.movieID != c.movieID || client.id == c.id {
			continue
		}
		c.egress <- newEvent
	}
	return nil
}
