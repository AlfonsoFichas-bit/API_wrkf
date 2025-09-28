package services

import (
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// TaskBoardService handles business logic for the task board.
type TaskBoardService struct {
	taskRepo *storage.TaskRepository
}

// NewTaskBoardService creates a new instance of TaskBoardService.
func NewTaskBoardService(taskRepo *storage.TaskRepository) *TaskBoardService {
	return &TaskBoardService{taskRepo: taskRepo}
}

// GetTaskBoard retrieves all tasks for a given sprint, organized by status.
func (s *TaskBoardService) GetTaskBoard(sprintID uint) (map[string][]models.Task, error) {
	tasks, err := s.taskRepo.GetTasksBySprintID(sprintID)
	if err != nil {
		return nil, err
	}

	taskBoard := make(map[string][]models.Task)
	for _, task := range tasks {
		taskBoard[task.Status] = append(taskBoard[task.Status], task)
	}

	return taskBoard, nil
}