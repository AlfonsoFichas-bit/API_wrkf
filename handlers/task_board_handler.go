package handlers

import (
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
)

// TaskBoardHandler handles HTTP requests for the task board.
type TaskBoardHandler struct {
	taskBoardService *services.TaskBoardService
}

// NewTaskBoardHandler creates a new instance of TaskBoardHandler.
func NewTaskBoardHandler(taskBoardService *services.TaskBoardService) *TaskBoardHandler {
	return &TaskBoardHandler{taskBoardService: taskBoardService}
}

// GetTaskBoard retrieves the task board for a given sprint.
// @Summary Get a sprint's task board
// @Description Get a sprint's task board
// @Tags Task Board
// @Accept json
// @Produce json
// @Param sprintId path int true "Sprint ID"
// @Success 200 {object} map[string][]models.Task
// @Router /api/sprints/{sprintId}/taskboard [get]
func (h *TaskBoardHandler) GetTaskBoard(c echo.Context) error {
	sprintID, err := strconv.Atoi(c.Param("sprintId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid sprint ID")
	}

	taskBoard, err := h.taskBoardService.GetTaskBoard(uint(sprintID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, taskBoard)
}