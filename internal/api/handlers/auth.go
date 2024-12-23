package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"rssparser/internal/api/responses"
	"rssparser/internal/models/api"
	"strconv"
)

type AuthService interface {
	AddUser(user *api.User) (*int, error)
}

type AuthHandler struct {
	authService AuthService
	authLogger  *slog.Logger
}

func NewAuthHandler(authService AuthService, authLogger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		authLogger:  authLogger,
	}
}

func (a *AuthHandler) AddUser(c echo.Context) error {
	var user *api.User

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.Error(err, "failed decoding json body"),
		)
	}

	userID, err := a.authService.AddUser(user)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			responses.Error(err, "failed to add user"),
		)
	}

	return c.JSON(
		http.StatusOK,
		responses.OK(
			map[string]string{
				"message": "user added",
				"userID":  strconv.Itoa(*userID)},
		),
	)
}
