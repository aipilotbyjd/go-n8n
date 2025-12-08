package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database configuration
type Config struct {
	Driver                string        `mapstructure:"driver"`
	Host                  string        `mapstructure:"host"`
	Port                  int           `mapstructure:"port"`
	User                  string        `mapstructure:"user"`
	Password              string        `mapstructure:"password"`
	Name                  string        `mapstructure:"name"`
	SSLMode               string        `mapstructure:"ssl_mode"`
	MaxConnections        int           `mapstructure:"max_connections"`
	MaxIdleConnections    int           `mapstructure:"max_idle_connections"`
	ConnectionMaxLifetime time.Duration `mapstructure:"connection_max_lifetime"`
	LogLevel              string        `mapstructure:"log_level"`
}

// DB wraps the database connection
type DB struct {
	*gorm.DB
}

// Connect establishes a database connection
func Connect(cfg Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	// Set log level
	logLevel := logger.Silent
	switch cfg.LogLevel {
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logLevel),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	if cfg.MaxConnections > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxConnections)
	}
	if cfg.MaxIdleConnections > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	}
	if cfg.ConnectionMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Migrate runs database migrations
func (db *DB) Migrate(models ...interface{}) error {
	return db.AutoMigrate(models...)
}

// Transaction executes a function within a transaction
func (db *DB) Transaction(fn func(*gorm.DB) error) error {
	return db.DB.Transaction(fn)
}

// EnableUUID enables UUID extension in PostgreSQL
func (db *DB) EnableUUID() error {
	return db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error
}
