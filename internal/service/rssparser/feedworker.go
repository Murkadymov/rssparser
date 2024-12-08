package rssparser

import (
	"context"
	"fmt"
	"github.com/mmcdole/gofeed"
	"log/slog"
	"rssparser/internal/repository/interfaces"
)

type FeedWorker struct {
	chunkSize      int
	cache          interfaces.Cache
	feedparser     *gofeed.Parser
	repo           interfaces.Repository
	workerDoneChan <-chan struct{}
}

func NewFeedWorker(
	cache interfaces.Cache,
	repo interfaces.Repository,
	batchPartition int,
	workerDoneChan <-chan struct{},
) *FeedWorker {
	return &FeedWorker{
		chunkSize:      batchPartition,
		cache:          cache,
		feedparser:     gofeed.NewParser(),
		repo:           repo,
		workerDoneChan: workerDoneChan,
	}
}

func (fw *FeedWorker) FetchFeedLinks(ctx context.Context, log *slog.Logger) error {
	op := "service.FetchFeedLinks"

	feedLinks := fw.cache.Get()

	for low := 0; low < len(feedLinks); low += fw.chunkSize {
		high := low + fw.chunkSize

		if high > len(feedLinks) {
			high = len(feedLinks)
		}

		for _, link := range feedLinks[low:high] {
			feed, err := fw.feedparser.ParseURL(link)
			if err != nil {
				log.Error(
					"error parsing url",
					"method", op,
					"fn", "feedparser.ParseURL",
					"error", err.Error(),
				)
			}

			err = fw.repo.InsertFeedURLs(ctx, feed)
			if err != nil {
				return err
				//TODO::
			}
		}
	}

	log.Info("done fetching feed links")

	return nil
}

func (fw *FeedWorker) RunFeedWorker(ctx context.Context, log *slog.Logger) error {
	op := "service.RunFeedWorker"

	go func() {
		for {
			select {
			case <-fw.workerDoneChan:
				err := fw.FetchFeedLinks(ctx, log)
				if err != nil {
					panic(err)
				}

				fmt.Println("FETCH FEED LINKS STARTED", op)
			}
		}
	}()
	return nil
}
