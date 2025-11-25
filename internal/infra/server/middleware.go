package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ecelayes/pms-core/pkg/auth"
	"github.com/ecelayes/pms-core/pkg/response"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, http.StatusUnauthorized, echo.ErrUnauthorized)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Error(c, http.StatusUnauthorized, echo.ErrUnauthorized)
		}

		tokenString := parts[1]

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			return response.Error(c, http.StatusUnauthorized, err)
		}

		c.Set("user_id", claims.UserID)
		c.Set("hotel_id", claims.HotelID)
		c.Set("role", claims.Role)

		return next(c)
	}
}
