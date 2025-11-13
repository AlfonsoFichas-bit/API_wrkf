package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/buga/API_wrkf/models"
)

// Message defines the structure for messages sent over WebSocket.
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// WebSocketManager manages WebSocket clients, registration, unregistration, and message broadcasting.
type WebSocketManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

// NewWebSocketManager creates and returns a new WebSocketManager.
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the WebSocketManager's event loop.
func (m *WebSocketManager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			m.clients[client] = true
			m.mutex.Unlock()
			log.Printf("Client registered: %d", client.userID)

		case client := <-m.unregister:
			m.mutex.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.Send)
				log.Printf("Client unregistered: %d", client.userID)
			}
			m.mutex.Unlock()

		case message := <-m.broadcast:
			m.mutex.RLock()
			for client := range m.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(m.clients, client)
				}
			}
			m.mutex.RUnlock()
		}
	}
}

// RegisterTestClient allows registering a client directly for testing purposes.
func (m *WebSocketManager) RegisterTestClient(client *Client) {
	m.register <- client
}

// BroadcastToProject sends a message to all clients subscribed to a specific project.
func (m *WebSocketManager) BroadcastToProject(projectID uint, message Message) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling broadcast message: %v", err)
		return
	}

	for client := range m.clients {
		// Check if the client is subscribed to the project.
		if _, ok := client.projects[projectID]; ok {
			select {
			case client.Send <- messageBytes:
			default:
				// If the send channel is blocked, assume the client is dead or stuck.
				close(client.Send)
				delete(m.clients, client)
			}
		}
	}
}

// BroadcastTaskStatusUpdated prepares and broadcasts a task status update event.
func (m *WebSocketManager) BroadcastTaskStatusUpdated(projectID, taskID uint, oldStatus, newStatus string, updatedBy *models.User) {
	payload := map[string]interface{}{
		"taskId":    taskID,
		"oldStatus": oldStatus,
		"newStatus": newStatus,
		"updatedBy": map[string]interface{}{
			"id":   updatedBy.ID,
			"name": updatedBy.Nombre,
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	message := Message{
		Type:    "task_status_updated",
		Payload: payload,
	}
	m.BroadcastToProject(projectID, message)
}

// BroadcastTaskCreated prepares and broadcasts a task creation event.
func (m *WebSocketManager) BroadcastTaskCreated(projectID uint, task *models.Task) {
	payload := map[string]interface{}{
		"task":      task,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	message := Message{
		Type:    "task_created",
		Payload: payload,
	}
	m.BroadcastToProject(projectID, message)
}

// BroadcastTaskAssigned prepares and broadcasts a task assignment event.
func (m *WebSocketManager) BroadcastTaskAssigned(projectID, taskID uint, assignedTo, assignedBy *models.User) {
	payload := map[string]interface{}{
		"taskId": taskID,
		"assignedTo": map[string]interface{}{
			"id":   assignedTo.ID,
			"name": assignedTo.Nombre,
		},
		"assignedBy": map[string]interface{}{
			"id":   assignedBy.ID,
			"name": assignedBy.Nombre,
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	message := Message{
		Type:    "task_assigned",
		Payload: payload,
	}
	m.BroadcastToProject(projectID, message)
}

// BroadcastTaskDeleted prepares and broadcasts a task deletion event.
func (m *WebSocketManager) BroadcastTaskDeleted(projectID, taskID uint, deletedBy *models.User) {
	payload := map[string]interface{}{
		"taskId": taskID,
		"deletedBy": map[string]interface{}{
			"id":   deletedBy.ID,
			"name": deletedBy.Nombre,
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	message := Message{
		Type:    "task_deleted",
		Payload: payload,
	}
	m.BroadcastToProject(projectID, message)
}
