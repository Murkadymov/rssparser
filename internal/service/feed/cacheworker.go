package feed

import (
	"context"
	"fmt"
	"log/slog"
	"rssparser/internal/repository/interfaces"
	"sync"
	"time"
)

type CacheWorker struct {
	cache         interfaces.Cache
	repository    interfaces.Repository
	interval      time.Duration
	CacheWorkerWG *sync.WaitGroup
	DoneChannel   chan<- struct{}
}

func NewCacheWorker(
	cache interfaces.Cache,
	repository interfaces.Repository,
	CacheWorkerWG *sync.WaitGroup,
	interval time.Duration,
	doneChan chan struct{}) *CacheWorker {
	return &CacheWorker{
		cache:         cache,
		repository:    repository,
		CacheWorkerWG: CacheWorkerWG,
		interval:      interval,
		DoneChannel:   doneChan,
	}
}

func (c *CacheWorker) GetFeedURLs(ctx context.Context, log *slog.Logger) ([]string, error) {
	op := "cacheworker.GetFeedURLs"

	feedLinks, err := c.repository.GetFeedURLs(ctx)
	if err != nil {
		slog.Error(
			"failed receiving feed urls",
			"data", fmt.Sprintf("%s: repository.GetFeedURLs: %s", op, err.Error()))

		return nil, fmt.Errorf("%s: repository.GetFeedURLs: %s", op, err.Error())
	}

	slog.Info("success receiving urls", "op", op, "link list", feedLinks)
	return feedLinks, nil
}

func (c *CacheWorker) SetFeedURL(ctx context.Context, log *slog.Logger, feedURLs []string) {
	op := "cacheworker.SetFeedURL"

	fmt.Printf("received urls: \n%s ", feedURLs)

	c.cache.Set(feedURLs)

	slog.Info("url successfully added", slog.String("method", op))
}

func (c *CacheWorker) UpdateCache(ctx context.Context, log *slog.Logger) error {
	op := "cacheworker.UpdateCache"

	defer c.CacheWorkerWG.Done()

	feedLinks, err := c.GetFeedURLs(ctx, log)
	if err != nil {
		slog.Error(
			"failed to update feedcache",
			"caller", op,
			"feedcache.GetFeedURLs", err.Error(),
		)
	}

	fmt.Printf("received urls: \n%s ", feedLinks)

	c.cache.Update(feedLinks)
	log.Info("cache has been updated")

	return nil
}

func (c *CacheWorker) RunCacheWorker(ctx context.Context, log *slog.Logger) {
	op := "cacheworker.RunCacheWorker"
	go func() {
		cacheWorkerTicker := time.NewTicker(c.interval * time.Second)
		for {
			select {
			case <-cacheWorkerTicker.C:
				c.CacheWorkerWG.Add(1)
				if err := c.UpdateCache(ctx, log); err != nil {
					slog.Error(
						"error occured running cacheWorker",
						"method", op,
						"error", err.Error(),
					)
					return
				}
				slog.Info("success updating cache", "method", op)

				c.SendDoneSignal()
			case <-ctx.Done():
				slog.Info("worker stopped", "method", op)

				cacheWorkerTicker.Stop()

				return
			}
		}
	}()

	fmt.Println("CACHEWORKER STARTED")
}

func (c *CacheWorker) SendDoneSignal() {
	c.DoneChannel <- struct{}{}
}
func (c *CacheWorker) CloseDoneChannel() {
	close(c.DoneChannel)
}
