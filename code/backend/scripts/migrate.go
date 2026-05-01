package main

import (
	"context"
	"digital-innovation/stratego/db"
	"fmt"
	"os"
)

func main() {
	// Initialize DB connection using env vars
	if err := db.InitDB(); err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.CloseDB()

	fmt.Println("Checking for pending migrations...")
	if err := db.RunMigrations(context.Background()); err != nil {
		fmt.Printf("Migration error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database is up to date!")
}
