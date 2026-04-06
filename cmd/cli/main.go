package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnectDB(ctx context.Context, connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil , fmt.Errorf("unable to ping database: %w", err)
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

func EncodeBase62(num int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if num == 0 {
		return string(charset[0])
	}

	var result strings.Builder
	for num > 0 {
		result.WriteByte(charset[num % 62])
		num = num / 62
	}
	return reverse(result.String())
}

func reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes) - 1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func UpdateRecord(ctx context.Context, conn *pgx.Conn, id int, shortKey string) error {
	_, err := conn.Exec(ctx, `UPDATE urls SET short_key = $1 WHERE id = $2`, shortKey, id)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %v <url>\n", os.Args[0])
		os.Exit(1)
	}

	url := os.Args[1]

	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to load .env file: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	conn, err := ConnectDB(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	id, existingShortKey, err := GetOrInsertURL(ctx, conn, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var finalShortKey string
	if existingShortKey != nil {
		finalShortKey = *existingShortKey
	} else {
		finalShortKey = EncodeBase62(id)
		if err := UpdateRecord(ctx, conn, id, finalShortKey); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("URL: https://localhost:8080/%v\n", finalShortKey)
}
