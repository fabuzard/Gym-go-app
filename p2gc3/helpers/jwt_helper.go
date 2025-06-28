package helper

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func ExtractUserID(c echo.Context) (uint, error) {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return 0, errors.New("unauthorized from helper")
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		return 0, errors.New("invalid token")
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user_id in token")
	}
	return uint(userIDFloat), nil
}
