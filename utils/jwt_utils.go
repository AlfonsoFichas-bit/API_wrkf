package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GetUserIDFromContext extrae el ID de usuario del token JWT en el contexto Echo.
func GetUserIDFromContext(c echo.Context) (uint, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return 0, errors.New("falta el token JWT o no es válido")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("failed to parse JWT claims")
	}

	sub, ok := claims["sub"].(float64) // JWT standard claim for subject (user ID)
	if !ok {
		return 0, errors.New("no se ha encontrado el ID de usuario en las reclamaciones JWT")
	}

	return uint(sub), nil
}
