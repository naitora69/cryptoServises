package repository

import (
	"database/sql"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

func RunMigrations(dsn string) error {

	_, filename, _, _ := runtime.Caller(0)

	// internal/repository → в корень → migrations
	path := filepath.Join(filepath.Dir(filename), "../../migrations")

	// делаем URL-совместимый путь
	path = strings.ReplaceAll(path, "\\", "/")

	m, err := migrate.New(
		"file://"+path,
		dsn,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// PostgresDSN (Data Source Name) формирует строку подключения PostgreSQL
func PostgresDSN(config PostgresConfig) string {
	return "postgres://" +
		config.Username + ":" + config.Password + "@" +
		config.Host + ":" + config.Port + "/" +
		config.Name + "?" + "sslmode=" + config.SSLMode
}

func NewPostgresDB(config PostgresConfig) (*sql.DB, error) {
	dsn := PostgresDSN(config)
	if err := RunMigrations(dsn); err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
