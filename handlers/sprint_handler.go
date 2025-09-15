package handlers

import (
	"API_wrkf/models"
	"API_wrkf/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SprintHandler struct {
	Service *services.SprintService
}

func NewSprintHandler(service *services.SprintService) *SprintHandler {
	return &SprintHandler{Service: service}
}

func (h *SprintHandler) CreateSprint(c echo.Context) error {
	sprint := new(models.Sprint)
	if err := c.Bind(sprint); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	userID, _ := c.Get("userID").(float64)

	if err := h.Service.CreateSprint(sprint, uint(userID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create sprint"})
	}

	return c.JSON(http.StatusCreated, sprint)
}
