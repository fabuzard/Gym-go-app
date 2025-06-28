package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.Println("JWT Missing Bearer prefix")
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"message": "Authentication failed",
					"details": "Authorization header must start with 'Bearer ' followed by a valid token",
				})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
						"message": "Authentication failed",
						"details": "Unexpected or unsupported signing method in token",
					})
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				log.Printf("JWT Parsing error: %v\n", err)
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"message": "Authentication failed",
					"details": "The token is invalid, malformed, or has expired",
				})
			}

			c.Set("user", token)
			return next(c)
		}
	}
}
