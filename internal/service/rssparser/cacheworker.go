package rssparser

import (
	"context"
	"log/slog"
	"rssparser/internal/repository/interfaces"
	"time"
)

type Cache interface {
	Get() []string
	Set(items []string)
}

type CacheWorker struct {
	cache      Cache
	repository interfaces.Repository
	interval   time.Duration
}

func NewCacheWorker(cache Cache) *CacheWorker {
	return &CacheWorker{
		cache: cache,
	}
}

func (c *CacheWorker) GetFeedLinks(ctx context.Context, log slog.Logger) []string {
	op := "cache_worker.GetFeedLinks"
	feedLinks, err := c.repository.GetFeedURLs(ctx)
	if err != nil {

	}
}
