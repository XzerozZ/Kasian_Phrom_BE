package socket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/gofiber/websocket/v2"
)

var (
	clients sync.Map
)

func WebSocketHandler(c *websocket.Conn) {
	userID := c.Params("user_id")
	if userID == "" {
		log.Println("Missing user ID")
		c.Close()
		return
	}

	clients.Store(userID, c)
	log.Printf("User %s connected", userID)

	defer func() {
		clients.Delete(userID)
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

func SendNotificationToUser(userID string, noti entities.Notification) {
	client, exists := clients.Load(userID)
	if !exists {
		log.Printf("User %s not connected", userID)
		return
	}

	conn, ok := client.(*websocket.Conn)
	if !ok {
		log.Printf("Invalid WebSocket connection for user %s", userID)
		return
	}

	notiJSON, err := json.Marshal(noti)
	if err != nil {
		log.Printf("Error marshaling notification: %v", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, notiJSON)
	if err != nil {
		log.Printf("Error sending message to %s: %v", userID, err)
		conn.Close()
		clients.Delete(userID)
	}
}
