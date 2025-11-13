package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTAuthMiddleware crea y devuelve un middleware JWT que valida los tokens
// utilizando la clave secreta proporcionada.
func JWTAuthMiddleware(secret string) echo.MiddlewareFunc {
	jwtSecret := []byte(secret)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing authorization header"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header format"})
			}

			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if ok {
				c.Set("userID", claims["sub"])
				c.Set("userName", claims["nam"])
				c.Set("userRole", claims["rol"]) // <-- EXTRACT ROLE
			}

			return next(c)
		}
	}
}

// GetUserIDFromContext extracts the user ID from the Echo context.
func GetUserIDFromContext(c echo.Context) (uint, error) {
	userIDVal := c.Get("userID")
	if userIDVal == nil {
		return 0, errors.New("user ID not found in context")
	}

	// JWT claims are often float64, so we need to handle that.
	userIDFloat, ok := userIDVal.(float64)
	if !ok {
		return 0, errors.New("user ID in context is not of expected type")
	}

	return uint(userIDFloat), nil
}
