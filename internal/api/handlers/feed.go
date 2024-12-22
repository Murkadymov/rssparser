package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"rssparser/internal/models/api"
)

type FeedService interface {
	InsertFeedSource(ctx context.Context, feedSource *api.FeedSource) error
}

type FeedHandlers struct {
	feedService FeedService
}

func NewFeedHandlers(service FeedService) *FeedHandlers {
	return &FeedHandlers{
		feedService: service,
	}
}

func (h *FeedHandlers) InsertFeedService(c echo.Context) error {
	var feedSource *api.FeedSource

	if c.Request().Method != "POST" {
		return echo.NewHTTPError(
			http.StatusMethodNotAllowed,
			"Method not allowed",
		)
	}

	c.Response().Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(c.Request().Body).Decode(&feedSource); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			h.error(err, "reading decoding body"),
		)
	}

	ctx := c.Request().Context()
	//ctx, cancel := context.WithTimeout(ctx, 15) // TODO: timeout

	if err := h.feedService.InsertFeedSource(ctx, feedSource); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			h.error(err, nil),
		)
	}

	return c.JSON(
		http.StatusOK,
		h.ok(fmt.Sprintf(
			"inserted link: %s",
			feedSource.FeedLink),
		),
	)
}
