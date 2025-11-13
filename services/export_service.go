package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// ExportService handles the business logic for exporting data.
type ExportService struct {
	ProjectRepo   *storage.ProjectRepository
	UserStoryRepo *storage.UserStoryRepository
	TaskRepo      *storage.TaskRepository
}

// NewExportService creates a new instance of ExportService.
func NewExportService(
	projectRepo *storage.ProjectRepository,
	userStoryRepo *storage.UserStoryRepository,
	taskRepo *storage.TaskRepository,
) *ExportService {
	return &ExportService{
		ProjectRepo:   projectRepo,
		UserStoryRepo: userStoryRepo,
		TaskRepo:      taskRepo,
	}
}

// ExportProjectToCSV generates a CSV file in memory for a given project.
func (s *ExportService) ExportProjectToCSV(projectID uint) ([]byte, error) {
	// 1. Fetch all necessary data.
	project, err := s.ProjectRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	userStories, err := s.UserStoryRepo.GetUserStoriesByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch user stories: %w", err)
	}

	var allTasks []models.Task
	for _, us := range userStories {
		tasks, err := s.TaskRepo.GetTasksByUserStoryID(us.ID)
		if err != nil {
			return nil, fmt.Errorf("could not fetch tasks for user story %d: %w", us.ID, err)
		}
		allTasks = append(allTasks, tasks...)
	}

	// 2. Generate CSV content in memory.
	b := new(bytes.Buffer)
	w := csv.NewWriter(b)

	// Write header
	header := []string{
		"Project Name",
		"User Story ID",
		"User Story Title",
		"User Story Status",
		"Task ID",
		"Task Title",
		"Task Status",
		"Assigned To",
	}
	if err := w.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// 3. Write data rows.
	taskMap := make(map[uint][]models.Task)
	for _, task := range allTasks {
		taskMap[task.UserStoryID] = append(taskMap[task.UserStoryID], task)
	}

	for _, us := range userStories {
		tasksInStory := taskMap[us.ID]
		if len(tasksInStory) == 0 {
			// Write a row even for user stories with no tasks
			row := []string{
				project.Name,
				strconv.Itoa(int(us.ID)),
				us.Title,
				us.Status,
				"", // No Task ID
				"", // No Task Title
				"", // No Task Status
				"", // No Assigned To
			}
			if err := w.Write(row); err != nil {
				return nil, fmt.Errorf("failed to write CSV row: %w", err)
			}
		} else {
			for _, task := range tasksInStory {
				assignedTo := "Unassigned"
				if task.AssignedTo != nil {
					assignedTo = task.AssignedTo.Nombre
				}
				row := []string{
					project.Name,
					strconv.Itoa(int(us.ID)),
					us.Title,
					us.Status,
					strconv.Itoa(int(task.ID)),
					task.Title,
					string(task.Status),
					assignedTo,
				}
				if err := w.Write(row); err != nil {
					return nil, fmt.Errorf("failed to write CSV row: %w", err)
				}
			}
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, fmt.Errorf("error flushing CSV writer: %w", err)
	}

	return b.Bytes(), nil
}
