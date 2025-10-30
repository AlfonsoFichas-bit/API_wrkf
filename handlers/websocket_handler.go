package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/websocket"
	"github.com/golang-jwt/jwt/v5"
	gws "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = gws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for dev
	},
}

// WebsocketHandler handles HTTP requests for WebSocket connections.
type WebsocketHandler struct {
	Hub       *websocket.Hub
	JwtSecret []byte
}

// NewWebsocketHandler creates a new instance of WebsocketHandler.
func NewWebsocketHandler(hub *websocket.Hub, jwtSecret string) *WebsocketHandler {
	return &WebsocketHandler{
		Hub:       hub,
		JwtSecret: []byte(jwtSecret),
	}
}

// ServeWs handles websocket requests from the peer.
func (h *WebsocketHandler) ServeWs(c echo.Context) error {
	// 1. Get token from query parameter
	tokenString := c.QueryParam("token")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, "Missing token")
	}

	// 2. Validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return h.JwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, "Invalid token")
	}

	// 3. Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusInternalServerError, "Failed to parse token claims")
	}
	userID, ok := claims["sub"].(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Invalid user ID in token")
	}

	// 4. Get Project ID
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid project ID")
	}

	// 5. Upgrade connection
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

	go client.WritePump()
	go client.ReadPump()

	return nil
}
