package chatws

import (
	"backend/internal/models"
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
	// This is mae sure that secure authenticated socket is returning common agreed supprotocol
	Subprotocols: []string{"chat"},
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

func (cm *ChatManager) routeEvent(event models.Event, c *Client) error {
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

func (cm *ChatManager) ServeChat(w http.ResponseWriter, r *http.Request, movieID int) error {
	fmt.Println("connection recieved..")
	// This also works. h can be pass into Upgrade function instead of nill. However, i have preferred to set it at upgrader initialization
	//h := http.Header{}
	// Setting Command understanding protocol as chat
	//h.Set("Sec-Websocket-Protocol", "chat")

	// Get User name from Context of request
	userName := r.Context().Value("username").(string)

	// Upgrade the HTTP connection to a WebSocket connection
	// This will not return until the connection has been established.
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	//defer ws.Close()
	// Create a new client object
	client := NewClient(ws, cm, movieID, userName)

	// Add client to the list of clients
	cm.addClient(client)

	// Start client process
	go client.readMessages()
	go client.writeMessages()

	//Broadcast to other clients
	if err := UserJoinedHandler(client); err != nil {
		log.Printf("error in broadcasting user left event: %v", err)
	}
	return nil
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
		// Broadcast to other clients
		if err := UserLeftHandler(client); err != nil {
			log.Printf("error in broadcasting user left event: %v", err)
		}

		fmt.Println("client removed")
	}
	fmt.Println("client is already removed")
}

func checkOrigin(r *http.Request) bool {
	// Add your logic to check the origin of the request here
	return true
}
