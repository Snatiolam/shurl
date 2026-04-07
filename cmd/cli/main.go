package main

import (
	"context"
	"fmt"
	"os"
	"shortener/internal/db"
	"shortener/internal/utils"

	"github.com/joho/godotenv"
)

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
	conn, err := db.ConnectDB(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	id, existingShortKey, err := db.GetOrInsertURL(ctx, conn, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var finalShortKey string
	if existingShortKey != nil {
		finalShortKey = *existingShortKey
	} else {
		finalShortKey = utils.EncodeBase62(id)
		if err := db.UpdateRecord(ctx, conn, id, finalShortKey); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("URL: http://localhost:8080/%v\n", finalShortKey)
}
