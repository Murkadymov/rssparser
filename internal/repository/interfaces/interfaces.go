package interfaces

import (
	"context"
	"time"
)

type Repository interface {
	InsertFeedContent(
		ctx context.Context,
		feedPrimaryID int,
		feedTitle string,
		feedDescription string,
		feedPubDate *time.Time) error
	GetFeedURLs(ctx context.Context) ([]string, error)
	GetLinkPrimaryID(ctx context.Context, link string) (int, error)
}
type Cache interface {
	Get() []string
	Set(items []string)
	Update(items []string)
}
