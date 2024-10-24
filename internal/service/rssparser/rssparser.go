package rssparser

import "time"

type Cache interface {
	Get() []string
	Set(items []string)
}

type RSSParser struct {
	cache    Cache
	interval time.Duration
}

func NewFetchWorker(cache Cache, interval time.Duration) *RSSParser {
	return &RSSParser{
		cache:    cache,
		interval: interval,
	}
}

func (fw *RSSParser) FetchFeed() {

}
