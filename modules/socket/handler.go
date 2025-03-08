package socket

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

var (
	clients   = make(map[string]*websocket.Conn)
	clientsMu sync.Mutex
)

func WebSocketHandler(c *websocket.Conn) {
	userID := c.Params("user_id")
	if userID == "" {
		log.Println("Missing user ID")
		c.Close()
		return
	}

	clientsMu.Lock()
	clients[userID] = c
	clientsMu.Unlock()

	log.Printf("User %s connected", userID)

	defer func() {
		clientsMu.Lock()
		delete(clients, userID)
		clientsMu.Unlock()
		c.Close()
		log.Printf("User %s disconnected", userID)
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %v", userID, err)
			break
		}
		log.Printf("Received from %s: %s", userID, string(msg))
	}
}

func SendNotificationToUser(userID string, message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	client, exists := clients[userID]
	if !exists {
		log.Printf("User %s not connected", userID)
		return
	}

	err := client.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error sending message to %s: %v", userID, err)
		client.Close()
		delete(clients, userID)
	}
}
