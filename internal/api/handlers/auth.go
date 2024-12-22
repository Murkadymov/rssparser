package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"rssparser/internal/models/api"
)

func (h *FeedHandlers) Register(c echo.Context) error {
	var user *api.User

	json.NewDecoder(c.Request().Body).Decode(&user)

}
