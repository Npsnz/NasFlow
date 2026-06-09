package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type EventMessage struct {
	UserID  uint        `json:"user_id"`
	Event   string      `json:"event"` // e.g. task.created, task.updated, task.deleted, task.moved, tasks.overdue
	Payload interface{} `json:"payload"`
}

type SSEBroker struct {
	clients    map[chan EventMessage]uint // chan -> userID
	register   chan clientInfo
	unregister chan chan EventMessage
	broadcast  chan EventMessage
	mutex      sync.RWMutex
}

type clientInfo struct {
	ch     chan EventMessage
	userID uint
}

var Broker *SSEBroker

func InitSSEBroker() {
	Broker = &SSEBroker{
		clients:    make(map[chan EventMessage]uint),
		register:   make(chan clientInfo),
		unregister: make(chan chan EventMessage),
		broadcast:  make(chan EventMessage),
	}
	go Broker.run()
}

func (b *SSEBroker) run() {
	for {
		select {
		case client := <-b.register:
			b.mutex.Lock()
			b.clients[client.ch] = client.userID
			b.mutex.Unlock()
			log.Printf("[SSE] Client connected for User %d. Total clients: %d", client.userID, len(b.clients))

		case ch := <-b.unregister:
			b.mutex.Lock()
			if _, exists := b.clients[ch]; exists {
				delete(b.clients, ch)
				close(ch)
			}
			b.mutex.Unlock()
			log.Printf("[SSE] Client disconnected. Total clients: %d", len(b.clients))

		case msg := <-b.broadcast:
			b.mutex.RLock()
			for ch, clientUserID := range b.clients {
				if clientUserID == msg.UserID {
					select {
					case ch <- msg:
					default:
						// Skip client if channel buffer is full
					}
				}
			}
			b.mutex.RUnlock()
		}
	}
}

// BroadcastEvent publishes an event to the SSE channel of a user
func BroadcastEvent(userID uint, event string, payload interface{}) {
	if Broker == nil {
		return
	}
	Broker.broadcast <- EventMessage{
		UserID:  userID,
		Event:   event,
		Payload: payload,
	}
}

// HandleSSE establishes SSE connection for a user
func HandleSSE(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized connection",
			"code":  "UNAUTHORIZED",
		})
		return
	}
	userID := val.(uint)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// Each connected tab gets a channel
	ch := make(chan EventMessage, 20)
	Broker.register <- clientInfo{ch: ch, userID: userID}

	defer func() {
		Broker.unregister <- ch
	}()

	// Send an initial handshake comment
	c.SSEvent("handshake", "connected")
	c.Writer.Flush()

	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-ch:
			if !ok {
				return false
			}
			dataBytes, err := json.Marshal(msg.Payload)
			if err != nil {
				return true
			}
			c.SSEvent(msg.Event, string(dataBytes))
			c.Writer.Flush()
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}
