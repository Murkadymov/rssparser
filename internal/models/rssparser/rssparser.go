package rssparser

import "sync"

type RSSCache struct {
	FeedURL []string
	mu      *sync.RWMutex
}

func NewRSSCache(initialURLs []string) *RSSCache {
	return &RSSCache{
		FeedURL: initialURLs,
	}
}

func (r *RSSCache) UpdateCache(urls []string) {
	defer r.mu.Unlock()
	r.mu.Lock()
	r.FeedURL = urls
}

func (r *RSSCache) GetCache() []string {
	defer r.mu.Unlock()
	r.mu.Lock()
	return r.FeedURL
}
