package services

import (
	"fmt"
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// ActivityService handles the business logic for activities.
type ActivityService struct {
	Repo           *storage.ActivityRepository
	ProjectService *ProjectService // To get user's projects
}

// NewActivityService creates a new instance of ActivityService.
func NewActivityService(repo *storage.ActivityRepository, projectService *ProjectService) *ActivityService {
	return &ActivityService{
		Repo:           repo,
		ProjectService: projectService,
	}
}

// CreateActivity creates and stores a new activity.
func (s *ActivityService) CreateActivity(activityType, entityType string, entityID, userID, projectID uint, description string) error {
	activity := &models.Activity{
		Type:        activityType,
		UserID:      userID,
		EntityType:  entityType,
		EntityID:    entityID,
		ProjectID:   projectID,
		Description: description,
	}
	return s.Repo.CreateActivity(activity)
}

// GetRecentActivities retrieves recent activities for a user.
func (s *ActivityService) GetRecentActivities(requestingUserID uint, projectID *uint, filterUserID *uint, limit int) ([]models.Activity, error) {
	var projectIDs []uint

	if projectID != nil {
		// If a specific project is requested, check if the user is a member.
		isMember, err := s.ProjectService.IsUserMemberOfProject(requestingUserID, *projectID)
		if err != nil {
			return nil, fmt.Errorf("could not verify project membership")
		}
		if !isMember {
			return nil, fmt.Errorf("user is not a member of the requested project")
		}
		projectIDs = append(projectIDs, *projectID)
	} else {
		// If no specific project, get all projects the user is a member of.
		projects, err := s.ProjectService.GetProjectsByUserID(requestingUserID)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve user's projects")
		}
		for _, p := range projects {
			projectIDs = append(projectIDs, p.ID)
		}
	}

	if len(projectIDs) == 0 {
		return []models.Activity{}, nil // User is not in any projects, so no activities to show.
	}

	return s.Repo.GetRecentActivities(projectIDs, filterUserID, limit)
}
