package postgres

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var schemaSQL string

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfigFromEnv() (DatabaseConfig, error) {
	var missing []string

	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	if host == "" {
		missing = append(missing, "DB_HOST")
	}

	port := strings.TrimSpace(os.Getenv("DB_PORT"))
	if port == "" {
		missing = append(missing, "DB_PORT")
	}

	user := strings.TrimSpace(os.Getenv("DB_USER"))
	if user == "" {
		missing = append(missing, "DB_USER")
	}

	password := strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	if password == "" {
		missing = append(missing, "DB_PASSWORD")
	}

	dbname := strings.TrimSpace(os.Getenv("DB_NAME"))
	if dbname == "" {
		missing = append(missing, "DB_NAME")
	}

	if len(missing) > 0 {
		return DatabaseConfig{}, fmt.Errorf("missing required database environment variables: %s", strings.Join(missing, ", "))
	}

	sslmode := strings.TrimSpace(os.Getenv("DB_SSLMODE"))
	if sslmode == "" {
		sslmode = "disable"
	}

	return DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
		SSLMode:  sslmode,
	}, nil
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

func NewDatabase(dsn string) (*pgxpool.Pool, error) {
	if dsn == "" {
		cfg, err := LoadConfigFromEnv()
		if err != nil {
			return nil, err
		}
		dsn = cfg.DSN()
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres dsn config: %w", err)
	}

	config.MaxConns = 100
	config.MinConns = 10
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres with pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	if _, err := pool.Exec(ctx, schemaSQL); err != nil {
		pool.Close()
		return nil, fmt.Errorf("schema migration failed: %w", err)
	}

	log.Println("[INFO] PostgreSQL (pgxpool) connected and schema migrated successfully.")
	return pool, nil
}
