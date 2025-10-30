package websocket

import (
	"encoding/json"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[uint]*Client

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[uint]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserID] = client
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			var msgData struct {
				ProjectID uint    `json:"projectId"`
				Payload   Message `json:"payload"`
			}
			if err := json.Unmarshal(message, &msgData); err != nil {
				log.Printf("error unmarshalling broadcast message: %v", err)
				continue
			}

			// Re-marshal the inner payload to send to clients
			finalPayload, err := json.Marshal(msgData.Payload)
			if err != nil {
				log.Printf("error re-marshalling payload for client: %v", err)
				continue
			}

			for _, client := range h.Clients {
				if client.ProjectID == msgData.ProjectID {
					select {
					case client.Send <- finalPayload:
					default:
						close(client.Send)
						delete(h.Clients, client.UserID)
					}
				}
			}
		}
	}
}
