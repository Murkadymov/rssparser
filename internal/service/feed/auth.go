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
	repo      postgres.HTTPRepository
	svcLogger *slog.Logger
}

func NewAuthService(repo postgres.HTTPRepository, svcLogger *slog.Logger) *AuthService {
	return &AuthService{
		repo:      repo,
		svcLogger: svcLogger,
	}
}

func (a *AuthService) AddUser(user *api.User) (*int, error) {
	const op = "service.feed.AddUser"

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("%s: hashPassword: %w", op, err)
	}

	user.CreatedAt = time.Now()

	userID, err := a.repo.AddUser(user.Username, hashedPassword, user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: repo.AddUser: %w", op, err)
	}

	return userID, nil
}

func (a *AuthService) Login(user *api.User) error {
	const op = "service.feed.Login"

	existingPassword, err := a.repo.ValidateUser(user.Username)
	if err != nil {
		return fmt.Errorf("%s: repo.ValidateUser: %w", op, err)
	}

	if err = validatePassword([]byte(existingPassword), user.Password); err != nil {
		return fmt.Errorf("%s: validatePassword: %w", op, err)
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("GenerateFromPassword: %w", err)
	}
	log.Debug("successful hash generation", "password", password, "hash", string(hash))

	return string(hash), nil
}

func validatePassword(hashedPassword []byte, password string) error {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return fmt.Errorf("CompareHashAndPassword: %w", err)
	}

	log.Debug("success comparing hash and password")

	return nil
}
