package middleware

import (
	"net/http"

	"github.com/buga/API_wrkf/models"

	"github.com/labstack/echo/v4"
)

// AdminAuthMiddleware comprueba si el usuario autenticado tiene el rol «admin».
// Debe utilizarse DESPUÉS de JWTAuthMiddleware.
func AdminAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole, ok := c.Get("userRole").(string)
		if !ok || userRole != string(models.RoleAdmin) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Prohibido: Se requiere acceso de admin"})
		}
		return next(c)
	}
}
