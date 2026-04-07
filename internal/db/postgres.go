package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func ConnectDB(ctx context.Context, connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return conn, nil
}

func GetOrInsertURL(ctx context.Context, conn *pgx.Conn, url string) (int, *string, error) {
	var id int
	var shortKey *string

	query := `
		INSERT INTO urls (long_url)
		VALUES ($1)
		ON CONFLICT (long_url) DO UPDATE
		SET long_url = EXCLUDED.long_url
		RETURNING id, short_key
	`
	err := conn.QueryRow(ctx, query, url).Scan(&id, &shortKey)
	if err != nil {
		return -1, nil, err
	}
	return id, shortKey, nil
}

func UpdateRecord(ctx context.Context, conn *pgx.Conn, id int, shortKey string) error {
	_, err := conn.Exec(ctx, `UPDATE urls SET short_key = $1 WHERE id = $2`, shortKey, id)
	if err != nil {
		return err
	}
	return nil
}
