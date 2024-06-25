package chatws

import (
	"backend/internal/models"
	"encoding/json"
	"fmt"
	"time"
)

type EventHandler func(event models.Event, client *Client) error

const (
	EventTypeSendMessage = "send_message"
	EventTypeNewMessage  = "new_message"
	EventTypeUserLeft    = "user_left"
	EventTypeUserJoined  = "user_joined"
)

func SendMessageHandler(event models.Event, client *Client) error {
	// Decode the event payload
	var sendMessageEvent models.SendMessageEvent
	if err := json.Unmarshal(event.Payload, &sendMessageEvent); err != nil {
		return fmt.Errorf("bad event payload: %w", err)
	}

	var newMessageEvent models.NewMessageEvent
	newMessageEvent.Sent = time.Now()
	newMessageEvent.From = client.userName
	newMessageEvent.SendMessageEvent = sendMessageEvent

	//Marshal the new message event
	newMessageEventPayload, err := json.Marshal(newMessageEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal new message event: %w", err)
	}

	// Create a new event with the new message event payload
	newEvent := models.Event{
		Type:    EventTypeNewMessage,
		Payload: newMessageEventPayload,
	}

	// Broadcast the new event to all clients connected for same movied Id
	for c := range client.manager.clients {
		if client.movieID != c.movieID {
			continue
		}
		c.egress <- newEvent
	}

	/// Movie Chat persistance.
	var movieChat models.MovieChat
	movieChat.MovieID = client.movieID
	movieChat.UserID = client.userID
	movieChat.ChatText = newMessageEvent.Message
	movieChat.CreatedAt = newMessageEvent.Sent

	/// Save msg into database
	err = client.manager.db.InsertMovieChat(movieChat)
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	return nil
}

// Will be called by client removal event
func UserLeftHandler(client *Client) error {
	var sendUserLeftEvent models.SendMessageEvent
	sendUserLeftEvent.Message = "User left"
	var newUserLeftEvent models.NewMessageEvent
	newUserLeftEvent.SendMessageEvent = sendUserLeftEvent
	newUserLeftEvent.Sent = time.Now()
	newUserLeftEvent.From = client.userName

	//Marshal the user left event
	userLeftEventPayload, err := json.Marshal(newUserLeftEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal user left event: %w", err)
	}

	// Create a new event with the new message event payload
	newEvent := models.Event{
		Type:    EventTypeUserLeft,
		Payload: userLeftEventPayload,
	}

	// Broadcast the new event to all clients connected for same movied Id
	for c := range client.manager.clients {
		if client.movieID != c.movieID {
			continue
		}
		c.egress <- newEvent
	}
	return nil
}

// Will be called by client add event
func UserJoinedHandler(client *Client) error {
	var sendUserJoinedEvent models.SendMessageEvent
	sendUserJoinedEvent.Message = "User joined"
	var newUserJoinedEvent models.NewMessageEvent
	newUserJoinedEvent.SendMessageEvent = sendUserJoinedEvent
	newUserJoinedEvent.Sent = time.Now()
	newUserJoinedEvent.From = client.userName
	//Marshal the user left event
	userJoinedEventPayload, err := json.Marshal(newUserJoinedEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal user left event: %w", err)
	}

	// Create a new event with the new message event payload
	newEvent := models.Event{
		Type:    EventTypeUserJoined,
		Payload: userJoinedEventPayload,
	}

	// Broadcast the new event to all clients connected for same movied Id
	for c := range client.manager.clients {
		if client.movieID != c.movieID {
			continue
		}
		c.egress <- newEvent
	}
	return nil
}
