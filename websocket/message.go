package websocket

import "encoding/json"

// MessageType defines the type of a WebSocket message.
type MessageType string

const (
	// MessageTypeTaskCreated is sent when a new task is created.
	MessageTypeTaskCreated MessageType = "TASK_CREATED"
	// MessageTypeTaskUpdated is sent when a task is updated (e.g., status, assignment, title).
	MessageTypeTaskUpdated MessageType = "TASK_UPDATED"
)

// Message represents the structure of a message sent over WebSocket.
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
