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

	fmt.Println("PARARARSARASRS", "%"+parsedURL+"%")

	const getLinkPrimaryIDQuery = `SELECT id
								   FROM feed
	                               WHERE feed_link ILIKE $1;`
	err := d.db.QueryRow(getLinkPrimaryIDQuery, "%"+parsedURL+"%").Scan(&feedPrimaryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, sql.ErrNoRows // Возвращаем 0, если запись не найдена
		}
		return 0, err
	}
	fmt.Println("FEEED ID AND LINK: ", feedPrimaryID, parsedURL)
	return feedPrimaryID, nil
}

func (d *Repository) InsertFeedContent(
	ctx context.Context,
	feedPrimaryID int,
	feedTitle string,
	feedDescription string,
	feedPubDate *time.Time,
	feedLink string) error {
	const insertContentQuery = `INSERT INTO feed_content(feed_id, title, description, published_at, pub_link)
								VALUES ($1,$2,$3,$4,$5)
								ON CONFLICT DO NOTHING`

	_, err := d.db.Exec(
		insertContentQuery,
		feedPrimaryID,
		feedTitle,
		feedDescription,
		feedPubDate,
		feedLink,
	)

	if err != nil {
		return fmt.Errorf("insert feed content info: %w", err)
	}

	return nil
}

func (d *Repository) GetExistingPubDate(feedLink string) (string, error) {
	const GetPubDateQuery = `SELECT published_at
							 FROM feed_content
							 WHERE pub_link = $1`

	var existingPubDate string

	err := d.db.QueryRow(GetPubDateQuery, feedLink).Scan(&existingPubDate)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return "", nil
		}
		return "", fmt.Errorf("existing pub date info: %w", err)

	}
	fmt.Println("RECEIVED EXISTITNG PUB DATE", existingPubDate)

	return existingPubDate, nil
}

func (d *Repository) InsertFeedSource(ctx context.Context, feedLink string) error {
	const InsertFeedSourceQuery = `INSERT INTO feed(feed_link)
    							   VALUES ($1)
								   ON CONFLICT DO NOTHING;
    							   `

	if _, err := d.db.ExecContext(ctx, InsertFeedSourceQuery, feedLink); err != nil {
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("insert feed source info: %w", err)
		}
		//TODO:контекст таймаут

		return fmt.Errorf("insert feed source info: %w", err)
	}

	return nil
}
