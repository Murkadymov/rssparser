package rssparser

import (
	"context"
	"time"
)

type Cache interface {
	Get() []string
	Set(items []string)
}

type Repository interface {
	InsertFeedURLs(ctx context.Context, feedURLs []string) error
	GetFeedURLs(ctx context.Context) []string
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

func (fw *RSSParser) LoadFeedURLsToCache(ctx context.Context) {

	feedURLs := fw.repository.GetFeedURLs(ctx)

}
