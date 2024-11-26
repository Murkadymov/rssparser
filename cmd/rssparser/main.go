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
		slog.E
	}

	_, err = tx.Exec(`	INSERT INTO feed(feed_link) 
	VALUES ('https://habr.com/ru/rss/all/all/'),
		('https://dtf.ru/rss/')
	ON CONFLICT DO NOTHING;
		`)

	if err != nil {
		slog.Error("error inserting into table")
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
