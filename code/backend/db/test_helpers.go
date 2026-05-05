package db

import (
	"digital-innovation/stratego/models"
	"fmt"
	"regexp"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var nameSanitizer = regexp.MustCompile(`[^a-zA-Z0-9_]`)

func setupSQLiteDB(t *testing.T) *gorm.DB {
	// Sanitize test name to avoid URI issues with slashes/special characters
	safeName := nameSanitizer.ReplaceAllString(t.Name(), "_")

	// Use a unique name per test to ensure isolation while allowing shared cache for concurrency
	dbName := fmt.Sprintf("file:%s?mode=memory&cache=shared", safeName)
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// AutoMigrate all models
	err = db.AutoMigrate(
		&models.User{},
		&models.UserStats{},
		&models.BoardSetup{},
		&models.RefreshToken{},
		&models.Game{},
		&models.GameMove{},
	)
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}
