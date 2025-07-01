package tui

import (
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

type WebSocketManager struct {
	conn      *websocket.Conn
	state     *AppState
	app       *tview.Application
	connected bool
	mu        sync.RWMutex
	stopCh    chan struct{}
}

func NewWebSocketManager(state *AppState, app *tview.Application) *WebSocketManager {
	return &WebSocketManager{
		state:  state,
		app:    app,
		stopCh: make(chan struct{}),
	}
}

func (wsm *WebSocketManager) Connect(userID string) {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	if wsm.connected {
		return
	}

	wsURL := "ws://localhost:8080/ws/" + userID
	u, err := url.Parse(wsURL)
	if err != nil {
		log.Printf("WebSocket URL parse error: %v", err)
		return
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("WebSocket connection error: %v", err)
		return
	}

	wsm.conn = conn
	wsm.connected = true

	go wsm.readMessages()
}

func (wsm *WebSocketManager) readMessages() {
	defer func() {
		wsm.mu.Lock()
		if wsm.conn != nil {
			wsm.conn.Close()
		}
		wsm.connected = false
		wsm.mu.Unlock()
	}()

	for {
		select {
		case <-wsm.stopCh:
			return
		default:
			var msg map[string]interface{}
			err := wsm.conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				return
			}

			// Convert to notification
			notif := Notification{
				ID:      time.Now().Format("20060102150405"),
				Type:    "notification",
				Message: "New notification received",
				Time:    time.Now().Format("15:04:05"),
			}

			if msgType, ok := msg["type"].(string); ok {
				notif.Type = msgType
			}
			if message, ok := msg["message"].(string); ok {
				notif.Message = message
			}

			wsm.state.AddNotification(notif)

			// Update UI
			wsm.app.QueueUpdateDraw(func() {
				log.Printf("🔔 Notification: %s - %s", notif.Type, notif.Message)
			})
		}
	}
}

func (wsm *WebSocketManager) IsConnected() bool {
	wsm.mu.RLock()
	defer wsm.mu.RUnlock()
	return wsm.connected
}

func (wsm *WebSocketManager) Disconnect() {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	if !wsm.connected {
		return
	}

	close(wsm.stopCh)
	if wsm.conn != nil {
		wsm.conn.Close()
	}
	wsm.connected = false
}
