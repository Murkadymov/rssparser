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
			"error", err,
		)
		return
	}

	postgres.MustStartDB(db, logger)

	doneChannel := make(chan struct{})
	cacheWorkerWG := &sync.WaitGroup{}

	feedRepository := postgres.NewFeedRepository(db)
	HTTPRepository := postgres.NewHTTPRepository(db)

	feedServiceHTTP := feed.NewService(HTTPRepository, logger)
	authService := feed.NewAuthService(HTTPRepository, logger)

	feedHandlersHTTP := handlers.NewFeedHandlers(feedServiceHTTP, logger)
	authHandlers := handlers.NewAuthHandler(authService, logger)

	cacheFeed := cache.NewCache[string]()
	cacheWorker := feed.NewCacheWorker(
		cacheFeed,
		feedRepository,
		cacheWorkerWG,
		time.Duration(cfg.WorkerInterval),
		doneChannel,
	)

	feedItemChannel := make(chan *feed.FeedTask)

	feedWorker := feed.NewFeedWorker(
		cacheFeed,
		feedRepository,
		2,
		doneChannel,
		feedItemChannel,
	)

	e := echo.New()

	e.POST("/feed", middleware.AuthMiddleware(feedHandlersHTTP.InsertFeedService, cfg))
	e.POST("/feed/register", authHandlers.AddUser)
	e.POST("/feed/login", authHandlers.Login)

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
