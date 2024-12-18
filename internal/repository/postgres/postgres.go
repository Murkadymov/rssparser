package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"rssparser/internal/config"
)

const rollbackError = "failed transaction rollback"

func ConnectToDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
			cfg.User,
			cfg.Password,
			cfg.DB,
			cfg.Host,
			cfg.Port,
			cfg.SSLMode,
		),
	)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	slog.Debug("successful connection to DB")

	return db, nil
}

func MustStartDB(db *sql.DB) {
	tx, err := db.Begin()

	defer func() {
		if p := recover(); p != nil {
			if err = tx.Rollback(); err != nil {
				slog.Error(rollbackError, "error", err.Error())
			}
			panic(p)

		} else if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				slog.Error(rollbackError)
			} else {
				slog.Info("success transaction rollback")
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				slog.Error("error transaction commit")
				return
			} else {
				slog.Info("success transaction commit")
			}

		}

		slog.Info("db has been created")
	}()

	const createFeedTableQuery = `CREATE TABLE IF NOT EXISTS feed(
									id SERIAL PRIMARY KEY,
									feed_link VARCHAR(255) UNIQUE);`

	result, err := tx.Exec(createFeedTableQuery)

	fmt.Println("RESULT: ", result)
	if err != nil {
		slog.Error("failed creating table", "error", err.Error())
		return
	}

	const createFeedContentQuery = `CREATE TABLE IF NOT EXISTS feed_content(
										item_id SERIAL PRIMARY KEY,
										feed_id INTEGER NOT NULL,
										title TEXT NOT NULL,
										description TEXT NOT NULL,
    									published_at TIMESTAMP NOT NULL,
    									pub_link TEXT NOT NULL,
										is_read BOOL DEFAULT FALSE NOT NULL,
				                 	FOREIGN KEY (feed_id) REFERENCES feed(id) ON DELETE CASCADE);`
	_, err = tx.Exec(createFeedContentQuery)
	if err != nil {
		slog.Error("failed creating table", "error", err.Error())
		return
	}

	const insertFeedLinksQuery = `INSERT INTO feed(feed_link) 
								  VALUES ('https://habr.com/ru/rss/all/all/'),
										 ('https://dtf.ru/rss/'),
										 ('https://www.it-world.ru/tech/products/rss/'),
										 ('https://www.theverge.com/rss/index.xml'),
										 ('https://www.engadget.com/rss.xml'),
										 ('https://www.cnet.com/rss/all/'),
										 ('https://www.zdnet.com/news/rss.xml'),
										 ('https://www.wired.com/feed/rss') ON CONFLICT DO NOTHING;`
	_, err = tx.Exec(insertFeedLinksQuery)
	if err != nil {
		slog.Error("failed to insert into table during transaction", "error", err.Error())
		return
	}

}
