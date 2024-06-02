package main

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
)

type SendMessageEvent struct {
	Message string `json:"message"`
	Room    string `json:"room"`
	From    string `json:"from"`
}

// NewMessageEvent is returned when responding to send_message
type NewMessageEvent struct {
	SendMessageEvent
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

	// Broadcast the new event to all clients in the room
	for client := range client.manager.clients {
		client.egress <- newEvent
	}
	return nil
}
