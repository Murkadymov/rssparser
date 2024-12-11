package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"rssparser/internal/config"
	"rssparser/internal/pkg/log"
	"rssparser/internal/repository/feedcache"
	"rssparser/internal/repository/postgres"
	"rssparser/internal/service/rssparser"
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

	postgres.MustStartDB(db)

	doneChannel := make(chan struct{})
	cacheWorkerWG := &sync.WaitGroup{}

	repository := postgres.NewRepository(db)
	cache := feedcache.NewCache[string]()
	cacheWorker := rssparser.NewCacheWorker(cache, repository, cacheWorkerWG, time.Duration(cfg.WorkerInterval), doneChannel)
	feedWorker := rssparser.NewFeedWorker(cache, repository, 2, doneChannel)

	cacheWorker.RunCacheWorker(ctx, logger)

	err = feedWorker.RunFeedWorker(ctx, logger)
	if err != nil {
		return
	}

	<-ctx.Done()

	cacheWorker.CacheWorkerWG.Wait()

	slog.Info("service stopped")

}
