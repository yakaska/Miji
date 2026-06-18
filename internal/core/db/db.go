package db

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Close()
	OpTimeout() time.Duration
}

type ConnectionPool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewConnectionPool(ctx context.Context, config *Config) (*ConnectionPool, error) {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.User, config.Password),
		Host:   net.JoinHostPort(config.Host, strconv.Itoa(config.Port)),
		Path:   config.Database,
	}
	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()

	pgxConfig, err := pgxpool.ParseConfig(u.String())
	if err != nil {
		return nil, fmt.Errorf("parse pgx config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping connection pool: %w", err)
	}

	return &ConnectionPool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (c *ConnectionPool) OpTimeout() time.Duration {
	return c.opTimeout
}
