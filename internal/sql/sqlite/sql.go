package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/larwef/base/internal/sql/sqlite/gen"
	"github.com/larwef/base/internal/sql/sqlite/schema"
	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

type Config struct {
	DBConnectionString string `envconfig:"DB_CONNECTION_STRING" required:"true"`
}

type SQLite struct {
	db *sql.DB
	q  *gen.Queries
}

func NewSqlite(ctx context.Context, conf Config) (*SQLite, error) {
	db, err := sql.Open("sqlite3", conf.DBConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	pg := &SQLite{
		db: db,
		q:  gen.New(db),
	}
	if err := pg.migrateUp(); err != nil {
		return nil, err
	}
	return pg, nil
}

func (s *SQLite) Close() {
	s.db.Close()
}

func (s *SQLite) migrateUp() error {
	n, err := migrate.Exec(s.db, "sqlite3", migrate.HttpFileSystemMigrationSource{
		FileSystem: http.FS(schema.Migrations),
	}, migrate.Up)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("database migration: %d migrations applied up\n", n))
	return nil
}

func (s *SQLite) Queries() *gen.Queries {
	return s.q
}
