package database

import (
	"context"
	"fmt"

	"github.com/dinizgab/booking-mvp/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Database interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row
	Close() error
}

type databaseImpl struct {
	conn *pgx.Conn
}

func New(config *config.DBConfig) (Database, error) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, config.DBUrl)
	if err != nil {
		return nil, fmt.Errorf("Database.New: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Database.New: %w", err)
	}

	return &databaseImpl{
		conn: conn,
	}, nil
}

func (d *databaseImpl) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return d.conn.Exec(ctx, sql, arguments...)
}

func (d *databaseImpl) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	return d.conn.Query(ctx, sql, arguments...)
}

func (d *databaseImpl) QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row {
	return d.conn.QueryRow(ctx, sql, arguments...)
}

func (d *databaseImpl) Close() error {
	return d.conn.Close(context.Background())
}
