package handlers

import (
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"

	"github.com/labstack/echo/v4"
)

// LoginRequest defines the structure for a login request.
type LoginRequest struct {
	Correo     string `json:"correo" example:"admin@example.com"`
	Contraseña string `json:"contraseña" example:"admin123"`
}

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// Login godoc
// @Summary      User Login
// @Description  Authenticates a user and returns a JWT token.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest  true  "User Credentials"
// @Success      200          {object}  map[string]string
// @Failure      400          {object}  map[string]string
// @Failure      401          {object}  map[string]string
// @Router       /login [post]
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

	return c.JSON(http.StatusOK, user)
}

// GetCurrentUser godoc
// @Summary      Get Current User
// @Description  Retrieves the details of the currently authenticated user.
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200          {object}  models.User
// @Failure      401          {object}  map[string]string
// @Failure      404          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /me [get]
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User ID not found in context"})
	}

	user, err := h.Service.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	user.Contraseña = "" // Never return the password hash
	return c.JSON(http.StatusOK, user)
}
