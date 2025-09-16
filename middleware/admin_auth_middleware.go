package middleware

import (
	"API_wrkf/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// AdminAuthMiddleware checks if the authenticated user has the 'admin' role.
// It should be used AFTER the JWTAuthMiddleware.
func AdminAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole, ok := c.Get("userRole").(string)
		if !ok || userRole != string(models.RoleAdmin) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Forbidden: Administrator access required"})
		}
		return next(c)
	}
}
