package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"shortener/internal/db"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func handleEncodedSite(w http.ResponseWriter, req *http.Request) {
	shortKey := req.URL.Path[1:]
	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to load .env file: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	conn, err := db.ConnectDB(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	longURL, err := db.GetURLFromShortKey(ctx, conn, shortKey)
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
	http.HandleFunc("/", handleEncodedSite)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintln(os.Stderr, "error with server", err)
		os.Exit(1)
	}
}
