package models

// --- Platform Roles ---

// PlatformRole represents a user's role at the application level.
type PlatformRole string

const (
	RoleAdmin PlatformRole = "admin"
	RoleUser  PlatformRole = "user"
)

// --- Project Roles ---

// ProjectRole represents a user's role within a specific project.
type ProjectRole string

const (
	RoleScrumMaster   ProjectRole = "scrum_master"
	RoleProductOwner  ProjectRole = "product_owner"
	RoleTeamDeveloper ProjectRole = "team_developer"
)

// IsValid checks if the project role is a defined valid role.
func (r ProjectRole) IsValid() bool {
	switch r {
	case RoleScrumMaster, RoleProductOwner, RoleTeamDeveloper:
		return true
	default:
		return false
	}
}
