package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"log/slog"
	"os"
	"os/signal"
	"rssparser/internal/api/handlers"
	"rssparser/internal/api/middleware"
	"rssparser/internal/config"
	"rssparser/internal/repository/feedcache"
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
		slog.Error("repository.ConnectToDB: ", "error", err.Error())
		return
	}

	postgres.MustStartDB(db, logger)

	doneChannel := make(chan struct{})
	cacheWorkerWG := &sync.WaitGroup{}

	repository := postgres.NewRepository(db)
	service := feed.NewFeedService(repository, logger)
	handlers := handlers.NewFeedHandlers(service)

	cache := feedcache.NewCache[string]()
	cacheWorker := feed.NewCacheWorker(
		cache,
		repository,
		cacheWorkerWG,
		time.Duration(cfg.WorkerInterval),
		doneChannel,
	)
	feedWorker := feed.NewFeedWorker(
		cache,
		repository,
		2,
		doneChannel,
	)

	e := echo.New()

	e.POST("/feed", middleware.AuthMiddleware(handlers.InsertFeedService, cfg))

	go func() {
		if err := e.Start("localhost:8080"); err != nil {
			panic("error starting server")
		}
	}()

	cacheWorker.RunCacheWorker(ctx, logger)

	err = feedWorker.RunFeedWorker(ctx, logger)
	if err != nil {
		return
	}

	<-ctx.Done()

	cacheWorker.CacheWorkerWG.Wait()

	slog.Info("service stopped")

}
