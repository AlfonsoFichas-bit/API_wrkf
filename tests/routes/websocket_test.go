package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/websocket"
	gws "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestWebSocketBoardUpdate(t *testing.T) {
	e, db := setupTestApp()
	server := httptest.NewServer(e)
	defer server.Close()

	// --- Setup Test Data ---
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{Nombre: "ws_user", Correo: "ws@example.com", Contraseña: string(hashedPassword)}
	db.Create(&user)
	project := models.Project{Name: "WS Project", CreatedByID: user.ID}
	db.Create(&project)
	userStory := models.UserStory{Title: "WS Story", ProjectID: project.ID, CreatedByID: user.ID}
	db.Create(&userStory)
	task := models.Task{Title: "WS Task", UserStoryID: userStory.ID, Status: string(models.StatusTodo), CreatedByID: user.ID}
	db.Create(&task)

	// --- Create a WebSocket client ---
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/ws/projects/" + strconv.Itoa(int(project.ID)) + "/board"

	token, err := loginAndGetToken(server.URL, "ws@example.com", password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token, "login should return a token")
	header := http.Header{"Authorization": {"Bearer " + token}}

	conn, _, err := gws.DefaultDialer.Dial(wsURL, header)
	assert.NoError(t, err)
	defer conn.Close()

	// --- Listen for messages in a goroutine ---
	msgChan := make(chan []byte)
	go func() {
		defer close(msgChan)
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Logf("read error: %v", err)
			return
		}
		msgChan <- msg
	}()

	// --- Trigger the update via HTTP ---
	updateStatus := map[string]string{"status": string(models.StatusInProgress)}
	jsonBody, _ := json.Marshal(updateStatus)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/tasks/"+strconv.Itoa(int(task.ID))+"/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// --- Assert the WebSocket message ---
	select {
	case msg := <-msgChan:
		var wsMsg websocket.Message
		err := json.Unmarshal(msg, &wsMsg)
		assert.NoError(t, err)
		assert.Equal(t, websocket.MessageTypeTaskUpdated, wsMsg.Type)

		var updatedTask models.Task
		err = json.Unmarshal(wsMsg.Payload, &updatedTask)
		assert.NoError(t, err)
		assert.Equal(t, task.ID, updatedTask.ID)
		assert.Equal(t, string(models.StatusInProgress), updatedTask.Status)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for WebSocket message")
	}
}

func TestWebSocketTaskCreation(t *testing.T) {
	e, db := setupTestApp()
	server := httptest.NewServer(e)
	defer server.Close()

	// --- Setup Test Data ---
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{Nombre: "ws_creator", Correo: "ws_creator@example.com", Contraseña: string(hashedPassword)}
	db.Create(&user)
	project := models.Project{Name: "WS Create Project", CreatedByID: user.ID}
	db.Create(&project)
	userStory := models.UserStory{Title: "WS Create Story", ProjectID: project.ID, CreatedByID: user.ID}
	db.Create(&userStory)

	// --- Create a WebSocket client ---
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/ws/projects/" + strconv.Itoa(int(project.ID)) + "/board"
	token, err := loginAndGetToken(server.URL, "ws_creator@example.com", password)
	assert.NoError(t, err)
	header := http.Header{"Authorization": {"Bearer " + token}}
	conn, _, err := gws.DefaultDialer.Dial(wsURL, header)
	assert.NoError(t, err)
	defer conn.Close()

	// --- Listen for messages ---
	msgChan := make(chan []byte)
	go func() {
		defer close(msgChan)
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Logf("read error: %v", err)
			return
		}
		msgChan <- msg
	}()

	// --- Trigger the creation via HTTP ---
	newTask := map[string]string{"title": "Newly Created Task"}
	jsonBody, _ := json.Marshal(newTask)
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/userstories/"+strconv.Itoa(int(userStory.ID))+"/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// --- Assert the WebSocket message ---
	select {
	case msg := <-msgChan:
		var wsMsg websocket.Message
		err := json.Unmarshal(msg, &wsMsg)
		assert.NoError(t, err)
		assert.Equal(t, websocket.MessageTypeTaskCreated, wsMsg.Type)

		var createdTask models.Task
		err = json.Unmarshal(wsMsg.Payload, &createdTask)
		assert.NoError(t, err)
		assert.Equal(t, "Newly Created Task", createdTask.Title)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for WebSocket message")
	}
}

// Helper function to log in and get a JWT token for tests
func loginAndGetToken(serverURL, email, password string) (string, error) {
	loginCreds := map[string]string{"correo": email, "contraseña": password}
	jsonBody, _ := json.Marshal(loginCreds)
	resp, err := http.Post(serverURL+"/login", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenMap map[string]string
	json.NewDecoder(resp.Body).Decode(&tokenMap)
	return tokenMap["token"], nil
}
