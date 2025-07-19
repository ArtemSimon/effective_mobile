package repository

import (
	"effective_mobile/internal/config"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnectPostgresDB(config *config.Config_PG) (*gorm.DB, error) {
	connection_db := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
		config.DBSSLMode,
	)

	// Открываем подключение к базе через ORM
	db, err := gorm.Open(postgres.Open(connection_db), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	return db, nil
}
