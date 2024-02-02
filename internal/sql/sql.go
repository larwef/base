package sql

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/larwef/base/internal/sql/gen"
	"github.com/larwef/base/internal/sql/schema"
	migrate "github.com/rubenv/sql-migrate"
)

type Config struct {
	DBConnectionString string        `envconfig:"DB_CONNECTION_STRING" required:"true"`
	MaxDBConnections   int32         `envconfig:"DB_MAX_CONNECTIONS" default:"10"`
	PingInterval       time.Duration `envconfig:"DB_PING_INTERVAL" default:"3s"`
	PingTimeout        time.Duration `envconfig:"DB_PING_TIMEOUT" default:"15s"`
}

type Postgres struct {
	db     *pgxpool.Pool
	q      *gen.Queries
	logger *slog.Logger
}

type PostgresOption func(*Postgres)

func NewPostgres(ctx context.Context, conf Config, options ...PostgresOption) (*Postgres, error) {
	dbConf, err := pgxpool.ParseConfig(conf.DBConnectionString)
	if err != nil {
		return nil, err
	}
	dbConf.MaxConns = conf.MaxDBConnections
	db, err := pgxpool.NewWithConfig(ctx, dbConf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	pg := &Postgres{
		db:     db,
		q:      gen.New(db),
		logger: slog.Default(),
	}
	for _, option := range options {
		option(pg)
	}
	if err := pg.pingRetry(ctx, conf.PingInterval, conf.PingTimeout); err != nil {
		return nil, err
	}
	if err := pg.migrateUp(); err != nil {
		return nil, fmt.Errorf("failed to migrate up: %w", err)
	}
	return pg, nil
}

func WithLogger(logger *slog.Logger) PostgresOption {
	return func(p *Postgres) {
		p.logger = logger
	}
}

func (p *Postgres) Close() {
	p.db.Close()
}

func (p *Postgres) pingRetry(ctx context.Context, pingInterval, timeout time.Duration) error {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()
	timeoutExceeded := time.After(timeout)
	for {
		if err := p.db.Ping(ctx); err == nil {
			return nil
		} else {
			p.logger.Info(fmt.Sprintf("connecting to db failed: %v. Retrying in %s", err, pingInterval))
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeoutExceeded:
			return fmt.Errorf("db connection timed out after %s", timeout)
		case <-ticker.C:
			continue
		}
	}
}

func (pg *Postgres) migrateUp() error {
	conf := pg.db.Config()
	sqlDb := stdlib.OpenDB(*conf.ConnConfig)
	defer sqlDb.Close()
	n, err := migrate.Exec(sqlDb, "postgres", migrate.HttpFileSystemMigrationSource{
		FileSystem: http.FS(schema.Migrations),
	}, migrate.Up)
	if err != nil {
		return err
	}
	pg.logger.Info(fmt.Sprintf("database migration: %d migrations applied up\n", n))
	return nil
}
