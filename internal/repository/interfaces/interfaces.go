package interfaces

import (
	"context"
	"github.com/mmcdole/gofeed"
)

type Repository interface {
	InsertFeedURLs(ctx context.Context, item *gofeed.Feed) error
	GetFeedURLs(ctx context.Context) ([]string, error)
}
type Cache interface {
	Get() []string
	Set(items []string)
	Update(items []string)
}
