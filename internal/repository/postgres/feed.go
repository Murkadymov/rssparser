package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed"
	"net/url"
	"strings"
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

func (d *Repository) GetLinkPrimaryID(ctx context.Context, link string) (int, error) {

	var feedID int

	urlParsed, err := url.Parse(strings.TrimSpace(link))
	if err != nil {
		fmt.Println(link, urlParsed)
		return 0, err
	}

	err = d.db.QueryRow(`
	SELECT id
	FROM feed
	WHERE feed_link ILIKE $1;
	`, "%"+urlParsed.Host+"%").Scan(&feedID)

	fmt.Println("FEEED ID AND LINK: ", feedID, link)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil // Возвращаем 0, если запись не найдена
		}
		return 0, err
	}

	return feedID, nil
}

func (d *Repository) InsertFeedURLs(ctx context.Context, item *gofeed.Feed) error {
	id, err := d.GetLinkPrimaryID(ctx, item.Link)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(`
	INSERT INTO feed_content(feed_id, title, description)
	VALUES ($1,$2,$3)
	ON CONFLICT DO NOTHING
	`, id, item.Title, item.Description,
	)
	if err != nil {
		return err
	}

	return nil
}
