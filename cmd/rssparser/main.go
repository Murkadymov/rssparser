package main

import (
	"context"
	"fmt"
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

	doneChannel := make(chan struct{})

	cacheWorkerWG := &sync.WaitGroup{}
	repository := postgres.NewRepository(db)
	cache := feedcache.NewCache[string]()
	cacheWorker := rssparser.NewCacheWorker(cache, repository, cacheWorkerWG, time.Duration(cfg.WorkerInterval), doneChannel)
	feedWorker := rssparser.NewFeedWorker(cache, repository, 2, doneChannel)

	tx, err := db.Begin()

	defer func() {
		if p := recover(); p != nil {
			if err = tx.Rollback(); err != nil {
				slog.Error("error query rollback", "error", err.Error())
			}
			panic(p)

		} else if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				slog.Error("error rollbacking transaction")
			} else {
				slog.Info("transaction has been rolled back")
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				slog.Error("error commiting transaction")
				return
			} else {
				slog.Info("transaction commited succesfuly")
			}

		}

		slog.Info("db has been created")
	}()

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS feed(
		id SERIAL PRIMARY KEY,
	 	feed_link VARCHAR(255) UNIQUE);
	`)
	if err != nil {
		slog.Error("error creating table", "error", err.Error())
		return
	}

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS feed_content(
    item_id SERIAL PRIMARY KEY,
    feed_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    is_read BOOL DEFAULT FALSE NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feed(id) ON DELETE CASCADE);
`)
	if err != nil {
		slog.Error("error creating table", "error", err.Error())
		return
	}

	insertQuery := `INSERT INTO feed(feed_link) 
					VALUES ('https://habr.com/ru/rss/all/all/'),
						('https://dtf.ru/rss/'),
						('https://www.it-world.ru/tech/products/rss/'),
						('https://www.techcrunch.com/feed/'),
						('https://www.theverge.com/rss/index.xml'),
						('https://www.engadget.com/rss.xml'),
						('https://www.cnet.com/rss/all/'),
						('https://www.mashable.com/feed/'),
						('https://www.zdnet.com/news/rss.xml'),
						('https://www.feeds.arstechnica.com/arstechnica/index/'),
						('http://www.rss.slashdot.org/Slashdot/slashdotMain'),
						('https://www.news.ycombinator.com/rss'),
						('https://www.wired.com/feed/rss'),
						('https://www.itc.ua/feed/'),
						('https://www.computerworld.com/index.rss'),
						('https://www.readwrite.com/feed/'),
						('https://www.itpro.co.uk/feeds/all'),
						('https://www.digitaltrends.com/feed/'),
						('https://www.infoworld.com/index.rss')
					ON CONFLICT DO NOTHING;
					`
	_, err = tx.Exec(insertQuery)
	if err != nil {
		slog.Error("error inserting into table during transaction", "error", err.Error())
		return
	}
	_, err = fmt.Println("privet")
	if err != nil {
		fmt.Println("ne raven nil")
	}

	cacheWorker.RunCacheWorker(ctx, logger)

	err = feedWorker.RunFeedWorker(ctx, logger)
	if err != nil {
		return
	}

	<-ctx.Done()

	cacheWorker.CacheWorkerWG.Wait()

	slog.Info("service stopped")

	//fp := gofeed.NewParser()
	//
	//feedsURL := []string{
	//	"https://habr.com/ru/rss/all/all/",
	//	"https://dtf.ru/rss/",
	//}
	//
	//<-ctx.Done()
	//
	//wp.StopPool()
	//
	//wp.Wait()

}
