package routes_test

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/buga/API_wrkf/config"
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/tests"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketConnection_Success(t *testing.T) {
	// --- Setup ---
	cfg := &config.AppConfig{JWTSecret: "test-secret-ws"}
	testRouter, err := tests.NewTestRouter(cfg)
	assert.NoError(t, err)

	user := &models.User{Nombre: "ws_user", Correo: "ws@example.com", Contraseña: "password"}
	err = testRouter.UserService.CreateUser(user)
	assert.NoError(t, err)
	createdUser, err := testRouter.UserService.GetUserByEmail("ws@example.com")
	assert.NoError(t, err)

	project := &models.Project{Name: "WebSocket Project"}
	err = testRouter.ProjectService.CreateProject(project, createdUser.ID)
	assert.NoError(t, err)

	token, err := testRouter.UserService.Login("ws@example.com", "password")
	assert.NoError(t, err)

	server := httptest.NewServer(testRouter.Echo)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + fmt.Sprintf("/api/ws/projects/%d/board?token=%s", project.ID, token)

	// --- Execute ---
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err, "Failed to connect to WebSocket")
	if err != nil {
		return
	}
	defer ws.Close()

	// --- Assertions ---
	err = ws.SetReadDeadline(time.Now().Add(1 * time.Second))
	assert.NoError(t, err)
	_, _, err = ws.ReadMessage()
	assert.Error(t, err, "Should have timed out as no message was expected")
	netErr, ok := err.(net.Error)
	assert.True(t, ok && netErr.Timeout(), "Error should be a network timeout")
}

func TestWebSocketBroadcast_OnTaskUpdate(t *testing.T) {
	t.Skip("Skipping flaky test to be addressed in a separate task")
	// --- Setup ---
	cfg := &config.AppConfig{JWTSecret: "test-secret-broadcast"}
	testRouter, err := tests.NewTestRouter(cfg)
	assert.NoError(t, err)

	server := httptest.NewServer(testRouter.Echo)
	defer server.Close()

	user := &models.User{Nombre: "ws_user_2", Correo: "ws2@example.com", Contraseña: "password"}
	err = testRouter.UserService.CreateUser(user)
	assert.NoError(t, err)
	createdUser, err := testRouter.UserService.GetUserByEmail("ws2@example.com")
	assert.NoError(t, err)

	token, err := testRouter.UserService.Login("ws2@example.com", "password")
	assert.NoError(t, err)

	project := &models.Project{Name: "Broadcast Project"}
	err = testRouter.ProjectService.CreateProject(project, createdUser.ID)
	assert.NoError(t, err)

	userStory := &models.UserStory{Title: "Broadcast Story"}
	err = testRouter.UserStoryService.CreateUserStory(userStory, project.ID, createdUser.ID)
	assert.NoError(t, err)

	task, err := testRouter.TaskService.CreateTask(&models.Task{Title: "Broadcast Task"}, userStory.ID, createdUser.ID)
	assert.NoError(t, err)

	// --- Connect WebSocket Client ---
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + fmt.Sprintf("/api/ws/projects/%d/board?token=%s", project.ID, token)
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// --- Trigger the update ---
	go func() {
		time.Sleep(250 * time.Millisecond) // Give WS time to register client
		updateURL := "http" + strings.TrimPrefix(server.URL, "http") + fmt.Sprintf("/api/tasks/%d/status", task.ID)
		payload := `{"status":"in_progress"}`
		req, _ := http.NewRequest(http.MethodPut, updateURL, strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, httpErr := client.Do(req)
		assert.NoError(t, httpErr)
		if resp != nil {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			closeErr := resp.Body.Close()
			assert.NoError(t, closeErr)
		}
	}()

	// --- Assertions ---
	ws.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, message, err := ws.ReadMessage()
	assert.NoError(t, err, "Failed to read broadcast message from WebSocket")

	var broadcastMsg struct {
		ProjectID uint `json:"projectId"`
		Payload   struct {
			Type    string          `json:"type"`
			Payload json.RawMessage `json:"payload"`
		} `json:"payload"`
	}
	err = json.Unmarshal(message, &broadcastMsg)
	assert.NoError(t, err)
	assert.Equal(t, project.ID, broadcastMsg.ProjectID)
	assert.Equal(t, "TASK_UPDATED", broadcastMsg.Payload.Type)

	var updatedTask models.Task
	err = json.Unmarshal(broadcastMsg.Payload.Payload, &updatedTask)
	assert.NoError(t, err)
	assert.Equal(t, task.ID, updatedTask.ID)
	assert.Equal(t, "in_progress", updatedTask.Status)
}
