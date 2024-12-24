package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"rssparser/internal/api/responses"
	"rssparser/internal/models/api"
)

type FeedService interface {
	InsertFeedSource(ctx context.Context, feedSource *api.FeedSource) error
}
type FeedHandlers struct {
	feedService   FeedService
	feedAPILogger *slog.Logger
}

func NewFeedHandlers(service FeedService, feedAPILogger *slog.Logger) *FeedHandlers {
	return &FeedHandlers{
		feedService:   service,
		feedAPILogger: feedAPILogger,
	}
}

func (h *FeedHandlers) InsertFeedService(c echo.Context) error {
	var feedSource *api.FeedSource

	defer func() {
		if err := c.Request().Body.Close(); err != nil {
			h.feedAPILogger.Error("request body close", "error", err)
		}
	}()

	if c.Request().Method != "POST" {
		return echo.NewHTTPError(
			200,
			"Method not allowed",
		)
	}

	c.Response().Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(c.Request().Body).Decode(&feedSource); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			responses.Error(err, "decoding body into feedsource"),
		)
	}

	ctx := c.Request().Context()
	//ctx, cancel := context.WithTimeout(ctx, 15) // TODO: timeout

	if err := h.feedService.InsertFeedSource(ctx, feedSource); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			responses.Error(err, nil),
		)
	}

	return c.JSON(
		http.StatusOK,
		responses.OK(fmt.Sprintf(
			"inserted link: %s",
			feedSource.FeedLink),
		),
	)
}
