package websocket

import (
	"errors"
	"log"
	"net/http"

	"github.com/buga/API_wrkf/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default. In production, you should have a whitelist.
		return true
	},
}

// WebSocketHandler handles upgrading HTTP connections to WebSocket connections.
type WebSocketHandler struct {
	hub           *WebSocketManager
	jwtSecret     string
	userService   *services.UserService
	projectService *services.ProjectService // Added ProjectService
}

// NewWebSocketHandler creates a new WebSocketHandler.
func NewWebSocketHandler(hub *WebSocketManager, jwtSecret string, userService *services.UserService, projectService *services.ProjectService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:           hub,
		jwtSecret:     jwtSecret,
		userService:   userService,
		projectService: projectService, // Initialize ProjectService
	}
}

// validateTokenAndGetUserID parses a JWT token string, validates it, and returns the user ID.
func (h *WebSocketHandler) validateTokenAndGetUserID(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(h.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(float64); ok {
			return uint(sub), nil
		}
	}

	return 0, errors.New("invalid token claims")
}

// HandleConnection handles the WebSocket connection request.
func (h *WebSocketHandler) HandleConnection(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return c.String(http.StatusUnauthorized, "Missing token")
	}

	userID, err := h.validateTokenAndGetUserID(token)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Invalid token")
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return err
	}

	// Get the list of projects the user is a member of.
	userProjects, err := h.projectService.GetProjectsByUserID(userID)
	if err != nil {
		log.Printf("Failed to get projects for user %d: %v", userID, err)
		// Continue with an empty project list if there's an error, or close connection.
		// For now, we'll proceed with an empty list to avoid breaking the connection entirely.
	}

	projects := make(map[uint]bool)
	for _, p := range userProjects {
		projects[p.ID] = true
	}

	client := &Client{
		hub:      h.hub,
		conn:     conn,
		Send:     make(chan []byte, 256),
		userID:   userID,
		projects: projects,
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.writePump()
	go client.readPump()

	return nil
}
