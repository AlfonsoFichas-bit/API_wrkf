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

	user.Contraseña = "" // Never return the password hash
	return c.JSON(http.StatusOK, user)
}

// GetAllUsers handles retrieving all users.
// @Summary Get All Users
// @Description Retrieves a list of all users. This is an admin-only endpoint.
// @Tags Users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Router /api/admin/users [get]
func (h *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not retrieve users"})
	}
	return c.JSON(http.StatusOK, users)
}

// UpdateUser handles updating a user's information.
// @Summary Update User
// @Description Updates a user's information. This is an admin-only endpoint.
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User data to update"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/users/{id} [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	updatedUser, err := h.Service.UpdateUser(uint(id), user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update user", "details": err.Error()})
	}

	return c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser handles deleting a user by their ID.
// @Summary Delete User
// @Description Deletes a user by their ID. This is an admin-only endpoint.
// @Tags Users
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	if err := h.Service.DeleteUser(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete user", "details": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
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
	userIDFloat, ok := c.Get("userID").(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User ID not found or is of incorrect type in context"})
	}
	userID := uint(userIDFloat)

	user, err := h.Service.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	user.Contraseña = "" // Never return the password hash
	return c.JSON(http.StatusOK, user)
}

// Logout godoc
// @Summary      User Logout
// @Description  Logs out the user. In a JWT-based system, this is primarily a client-side action. This endpoint is provided for completeness.
// @Tags         Authentication
// @Security     ApiKeyAuth
// @Success      200  {object}  map[string]string
// @Router       /api/logout [post]
func (h *UserHandler) Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}
