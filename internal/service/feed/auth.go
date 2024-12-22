package feed

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"rssparser/internal/models/api"
	"rssparser/internal/repository/postgres"
	"time"
)

type AuthService struct {
	svcLogger *slog.Logger
	repo      postgres.Repository
}

func (a *AuthService) Register(user *api.User) error {

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("hash user password: %w", err)
	}

	user.CreatedAt = time.Now()

	a.repo.adduser(user.Username, string(hashedPassword), user.CreatedAt)

}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate hash from passowrd: %w", err)
	}
	log.Debug("successful hash generation", "password", password, "hash", string(hash))

	return string(hash), nil
}
