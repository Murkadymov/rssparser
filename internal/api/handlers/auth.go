package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"rssparser/internal/api/responses"
	"rssparser/internal/models/api"
	"rssparser/internal/service/feed"
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

	token, err := h.authService.Login(h.secret, user)
	if err != nil {
		if err = errors.Unwrap(err); errors.Is(err, sql.ErrNoRows) {
			return c.JSON(
				http.StatusUnauthorized,
				responses.Error(
					err,
					ErrUserNotExist,
				),
			)
		}
		var jwtErr *feed.JWTGenerationError
		if errors.As(err, &jwtErr) {
			return c.JSON(
				http.StatusUnauthorized,
				responses.Error(
					err,
					nil,
				),
			)
		}

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
		responses.OK(echo.Map{
			"additionalText": "successful auth",
			"token":          token,
		}),
	)
}
