package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/websocket"
	gws "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = gws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default for development.
		// For production, you should implement a proper origin check.
		return true
	},
}

// WebsocketHandler handles HTTP requests for WebSocket connections.
type WebsocketHandler struct {
	Hub *websocket.Hub
}

// NewWebsocketHandler creates a new instance of WebsocketHandler.
func NewWebsocketHandler(hub *websocket.Hub) *WebsocketHandler {
	return &WebsocketHandler{Hub: hub}
}

// ServeWs handles websocket requests from the peer.
func (h *WebsocketHandler) ServeWs(c echo.Context) error {
	userID, ok := c.Get("userID").(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Invalid or missing user ID from token")
	}

	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid project ID")
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	client := &websocket.Client{
		Hub:       h.Hub,
		Conn:      conn,
		Send:      make(chan []byte, 256),
		UserID:    uint(userID),
		ProjectID: uint(projectID),
	}
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()

	return nil
}
