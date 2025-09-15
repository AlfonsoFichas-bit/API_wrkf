
package handlers

import (
	"API_wrkf/models"
	"API_wrkf/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// LoginRequest defines the structure for a login request.
type LoginRequest struct {
	Correo     string `json:"correo"`
	Contraseña string `json:"contraseña"`
}

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// Login handles user login and token generation.
func (h *UserHandler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	token, err := h.Service.Login(req.Correo, req.Contraseña)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

// CreateAdminUser handles the creation of a new admin user.
func (h *UserHandler) CreateAdminUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.Service.CreateAdminUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create admin user", "details": err.Error()})
	}

	user.Contraseña = ""
	return c.JSON(http.StatusCreated, user)
}

// CreateUser handles the creation of a new standard platform user.
func (h *UserHandler) CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.Service.CreateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create user", "details": err.Error()})
	}

	user.Contraseña = ""
	return c.JSON(http.StatusCreated, user)
}

// GetUser handles retrieving a user by their ID.
func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	user, err := h.Service.GetUserByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	user.Contraseña = ""
	return c.JSON(http.StatusOK, user)
}
