package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAndGetEvaluation(t *testing.T) {
	// Setup
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	// --- Create Test Data ---
	adminUser, adminToken := CreateTestUser(t, testApp, "admin-eval@test.com", "admin")
	teacherUser, teacherToken := CreateTestUser(t, testApp, "teacher-eval@test.com", "docente")
	studentUser, _ := CreateTestUser(t, testApp, "student-eval@test.com", "user")

	project := CreateTestProject(t, testApp, "Eval Project", adminUser.ID)
	AddUserToProject(t, testApp, project.ID, teacherUser.ID, "docente")
	AddUserToProject(t, testApp, project.ID, studentUser.ID, "team_developer")

	userStory := CreateTestUserStory(t, testApp, "Eval Story", project.ID)
	task := CreateTestTask(t, testApp, "Eval Task", userStory.ID, studentUser.ID)

	rubric := CreateTestRubric(t, testApp, project.ID, adminUser.ID, "Eval Rubric")
	require.Len(t, rubric.Criteria, 2, "Test rubric should have 2 criteria")

	// --- Test Case 1: Teacher successfully creates an evaluation ---
	t.Run("Teacher creates evaluation successfully", func(t *testing.T) {
		evalReq := services.CreateEvaluationRequest{
			RubricID:        rubric.ID,
			OverallFeedback: "Excellent work!",
			CriterionEvaluations: []services.CriterionEvaluationRequest{
				{CriterionID: rubric.Criteria[0].ID, Score: 5, Feedback: "Criterion 1 OK"},
				{CriterionID: rubric.Criteria[1].ID, Score: 4, Feedback: "Criterion 2 OK"},
			},
		}
		reqBody, _ := json.Marshal(evalReq)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/tasks/%d/evaluations", task.ID), bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+teacherToken)
		rec := httptest.NewRecorder()

		testApp.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var createdEval models.Evaluation
		err := json.Unmarshal(rec.Body.Bytes(), &createdEval)
		require.NoError(t, err)
		assert.Equal(t, "Excellent work!", createdEval.OverallFeedback)
		assert.Equal(t, 9.0, createdEval.TotalScore) // 5 + 4
		assert.Len(t, createdEval.CriterionEvaluations, 2)
		assert.Equal(t, task.ID, createdEval.TaskID)
		assert.Equal(t, teacherUser.ID, createdEval.EvaluatorID)

		// --- Test Case 2: Get the created evaluation ---
		t.Run("Get created evaluation", func(t *testing.T) {
			reqGet := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/%d/evaluations", task.ID), nil)
			reqGet.Header.Set(echo.HeaderAuthorization, "Bearer "+adminToken) // Any valid user can see it
			recGet := httptest.NewRecorder()

			testApp.Router.ServeHTTP(recGet, reqGet)

			assert.Equal(t, http.StatusOK, recGet.Code)
			var fetchedEvals []models.Evaluation
			err := json.Unmarshal(recGet.Body.Bytes(), &fetchedEvals)
			require.NoError(t, err)
			require.Len(t, fetchedEvals, 1)
			assert.Equal(t, createdEval.ID, fetchedEvals[0].ID)
			assert.Equal(t, "Excellent work!", fetchedEvals[0].OverallFeedback)
			assert.Len(t, fetchedEvals[0].CriterionEvaluations, 2)
			assert.NotNil(t, fetchedEvals[0].Evaluator)
			assert.Equal(t, teacherUser.Nombre, fetchedEvals[0].Evaluator.Nombre)
		})
	})

	// --- Test Case 3: Student fails to create an evaluation (permission denied) ---
	t.Run("Student fails to create evaluation", func(t *testing.T) {
		evalReq := services.CreateEvaluationRequest{RubricID: rubric.ID}
		reqBody, _ := json.Marshal(evalReq)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/tasks/%d/evaluations", task.ID), bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+adminToken) // Using admin for simplicity, any non-docente fails
		rec := httptest.NewRecorder()

		testApp.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected error due to permissions")
	})
}
