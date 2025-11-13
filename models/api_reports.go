package models

// VelocityReport represents the team's velocity for a project.
type VelocityReport struct {
	ProjectID         uint             `json:"project_id"`
	AverageVelocity   float64          `json:"average_velocity"`
	SprintsConsidered int              `json:"sprints_considered"`
	VelocityPerSprint []SprintVelocity `json:"velocity_per_sprint"`
}

// SprintVelocity shows the points completed in a single sprint.
type SprintVelocity struct {
	SprintID        uint   `json:"sprint_id"`
	SprintName      string `json:"sprint_name"`
	CompletedPoints int    `json:"completed_points"`
}

// BurndownReport represents the data for a sprint's burndown chart.
type BurndownReport struct {
	SprintID      uint            `json:"sprint_id"`
	SprintName    string          `json:"sprint_name"`
	TotalPoints   int             `json:"total_points"`
	BurndownData  []BurndownPoint `json:"burndown_data"`
}

// BurndownPoint represents the remaining points on a specific day.
type BurndownPoint struct {
	Date            string  `json:"date"` // "YYYY-MM-DD"
	RemainingPoints float64 `json:"remaining_points"`
	IdealPoints     float64 `json:"ideal_points"`
}

// CommitmentReport represents the comparison between committed and completed points in a sprint.
type CommitmentReport struct {
	SprintID        uint    `json:"sprintId"`
	SprintName      string  `json:"sprintName"`
	CommittedPoints int     `json:"committedPoints"`
	CompletedPoints int     `json:"completedPoints"`
	CompletionRate  float64 `json:"completionRate"`
}
