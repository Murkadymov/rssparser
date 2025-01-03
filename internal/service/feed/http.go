package feed

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"log/slog"
	"rssparser/internal/models/api"
)

type HTTPService struct {
	repo       HTTPRepository
	feedLogger *slog.Logger
}

func NewService(repo HTTPRepository, log *slog.Logger) *HTTPService {
	return &HTTPService{
		repo:       repo,
		feedLogger: log,
	}
}

func (s *HTTPService) InsertFeedSource(
	ctx context.Context,
	feedSource *api.FeedSource,
) error {

	const op = "feed.InsertFeedSource"

	if err := s.repo.InsertFeedSource(ctx, feedSource.FeedLink); err != nil {
		log.Error(op, "errorText", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("feed source added", "link", feedSource.FeedLink)

	return nil
}
