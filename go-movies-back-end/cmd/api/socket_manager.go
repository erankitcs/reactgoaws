package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

type ChatManager struct {
	clients ClientList
	sync.RWMutex
	handlers map[string]EventHandler
}

func NewChatManager() *ChatManager {
	cm := &ChatManager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	// Setup event handlers
	cm.SetupEventHandlers()
	return cm
}

func (cm *ChatManager) SetupEventHandlers() {
	cm.handlers[EventTypeSendMessage] = SendMessageHandler
}

func (cm *ChatManager) routeEvent(event Event, c *Client) error {
	// Check if event type is available then excutre its handler.
	if handler, ok := cm.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			log.Printf("Error handling event: %v", err)
		}
		return nil
	} else {
		log.Printf("No handler for event type %s", event.Type)
		return errors.New("no handler for event type")
	}
}

func (cm *ChatManager) serveChat(w http.ResponseWriter, r *http.Request) {
	fmt.Println("connection recieved..")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//defer ws.Close()

	client := NewClient(ws, cm)
	cm.addClient(client)

	// Start client process
	go client.readMessages()
	go client.writeMessages()

}

func (cm *ChatManager) addClient(client *Client) {
	cm.Lock()
	defer cm.Unlock()
	cm.clients[client] = true
	log.Printf("Client added: %v", client)
}

func (cm *ChatManager) removeClient(client *Client) {
	cm.Lock()
	defer cm.Unlock()
	if _, ok := cm.clients[client]; ok {
		client.wsConn.Close()
		delete(cm.clients, client)
	}
	fmt.Println("client removed")
}

func checkOrigin(r *http.Request) bool {
	// Add your logic to check the origin of the request here
	return true
}
