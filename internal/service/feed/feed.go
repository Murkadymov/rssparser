package feed

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"log/slog"
	"rssparser/internal/models/api"
	"rssparser/internal/repository/interfaces"
)

type Service struct {
	repo interfaces.Repository
	log  *slog.Logger
}

func NewService(repo interfaces.Repository, log *slog.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

func (s *Service) InsertFeedSource(
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
