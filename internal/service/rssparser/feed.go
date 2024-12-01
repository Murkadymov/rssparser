package rssparser

import (
	"context"
	"fmt"
	"time"
)

type Cache interface {
	Get() []string
	Set(items []string)
}

type Repository interface {
	InsertFeedURLs(ctx context.Context, feedURLs []string) error
	GetFeedURLs(ctx context.Context) ([]string, error)
}

type RSSParser struct {
	cache      Cache
	repository Repository
	interval   time.Duration
}

func NewFetchWorker(cache Cache, repository Repository, interval time.Duration) *RSSParser {
	return &RSSParser{
		cache:      cache,
		repository: repository,
		interval:   interval,
	}
}

func (fw *RSSParser) LoadFeedURLsToCache(ctx context.Context) error {
	op := "service.LoadFeedURLsToCache"
	feedURLs, err := fw.repository.GetFeedURLs(ctx)

	if err != nil {
		return fmt.Errorf("%s: repository.GetFeedURLs: %w", op, err)
	}
	fmt.Errorf("error scanning into string: %s, %w", op, err)
}
