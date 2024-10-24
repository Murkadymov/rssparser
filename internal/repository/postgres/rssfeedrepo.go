package postgres

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (d *Repository) InsertFeedURLs(ctx context.Context, feedURLs []string) error {

	d.db.Begin()
	return nil
}
