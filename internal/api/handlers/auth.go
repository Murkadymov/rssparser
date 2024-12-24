package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"rssparser/internal/api/responses"
	"rssparser/internal/models/api"
	"strconv"
)

func (h *Handler) AddUser(c echo.Context) error {
	var user *api.User

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.Error(err, "failed decoding json body"),
		)
	}
	defer func() {
		if err := c.Request().Body.Close(); err != nil {
			h.log.Error("request body close", "error", err)
		}
	}()

	userID, err := h.authService.AddUser(user)
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
				"userID":  strconv.Itoa(*userID),
			},
		),
	)
}

func (h *Handler) Login(c echo.Context) error {
	var user *api.User

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			responses.Error(err, nil),
		)
	}

	if err := h.authService.Login(user); err != nil {
		return c.JSON(
			http.StatusUnauthorized,
			responses.Error(
				err,
				"incorrect password",
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		responses.OK("successful login"),
	)
}
