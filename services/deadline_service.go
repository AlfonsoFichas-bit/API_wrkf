package services

import (
	"fmt"
	"time"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// DeadlineService handles the business logic for deadlines.
type DeadlineService struct {
	Repo           *storage.DeadlineRepository
	ProjectService *ProjectService
}

// NewDeadlineService creates a new instance of DeadlineService.
func NewDeadlineService(repo *storage.DeadlineRepository, projectService *ProjectService) *DeadlineService {
	return &DeadlineService{
		Repo:           repo,
		ProjectService: projectService,
	}
}

// GetUpcomingDeadlines retrieves upcoming deadlines based on user role and filters.
func (s *DeadlineService) GetUpcomingDeadlines(userID uint, userRole string, days int, projectID *uint, deadlineType *string) ([]models.Deadline, error) {
	var projectIDs []uint

	// Determine which projects the user can see deadlines for.
	if projectID != nil {
		// If a specific project is requested, verify user membership.
		isMember, err := s.ProjectService.IsUserMemberOfProject(userID, *projectID)
		if err != nil {
			return nil, fmt.Errorf("could not verify project membership")
		}
		if !isMember {
			return nil, fmt.Errorf("user is not a member of the requested project")
		}
		projectIDs = append(projectIDs, *projectID)
	} else {
		// If no specific project, get all projects the user is a member of.
		// This applies to both students and teachers.
		projects, err := s.ProjectService.GetProjectsByUserID(userID)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve user's projects")
		}
		for _, p := range projects {
			projectIDs = append(projectIDs, p.ID)
		}
	}

	from := time.Now()
	to := from.AddDate(0, 0, days)

	return s.Repo.GetUpcomingDeadlines(from, to, projectIDs, deadlineType)
}
