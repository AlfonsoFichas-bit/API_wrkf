package middleware

import (
	"net/http"

	"github.com/buga/API_wrkf/models"
	"github.com/labstack/echo/v4"
)

// TeacherAuthMiddleware checks if the authenticated user has the 'admin' or 'teacher' role.
// It must be used AFTER the JWTAuthMiddleware.
func TeacherAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole, ok := c.Get("userRole").(string)
		if !ok {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Forbidden: Role not found in token"})
		}

		if userRole != string(models.RoleAdmin) && userRole != string(models.RoleTeacher) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Forbidden: Admin or Teacher access required"})
		}

		return next(c)
	}
}
