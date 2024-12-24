package interfaces

import (
	"context"
	"time"
)

type Cache interface {
	Get() []string
	Set(items []string)
	Update(items []string)
}

type FeedRepository interface {
	InsertFeedContent(
		ctx context.Context,
		feedPrimaryID int,
		feedTitle string,
		feedDescription string,
		feedPubDate *time.Time,
		feedLink string,
	) error
	GetFeedURLs(ctx context.Context) ([]string, error)
	GetLinkPrimaryID(ctx context.Context, link string) (int, error)
	GetExistingPubDate(feedLink string) (string, error)
}

type HTTPRepository interface {
	InsertFeedSource(ctx context.Context, feedLink string) error
	AddUser(name string, hashedPassword string, createdAt time.Time) (*int, error)
	ValidateUser(username string) (string, error)
}
