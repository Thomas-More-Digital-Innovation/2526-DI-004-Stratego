// Package main provides a script to run database migrations
package main

import (
	"context"
	"digital-innovation/stratego/db"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Initialize DB connection using env vars
	if err := db.InitDB(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() {
		err := db.CloseDB()
		if err != nil {
			fmt.Printf("Error closing database: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Println("Checking for pending migrations...")
	if err := db.RunMigrations(context.Background()); err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	fmt.Println("Database is up to date!")
	return nil
}
