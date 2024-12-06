package storage

import (
	"context"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"

	"github.com/tonytcb/inventory-management-system/internal/app/config"
)

func NewPostgresConnection(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "error to parse database url")
	}

	pgxConfig.AfterConnect = func(_ context.Context, c *pgx.Conn) error {
		pgxdecimal.Register(c.TypeMap())
		pgxUUID.Register(c.TypeMap())
		return nil
	}

	// optional configs
	if cfg.DatastoreMaxOpenConn != nil {
		pgxConfig.MaxConns = *cfg.DatastoreMaxOpenConn
	}
	if cfg.DatastoreMinOpenConn != nil {
		pgxConfig.MinConns = *cfg.DatastoreMinOpenConn
	}

	conn, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error to create database connection")
	}

	return conn, nil
}
