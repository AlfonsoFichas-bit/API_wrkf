package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/websocket"
	"github.com/stretchr/testify/assert"
)

// TestWebSocketManager_BroadcastTaskStatusUpdated tests that the JSON message
// for a task status update is formatted correctly.
func TestWebSocketManager_BroadcastTaskStatusUpdated(t *testing.T) {
	// 1. Setup
	wsManager := websocket.NewWebSocketManager()
	go wsManager.Run() // Run the manager in a goroutine to handle channels

	// Create a mock client and register it
	// In a real scenario, this client would be created by an incoming WebSocket connection.
	client := websocket.NewTestClient(wsManager, 1, map[uint]bool{100: true})
	wsManager.RegisterTestClient(client)

	// Wait a moment for registration to complete
	time.Sleep(10 * time.Millisecond)

	// 2. Test Data
	projectID := uint(100)
	taskID := uint(200)
	oldStatus := "todo"
	newStatus := "in_progress"
	updater := &models.User{ID: 1, Nombre: "Test User"}

	// 3. Execute
	wsManager.BroadcastTaskStatusUpdated(projectID, taskID, oldStatus, newStatus, updater)

	// 4. Assert
	select {
	case msgBytes := <-client.Send:
		var receivedMsg websocket.Message
		err := json.Unmarshal(msgBytes, &receivedMsg)
		assert.NoError(t, err, "Should be able to unmarshal the message")

		// Check message type
		assert.Equal(t, "task_status_updated", receivedMsg.Type, "Message type should be correct")

		// Check payload content
		payload, ok := receivedMsg.Payload.(map[string]interface{})
		assert.True(t, ok, "Payload should be a map")

		assert.Equal(t, float64(taskID), payload["taskId"], "Task ID in payload should be correct")
		assert.Equal(t, oldStatus, payload["oldStatus"], "Old status in payload should be correct")
		assert.Equal(t, newStatus, payload["newStatus"], "New status in payload should be correct")

		updaterPayload, ok := payload["updatedBy"].(map[string]interface{})
		assert.True(t, ok, "Updater payload should be a map")
		assert.Equal(t, float64(updater.ID), updaterPayload["id"], "Updater ID should be correct")
		assert.Equal(t, updater.Nombre, updaterPayload["name"], "Updater name should be correct")

	case <-time.After(1 * time.Second):
		t.Fatal("Test timed out: did not receive a message on the client's send channel")
	}
}
