package services

import (
	"fmt"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// EvaluationService handles the business logic for evaluations.
type EvaluationService struct {
	EvalRepo       *storage.EvaluationRepository
	TaskRepo       *storage.TaskRepository
	RubricRepo     storage.RubricRepository // Use interface type
	ProjectService *ProjectService          // To check user roles
}

// NewEvaluationService creates a new instance of EvaluationService.
func NewEvaluationService(evalRepo *storage.EvaluationRepository, taskRepo *storage.TaskRepository, rubricRepo storage.RubricRepository, projectService *ProjectService) *EvaluationService {
	return &EvaluationService{
		EvalRepo:       evalRepo,
		TaskRepo:       taskRepo,
		RubricRepo:     rubricRepo,
		ProjectService: projectService,
	}
}

// CreateEvaluationRequest defines the structure for the evaluation payload.
type CreateEvaluationRequest struct {
	RubricID             uint                          `json:"rubricId"`
	OverallFeedback      string                        `json:"overallFeedback"`
	CriterionEvaluations []CriterionEvaluationRequest `json:"criterionEvaluations"`
}

// CriterionEvaluationRequest defines the structure for a single criterion's evaluation.
type CriterionEvaluationRequest struct {
	CriterionID uint    `json:"criterionId"`
	Score       float64 `json:"score"`
	Feedback    string  `json:"feedback"`
}

// CreateEvaluation handles the complex business logic for creating a new evaluation.
func (s *EvaluationService) CreateEvaluation(taskID uint, evaluatorID uint, req CreateEvaluationRequest) (*models.Evaluation, error) {
	// 1. Verify that the task exists.
	_, err := s.TaskRepo.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task with ID %d not found", taskID)
	}

	// 2. Get the project and verify the evaluator is a teacher ('docente').
	projectID, err := s.TaskRepo.GetProjectIDForTask(taskID)
	if err != nil {
		return nil, fmt.Errorf("could not find project for task %d", taskID)
	}
	role, err := s.ProjectService.GetUserRoleInProject(evaluatorID, projectID)
	if err != nil {
		return nil, fmt.Errorf("could not verify user role in project")
	}
	if role != "docente" { // Assuming 'docente' is the role for teachers.
		return nil, fmt.Errorf("user does not have permission to evaluate this task")
	}

	// 3. Verify that the rubric exists and belongs to the same project.
	rubric, err := s.RubricRepo.FindByID(req.RubricID)
	if err != nil {
		return nil, fmt.Errorf("rubric with ID %d not found", req.RubricID)
	}
	if rubric.ProjectID != projectID {
		return nil, fmt.Errorf("rubric does not belong to the same project as the task")
	}

	// 4. Build the Evaluation model from the request.
	evaluation := &models.Evaluation{
		TaskID:          taskID,
		EvaluatorID:     evaluatorID,
		RubricID:        req.RubricID,
		OverallFeedback: req.OverallFeedback,
		Status:          "published", // Or could be 'draft' initially
	}

	// 5. Process criterion evaluations and calculate total score.
	var totalScore float64
	for _, critEvalReq := range req.CriterionEvaluations {
		evaluation.CriterionEvaluations = append(evaluation.CriterionEvaluations, models.CriterionEvaluation{
			CriterionID: critEvalReq.CriterionID,
			Score:       critEvalReq.Score,
			Feedback:    critEvalReq.Feedback,
		})
		totalScore += critEvalReq.Score
	}
	evaluation.TotalScore = totalScore

	// 6. Save the evaluation to the database.
	if err := s.EvalRepo.CreateEvaluation(evaluation); err != nil {
		return nil, fmt.Errorf("could not save evaluation: %w", err)
	}

	return evaluation, nil
}

// GetEvaluationsByTaskID retrieves all evaluations for a specific task.
func (s *EvaluationService) GetEvaluationsByTaskID(taskID uint) ([]models.Evaluation, error) {
	return s.EvalRepo.GetEvaluationsByTaskID(taskID)
}
