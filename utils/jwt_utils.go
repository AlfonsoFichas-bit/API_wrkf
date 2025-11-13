package utils

import (
	"errors"

	"github.com/labstack/echo/v4"
)

// GetUserIDFromContext extrae el ID de usuario del token JWT en el contexto Echo.
func GetUserIDFromContext(c echo.Context) (uint, error) {
	userID, ok := c.Get("userID").(float64)
	if !ok {
		return 0, errors.New("invalid or missing user ID in context")
	}

	return uint(userID), nil
}
