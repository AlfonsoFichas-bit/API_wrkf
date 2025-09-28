package handlers

import (
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
)

// AdminHandler handles HTTP requests for admin-related actions.
type AdminHandler struct {
	userService *services.UserService
}

// NewAdminHandler creates a new instance of AdminHandler.
func NewAdminHandler(userService *services.UserService) *AdminHandler {
	return &AdminHandler{userService: userService}
}

// GetAllUsers retrieves all users.
// @Summary Get all users
// @Description Get all users
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /api/admin/users [get]
func (h *AdminHandler) GetAllUsers(c echo.Context) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

// DeleteUser deletes a user.
// @Summary Delete a user
// @Description Delete a user
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Router /api/admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user ID")
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// CreateAdminUser handles the creation of a new admin user.
func (h *AdminHandler) CreateAdminUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.userService.CreateAdminUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create admin user", "details": err.Error()})
	}

	user.Contraseña = ""
	return c.JSON(http.StatusCreated, user)
}

// CreateUser handles the creation of a new standard platform user.
func (h *AdminHandler) CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.userService.CreateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create user", "details": err.Error()})
	}

	user.Contraseña = ""
	return c.JSON(http.StatusCreated, user)
}