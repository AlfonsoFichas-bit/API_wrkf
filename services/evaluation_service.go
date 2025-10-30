package services

import (
	"fmt"
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// EvaluationService handles the business logic for evaluations.
type EvaluationService struct {
	Repo           *storage.EvaluationRepository
	TaskService    *TaskService    // Dependency to check task properties
	ProjectService *ProjectService // Dependency to check user permissions
}

// NewEvaluationService creates a new instance of EvaluationService.
func NewEvaluationService(repo *storage.EvaluationRepository, taskService *TaskService, projectService *ProjectService) *EvaluationService {
	return &EvaluationService{
		Repo:           repo,
		TaskService:    taskService,
		ProjectService: projectService,
	}
}

// CreateEvaluationRequest defines the data needed to create a new evaluation.
type CreateEvaluationRequest struct {
	Score    float32 `json:"score"`
	Comments string  `json:"comments"`
	// Criteria and Feedback can be added here if needed for creation
}

// CreateEvaluation handles the business logic for creating a new evaluation.
func (s *EvaluationService) CreateEvaluation(taskID uint, evaluatorID uint, req CreateEvaluationRequest) (*models.Evaluation, error) {
	// 1. Check if the task is a deliverable
	task, err := s.TaskService.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}
	if !task.IsDeliverable {
		return nil, fmt.Errorf("task is not a deliverable and cannot be evaluated")
	}

	// 2. Check user permissions (example: only product owner or scrum master can evaluate)
	projectID := task.UserStory.ProjectID
	role, err := s.ProjectService.GetUserRoleInProject(evaluatorID, projectID)
	if err != nil {
		return nil, fmt.Errorf("could not verify user permissions")
	}
	if role != string(models.RoleProductOwner) && role != string(models.RoleScrumMaster) {
		return nil, fmt.Errorf("user does not have permission to evaluate tasks in this project")
	}

	// 3. Create the evaluation object
	evaluation := &models.Evaluation{
		TaskID:      taskID,
		EvaluatorID: evaluatorID,
		Score:       req.Score,
		Comments:    req.Comments,
	}

	// 4. Save to the database
	if err := s.Repo.CreateEvaluation(evaluation); err != nil {
		return nil, fmt.Errorf("could not save evaluation: %w", err)
	}

	// 5. Return the fully loaded evaluation
	return s.Repo.GetEvaluationByTaskID(taskID)
}

// GetEvaluation retrieves a single evaluation for a task.
func (s *EvaluationService) GetEvaluation(taskID uint) (*models.Evaluation, error) {
	return s.Repo.GetEvaluationByTaskID(taskID)
}

// UpdateEvaluation handles the business logic for updating an evaluation.
func (s *EvaluationService) UpdateEvaluation(evaluationID uint, evaluatorID uint, req CreateEvaluationRequest) (*models.Evaluation, error) {
	// For now, we'll keep the permission logic simple: only the original evaluator can update.
	evaluation, err := s.Repo.GetEvaluationByTaskID(evaluationID) // Assuming evalID corresponds to taskID for simplicity
	if err != nil {
		return nil, fmt.Errorf("evaluation not found")
	}
	if evaluation.EvaluatorID != evaluatorID {
		return nil, fmt.Errorf("user does not have permission to update this evaluation")
	}

	evaluation.Score = req.Score
	evaluation.Comments = req.Comments

	if err := s.Repo.UpdateEvaluation(evaluation); err != nil {
		return nil, fmt.Errorf("could not update evaluation: %w", err)
	}
	return evaluation, nil
}
