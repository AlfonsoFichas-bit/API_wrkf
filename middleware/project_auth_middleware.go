package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
)

// ProjectRoleAuth crea un middleware que comprueba si un usuario tiene un rol específico dentro de un proyecto.
// Requiere que se inyecte el ProjectService para comprobar el rol del usuario en la base de datos.
func ProjectRoleAuth(ps *services.ProjectService, requiredRoles ...models.ProjectRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Paso 1: Obtener el ID de usuario del contexto del JWT.
			userIDClaim := c.Get("userID")
			if userIDClaim == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "ID de usuario no encontrado en el token"})
			}
			// El claim numérico de JWT se decodifica como float64, hay que convertirlo de forma segura.
			userID, err := strconv.ParseUint(fmt.Sprintf("%.f", userIDClaim), 10, 64)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Formato de ID de usuario inválido en el token"})
			}

			// Paso 2: Obtener el ID del proyecto del parámetro de la URL.
			// Se busca "projectId" y, como alternativa, "id" para cubrir rutas como /projects/:id.
			projectIDStr := c.Param("projectId")
			if projectIDStr == "" {
				projectIDStr = c.Param("id") 
			}
			if projectIDStr == "" {
				// Si no se encuentra el ID del proyecto, podría ser una ruta anidada más profunda (ej. /sprints/:sprintId/tasks).
				// Por ahora, este middleware solo funciona para rutas que contienen directamente el ID del proyecto.
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID de proyecto no encontrado en la URL"})
			}
			projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Formato de ID de proyecto inválido en la URL"})
			}

			// Paso 3: Obtener el rol real del usuario en el proyecto desde el servicio.
			actualRole, err := ps.GetUserRoleInProject(uint(userID), uint(projectID))
			if err != nil {
				// Esto puede ser un error de la base de datos o "registro no encontrado".
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Prohibido: No eres miembro de este proyecto"})
			}

			// Paso 4: Comprobar si el rol del usuario es uno de los roles requeridos.
			hasPermission := false
			for _, requiredRole := range requiredRoles {
				if models.ProjectRole(actualRole) == requiredRole {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": fmt.Sprintf("Prohibido: La acción requiere uno de los siguientes roles: %v. Tu rol es: '%s'", requiredRoles, actualRole),
				})
			}

			// Paso 5: El usuario tiene permiso, proceder al siguiente manejador.
			return next(c)
		}
	}
}
