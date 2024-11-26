package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"rssparser/internal/config"
	"rssparser/internal/repository/postgres"
	"syscall"
)

func init() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				s := a.Value.Any().(*slog.Source)
				s.File = path.Base(s.File)
			}
			return a
		},
	}))
	slog.SetDefault(logger)

}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	cfg := config.MustLoad()

	db, err := postgres.ConnectToDB(cfg)
	if err != nil {
		slog.Error("repository.ConnectToDB: ", "error", err.Error())
		return
	}

	repository := postgres.NewRepository(db)

	tx, err := db.Begin()

	defer func() {
		if p := recover(); p != nil {
			if err = tx.Rollback(); err != nil {
				slog.Error("error rollbacking query", "error", err.Error())
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
		feed_id SERIAL PRIMARY KEY,
	 	feed_link VARCHAR(255) UNIQUE);
	`)
	if err != nil {
		slog.Error("error creating table", "error", err.Error())
		return
	}

	insertQuery := `INSERT INTO feed(feed_link) 
					VALUES ('https://habr.com/ru/rss/all/all/'),
						('https://dtf.ru/rss/'),
						('https://www.it-world.ru/tech/products/rss/'),
						('https://techcrunch.com/feed/'),
						('https://www.theverge.com/rss/index.xml'),
						('https://www.engadget.com/rss.xml'),
						('https://www.cnet.com/rss/all/'),
						('https://mashable.com/feed/'),
						('https://www.zdnet.com/news/rss.xml'),
						('https://feeds.arstechnica.com/arstechnica/index/'),
						('http://rss.slashdot.org/Slashdot/slashdotMain'),
						('https://news.ycombinator.com/rss'),
						('https://www.wired.com/feed/rss'),
						('https://itc.ua/feed/'),
						('https://www.techrepublic.com/rssfeeds/articles/latest/'),
						('https://www.computerworld.com/index.rss'),
						('https://readwrite.com/feed/'),
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

	urls := repository.GetFeedURLs(ctx)

	fmt.Println(urls)
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
