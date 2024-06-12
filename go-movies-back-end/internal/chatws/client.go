package chatws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type Client struct {
	id      string
	wsConn  *websocket.Conn
	manager *ChatManager
	movieID int

	// egress is used to avoid concurrent writes on the websocket connection
	egress chan Event
}

func NewClient(wsConn *websocket.Conn, manager *ChatManager, movieID int) *Client {
	id := uuid.New().String()
	return &Client{id, wsConn, manager, movieID, make(chan Event)}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	c.wsConn.SetReadLimit(512)
	fmt.Printf("Read started..")
	if err := c.wsConn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("error setting read deadline: %v", err)
		return
	}
	c.wsConn.SetPongHandler(c.pongHandler)
	for {
		_, payload, err := c.wsConn.ReadMessage()
		fmt.Println("Event received...")
		if err != nil {
			log.Printf("errror reading event: %v", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("errror reading event: %v", err)
			}
			break
		}
		var event Event
		if err := json.Unmarshal(payload, &event); err != nil {
			log.Printf("error marshalling message: %v", err)
			break
		}
		fmt.Println("Event Routed...")
		if err := c.manager.routeEvent(event, c); err != nil {
			log.Printf("error routing event: %v", err)
		}
	}
}

// Socket function to write event to the client
func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.manager.removeClient(c)
	}()
	fmt.Println("Write started...")
	for {
		select {
		case event, ok := <-c.egress:
			fmt.Println("Write Event recieved...")
			if !ok {
				if err := c.wsConn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("connection already closed: %v", err)
				}
				return
			}
			fmt.Println("Sending event...")
			if err := c.wsConn.WriteJSON(event); err != nil {
				log.Printf("failed to send event: %v", err)
			}
			log.Println("event sent")
		case <-ticker.C:
			if err := c.wsConn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("failed to send ping: %v", err)
				return
			}
			fmt.Println("ping sent")
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	fmt.Println("pong received")
	if err := c.wsConn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("error setting read deadline: %v", err)
		return err
	}
	return nil
}
