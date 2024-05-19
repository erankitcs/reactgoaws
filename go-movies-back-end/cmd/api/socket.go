package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/internal/models"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.Message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (app *application) MovieChat(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer ws.Close()

	clients[ws] = true

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			//app.errorJSON(w, err)
			delete(clients, ws)
			break
		}
		//fmt.Println(msg)
		broadcast <- msg
	}

}

func handleMovieChatMessages() {
	for {
		msg := <-broadcast
		fmt.Println(msg)
		for client := range clients {

			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
