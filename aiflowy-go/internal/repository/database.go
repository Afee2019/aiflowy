package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/aiflowy/aiflowy-go/internal/config"
	"github.com/aiflowy/aiflowy-go/pkg/logger"
	"go.uber.org/zap"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB(cfg *config.DatabaseConfig) error {
	var err error

	dsn := cfg.DSN()
	db, err = sql.Open(cfg.Driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connected successfully",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Database),
	)

	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// TestQuery executes a test query to verify database connection
func TestQuery() (map[string]interface{}, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM tb_sys_account").Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to query tb_sys_account: %w", err)
	}

	// Get first account for verification
	var id, loginName string
	err = db.QueryRow("SELECT id, login_name FROM tb_sys_account LIMIT 1").Scan(&id, &loginName)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to query account: %w", err)
	}

	result := map[string]interface{}{
		"table":           "tb_sys_account",
		"total_count":     count,
		"first_user_id":   id,
		"first_login_name": loginName,
	}

	return result, nil
}
