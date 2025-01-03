package handlers

import (
	"context"
	"rssparser/internal/models/api"
)

type FeedService interface {
	InsertFeedSource(ctx context.Context, feedSource *api.FeedSource) error
}

type AuthService interface {
	AddUser(user *api.User) (*int, error)
	Login(secret string, user *api.User) (string, error)
}
