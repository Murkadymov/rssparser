package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (d *Repository) GetFeedURLs(ctx context.Context) ([]string, error) {
	getFeedQuery := `SELECT feed.feed_link
					 FROM feed`

	rows, err := d.db.QueryContext(ctx, getFeedQuery)
	if err != nil {
		return nil, err
	}

	var feedURLs []string

	for rows.Next() {
		var feedURL string

		err = rows.Scan(&feedURL)
		if err != nil {
			return nil, err
		}

		feedURLs = append(feedURLs, feedURL)
	}

	return feedURLs, nil
}

func (d *Repository) GetLinkPrimaryID(ctx context.Context, parsedURL string) (int, error) {
	var feedPrimaryID int

	const getLinkPrimaryIDQuery = `SELECT id
								   FROM feed
	                               WHERE feed_link ILIKE $1;`
	err := d.db.QueryRow(getLinkPrimaryIDQuery, "%"+parsedURL+"%").Scan(&feedPrimaryID)

	fmt.Println("FEEED ID AND LINK: ", feedPrimaryID, parsedURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil // Возвращаем 0, если запись не найдена
		}
		return 0, err
	}

	return feedPrimaryID, nil
}

func (d *Repository) InsertFeedContent(
	ctx context.Context,
	feedPrimaryID int,
	feedTitle string,
	feedDescription string,
	feedPubDate *time.Time) error {
	const insertContentQuery = `INSERT INTO feed_content(feed_id, title, description, published_at)
								VALUES ($1,$2,$3,$4)
								ON CONFLICT DO NOTHING`

	_, err := d.db.Exec(insertContentQuery, feedPrimaryID, feedTitle, feedDescription, feedPubDate)

	if err != nil {
		return err
	}

	return nil
}
