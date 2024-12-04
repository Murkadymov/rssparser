package interfaces

import "context"

type Repository interface {
	InsertFeedURLs(ctx context.Context, feedURLs []string) error
	GetFeedURLs(ctx context.Context) ([]string, error)
}
