package postgres

import (
	"context"
	"database/sql"
	"log/slog"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (d *Repository) InsertFeedURLs(ctx context.Context, feedURLs []string) error {

	d.db.Begin()
	return nil
}

func (d *Repository) GetFeedURLs(ctx context.Context) ([]string, error) {
	op := "repository.GetFeedURLs"
	getFeedQuery := `SELECT feed.feed_link
					 FROM feed
	`

	rows, err := d.db.QueryContext(ctx, getFeedQuery)
	if err != nil {
		slog.Error("error executing getFeedURLs query", "function", op, "error", err.Error())
		return nil, err
	}

	var feedURLs []string

	for rows.Next() {
		var feedURL string

		err = rows.Scan(&feedURL)
		if err != nil {
			slog.Error("error scanning into string", "error", err.Error())
			return nil, err
		}

		feedURLs = append(feedURLs, feedURL)
	}

	return feedURLs, nil
}
