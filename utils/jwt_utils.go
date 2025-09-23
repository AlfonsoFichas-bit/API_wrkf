package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GetUserIDFromContext extracts the user ID from the JWT token in the Echo context.
func GetUserIDFromContext(c echo.Context) (uint, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return 0, errors.New("JWT token missing or invalid")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("failed to parse JWT claims")
	}

	sub, ok := claims["sub"].(float64) // JWT standard claim for subject (user ID)
	if !ok {
		return 0, errors.New("user ID not found in JWT claims")
	}

	return uint(sub), nil
}
