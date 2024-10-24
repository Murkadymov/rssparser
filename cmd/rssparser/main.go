package main

import (
	"context"
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
	_, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	cfg := config.MustLoad()

	db, err := postgres.ConnectToDB(cfg)
	if err != nil {
		slog.Error("repository.ConnectToDB: ", "error", err.Error())
		return
	}

	_, err = db.Exec(
		`BEGIN TRANSACTION;
      
		CREATE TABLE IF NOT EXISTS feed(
    	feed_id SERIAL PRIMARY KEY,
    	feed_link VARCHAR(60) UNIQUE);
    	
    	INSERT INTO feed(
    	    feed_link
    	)
		VALUES
		('https://habr.com/ru/rss/all/all/'),
		('https://dtf.ru/rss/')
		ON CONFLICT DO NOTHING;
		
		END TRANSACTION;
`)
	if err != nil {
		slog.Error("error executing migration", "error", err.Error())
	}
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
