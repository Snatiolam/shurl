package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"shortener/internal/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type application struct {
	pool *pgxpool.Pool
}

func (app *application) handleEncodedSite(w http.ResponseWriter, req *http.Request) {
	shortKey := req.URL.Path[1:]
	ctx := context.Background()

	longURL, err := db.GetURLFromShortKey(ctx, app.pool, shortKey)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, longURL, http.StatusPermanentRedirect)
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to load .env file: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to database pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	app := application{
		pool: pool,
	}
	http.HandleFunc("/", app.handleEncodedSite)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintln(os.Stderr, "error with server", err)
		os.Exit(1)
	}
}
