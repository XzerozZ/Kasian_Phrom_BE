package socket

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]bool)

func WebSocketHandler(c *websocket.Conn) {
	log.Println("Client connected")
	clients[c] = true
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Println("Received message:", string(msg))
	}
}

func BroadcastNotification(message string) {
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error sending message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func CloseConnection(client *websocket.Conn) {
	delete(clients, client)
	client.Close()
}
