package feed

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/mmcdole/gofeed"
	log2 "log"
	"log/slog"
	"net/url"
	"rssparser/internal/repository/interfaces"
	"rssparser/internal/utils"
	"strings"
	"sync"
	"time"
)

type FeedWorker struct {
	chunkSize      int
	cache          interfaces.Cache
	feedparser     *gofeed.Parser
	repo           interfaces.FeedRepository
	workerDoneChan <-chan struct{}
	workerItemChan chan *FeedTask
}

type FeedTask struct {
	primaryFeedID int
	feedItem      *gofeed.Item
}

func NewFeedWorker(
	cache interfaces.Cache,
	repo interfaces.FeedRepository,
	batchPartition int,
	workerDoneChan <-chan struct{},
	workerItemChan chan *FeedTask,
) *FeedWorker {
	return &FeedWorker{
		chunkSize:      batchPartition,
		cache:          cache,
		feedparser:     gofeed.NewParser(),
		repo:           repo,
		workerDoneChan: workerDoneChan,
		workerItemChan: workerItemChan,
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
					"link", link,
				)
				continue
			}
			urlParsed, err := url.Parse(strings.TrimSpace(feed.Link))
			if err != nil {

				log.Error("error parsing url", op, err.Error())

				return fmt.Errorf("%s: %w", op, err)
			}
			log.Error("GOT ", "url", urlParsed.Host)
			feedPrimaryID, err := fw.repo.GetLinkPrimaryID(ctx, urlParsed.Host)
			if err != nil {
				log.Error(
					"get primary id",
					op, "repo.GetLinkPrimaryID",
					"errorText", err,
				)
				if feedPrimaryID == 0 {
					log.Error("GOT 0 ID", "url", urlParsed.Host)
				}
				return fmt.Errorf("%s: %w", op, err)
			}

			for _, feedItem := range feed.Items {

				select {
				case fw.workerItemChan <- &FeedTask{
					primaryFeedID: feedPrimaryID,
					feedItem:      feedItem,
				}:

				case <-ctx.Done():
					log.Debug("")
					return nil

				}
			}
		}

	}
	log.Debug("done fetching feed links")

	return nil
}

func (fw *FeedWorker) ProcessFeedItem(
	ctx context.Context,
	feedPrimaryID int,
	feedItem *gofeed.Item,
) error {
	const op = "feedworker.ProcessFeedItem"

	var pubDate *time.Time = feedItem.PublishedParsed //вынес получение pubDate в service

	rawExistingDate, err := fw.repo.GetExistingPubDate(feedItem.Link)
	if err != nil {
		log.Error("failed to get existing pubdate", "errorText", err)
		fmt.Errorf("get existing pub date: %w ", err)
	}

	if rawExistingDate != "" {
		isNewer, err := fw.compareDates(pubDate, rawExistingDate)
		if !isNewer {
			log.Debug(
				"publication is up to date",
				"feedItemPubDate", pubDate,
				"dbPubDate", rawExistingDate,
			)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	feedItem.Description = utils.TextCleaner(feedItem.Description)

	err = fw.repo.InsertFeedContent(
		ctx,
		feedPrimaryID,
		feedItem.Title,
		feedItem.Description,
		pubDate,
		feedItem.Link,
	)
	if err != nil {
		log.Error(
			"method", op,
			"errorText", err,
			"link", feedItem.Link,
			"feedPrimaryID", feedPrimaryID,
			"feedDate", pubDate,
			"existingDate", rawExistingDate,
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (fw *FeedWorker) compareDates(date *time.Time, dateString string) (bool, error) {

	layout := "2006-01-02 15:04:05"

	parsedExistingDate, err := time.Parse(layout, dateString)
	if err != nil {
		return false, fmt.Errorf("parsing existing date: %w", err)
	}

	if isNewer := parsedExistingDate.Before(*date); !isNewer {
		return false, nil

	}

	return true, nil
}

func (fw *FeedWorker) RunFeedWorkers(ctx context.Context, workersNumber int, log *slog.Logger) error {
	op := "service.RunFeedWorkers"

	errChan := make(chan error, workersNumber)
	var workerErrors []error

	var wg sync.WaitGroup = sync.WaitGroup{}

	for i := range workersNumber { //for i:=0; i<workersNumber; i++
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case feedTask, ok := <-fw.workerItemChan:
					if !ok {

						return
					}

					log.Debug("worker received item",
						"worker id", workerID,
						"feed link", feedTask.feedItem,
						"method", op,
					)
					if err := fw.ProcessFeedItem(ctx, feedTask.primaryFeedID, feedTask.feedItem); err != nil {
						errChan <- err
					}
					log.Debug("worker done processing item",
						"worker id", workerID,
						"feed link", feedTask.feedItem.Link,
						"method", op,
					)
				case <-ctx.Done():

					return
				}
			}
		}(i)

	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			workerErrors = append(workerErrors, err)
		}
	}

	if len(workerErrors) > 0 {
		return fmt.Errorf("encountered  %d \n errors: %v\n", len(errChan), workerErrors)
	}

	return nil
}

func (fw *FeedWorker) RunFeedFetch(ctx context.Context, log *slog.Logger) {
	go func() {
		for {
			select {
			case <-fw.workerDoneChan:
				err := fw.FetchFeedLinks(ctx, log)
				if err != nil {
					log2.Fatalf("run feed: %w", err)
				}
			case <-ctx.Done():
				return
			}

		}
	}()

}
