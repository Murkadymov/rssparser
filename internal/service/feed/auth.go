package feed

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"rssparser/internal/models/api"
	"time"
)

type JWTGenerationError struct {
	msg string
}

func (j *JWTGenerationError) Error() string {
	return fmt.Sprintf("JWT token generation error: %s", j.msg)
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthService struct {
	repo      HTTPRepository
	svcLogger *slog.Logger
}

func NewAuthService(repo HTTPRepository, svcLogger *slog.Logger) *AuthService {
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

func (a *AuthService) Login(secret string, user *api.User) (string, error) {
	// Проверка пользователя в базе
	existingPassword, err := a.repo.ValidateUser(user.Username)
	if err != nil {
		return "", fmt.Errorf("repo.ValidateUser: %w", err)
	}

	// Проверка пароля
	if err = validatePassword([]byte(existingPassword), user.Password); err != nil {
		return "", fmt.Errorf("validatePassword: %w", err)
	}

	// Генерация токена
	token, err := GenerateToken(user.Username, secret)
	if err != nil {
		return "", fmt.Errorf("GenerateToken: %w", &JWTGenerationError{
			msg: err.Error(),
		})
	}

	return token, nil
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

func GenerateToken(username, secret string) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Срок действия токена: 24 часа
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
