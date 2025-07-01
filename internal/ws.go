package internal

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func WatchNotifications(userID string) {
	url := "ws://localhost:8080/ws/" + userID
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("WebSocket error: %v", err)
	}
	defer conn.Close()

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("WebSocket read error:", err)
			time.Sleep(time.Second)
			continue
		}
		log.Printf("🔔 Notification: %+v", msg)
	}
}
