package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type HTTPRepository struct {
	db *sql.DB
}

func NewHTTPRepository(db *sql.DB) *HTTPRepository {
	return &HTTPRepository{
		db: db,
	}
}

func (h *HTTPRepository) AddUser(name string, hashedPassword string, createdAt time.Time) (*int, error) {
	const AddUserQuery = `INSERT INTO users(
							  name, 
							  password,
							  createdAt
						  ) VALUES (
						      $1,
						      $2,
						      $3
						  ) RETURNING ID;`

	var userID *int

	err := h.db.QueryRow(AddUserQuery, name, hashedPassword, createdAt).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("add user info: %w", err)
	}

	return userID, nil
}

func (h *HTTPRepository) InsertFeedSource(ctx context.Context, feedLink string) error {
	const InsertFeedSourceQuery = `INSERT INTO feed(feed_link)
    							   VALUES ($1)
								   ON CONFLICT DO NOTHING;
    							   `

	if _, err := h.db.ExecContext(ctx, InsertFeedSourceQuery, feedLink); err != nil {
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("insert feed source info: %w", err)
		}
		//TODO:контекст таймаут

		return fmt.Errorf("insert feed source info: %w", err)
	}

	return nil
}
