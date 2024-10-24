package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"rssparser/internal/config"
)

func ConnectToDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
			cfg.User,
			cfg.Password,
			cfg.DB,
			cfg.Host,
			cfg.Port,
			cfg.SSLMode,
		),
	)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	slog.Debug("successful connection to DB")

	return db, nil
}
