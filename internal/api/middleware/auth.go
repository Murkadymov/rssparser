package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"rssparser/internal/config"
)

func AuthMiddleware(next echo.HandlerFunc, cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {

		authID := c.Request().Header.Get("X-Auth-ID")

		switch authID {
		case "":
			return c.JSON(
				http.StatusForbidden,
				echo.Map{"message": "forbidden"},
			)
		case cfg.UserAuth:
			return c.JSON(
				http.StatusForbidden,
				echo.Map{"message": "forbidden"},
			)
		case cfg.AdminAuth:
			err := next(c)
			if err != nil {
				return fmt.Errorf("failed to call next handler: %w", err)
			}
		}
		return nil
	}
}
