package main

import (
	"context"
	"github.com/labstack/echo/v4"
	log2 "log"
	"log/slog"
	"os"
	"os/signal"
	"rssparser/internal/api/handlers"
	"rssparser/internal/api/middleware"
	"rssparser/internal/config"
	"rssparser/internal/repository/cache"
	"rssparser/internal/repository/postgres"
	"rssparser/internal/service/feed"
	"rssparser/pkg/log"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	logger := log.New()
	cfg := config.MustLoad()

	db, err := postgres.ConnectToDB(cfg)
	if err != nil {
		slog.Error(
			"repository.ConnectToDB: ",
			"error", err.Error())
		return
	}

	postgres.MustStartDB(db, logger)

	doneChannel := make(chan struct{})
	cacheWorkerWG := &sync.WaitGroup{}

	repository := postgres.NewRepository(db)
	service := feed.NewService(repository, logger)
	handlers := handlers.NewFeedHandlers(service)

	cache := cache.NewCache[string]()
	cacheWorker := feed.NewCacheWorker(
		cache,
		repository,
		cacheWorkerWG,
		time.Duration(cfg.WorkerInterval),
		doneChannel,
	)

	feedItemChannel := make(chan *feed.FeedTask)

	feedWorker := feed.NewFeedWorker(
		cache,
		repository,
		2,
		doneChannel,
		feedItemChannel,
	)

	e := echo.New()

	e.POST("/feed", middleware.AuthMiddleware(handlers.InsertFeedService, cfg))
	e.POST("/feed/register")

	go func() {
		if err := e.Start("localhost:8080"); err != nil {
			log2.Fatal("error running server")
		}
	}()

	cacheWorker.RunCacheWorker(ctx, logger)

	feedWorker.RunFeedFetch(ctx, logger)

	feedWorker.RunFeedWorkers(ctx, 2, logger)

	<-ctx.Done()
	close(doneChannel)
	close(feedItemChannel)
	cacheWorker.CacheWorkerWG.Wait()

	slog.Info("service stopped")

}
