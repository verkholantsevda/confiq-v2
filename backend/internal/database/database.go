package database

import (
	"log/slog"
	"os"
	"path/filepath"

	"confiq/internal/config"
	"confiq/internal/configs"
	"confiq/internal/configtypes"
	"confiq/internal/endpoints"
	"confiq/internal/groups"
	"confiq/internal/logger"
	"confiq/internal/relations"
	"confiq/internal/users"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dir := filepath.Dir(cfg.DatabasePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	_, err := os.Stat(cfg.DatabasePath)
	dbExists := err == nil

	db, err := gorm.Open(sqlite.Open(cfg.DatabasePath), &gorm.Config{
		Logger: logger.Gorm(),
	})
	if err != nil {
		return nil, err
	}

	if !dbExists {
		if err := migrate(db); err != nil {
			return nil, err
		}

		slog.Info("database initialized")

	} else {
		slog.Info("existing database detected")

		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}

		if err := RunMigrations(sqlDB); err != nil {
			return nil, err
		}
		slog.Info("database migrations applied")
	}

	slog.Info(
		"sqlite connected",
		"path", cfg.DatabasePath,
	)

	return db, nil
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&users.User{},
		&groups.Group{},
		&endpoints.Endpoint{},
		&configtypes.ConfigType{},
		&configs.Config{},
		&relations.EndpointConfigType{},
		&relations.GroupEndpoint{},
	)
}
