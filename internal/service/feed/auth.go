package feed

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"rssparser/internal/models/api"
	"rssparser/internal/repository/interfaces"
	"time"
)

type AuthService struct {
	repo      interfaces.HTTPRepository
	svcLogger *slog.Logger
}

func NewAuthService(repo interfaces.HTTPRepository, svcLogger *slog.Logger) *AuthService {
	return &AuthService{
		repo:      repo,
		svcLogger: svcLogger,
	}
}

func (a *AuthService) AddUser(user *api.User) (*int, error) {

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("hash user password: %w", err)
	}

	user.CreatedAt = time.Now()

	userID, err := a.repo.AddUser(user.Username, hashedPassword, user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("authService.AddUser: %w", err)
	}

	return userID, nil
}

func (a *AuthService) Login(user *api.User) error {

	existingPassword, err := a.repo.ValidateUser(user.Username)
	if err != nil {
		return fmt.Errorf("authService: %w", err)
	}

	if err = validatePassword([]byte(existingPassword), user.Password); err != nil {
		return fmt.Errorf("authService.Login: %w", err)
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate hash from passowrd: %w", err)
	}
	log.Debug("successful hash generation", "password", password, "hash", string(hash))

	return string(hash), nil
}

func validatePassword(hashedPassword []byte, password string) error {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return fmt.Errorf("validating password %w", err)
	}

	log.Debug("comparing hash and password successful")

	return nil
}
