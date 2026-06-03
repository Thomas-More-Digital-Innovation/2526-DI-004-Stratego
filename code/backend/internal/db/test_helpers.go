package db

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Test constants for use across the test suite to ensure consistency and satisfy goconst
const (
	TestPassword       = "StrongPassword1"
	TestPasswordAlt    = "NewStrongPassword2"
	TestUser           = "testuser"
	TestUserLogin      = "loginuser"
	TestUserHistory    = "historyuser"
	TestUserPassChange = "passuser"
)

var nameSanitizer = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// SetupTestDB initializes an in-memory SQLite database for testing and sets the global DB instance.
func SetupTestDB(t *testing.T) *gorm.DB {
	// Sanitize test name to avoid URI issues with slashes/special characters
	safeName := nameSanitizer.ReplaceAllString(t.Name(), "_")

	// Use a unique name per test to ensure isolation while allowing shared cache for concurrency
	dbName := fmt.Sprintf("file:%s?mode=memory&cache=shared", safeName)
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Set the global DB instance for handlers to use
	DB = db

	// AutoMigrate all models
	err = db.AutoMigrate(AllModels...)
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	return db
}

// SetupDBTest swaps the global DB pointer and restores it on test cleanup
func SetupDBTest(t *testing.T) *gorm.DB {
	testDB := SetupTestDB(t)
	oldDB := DB
	DB = testDB
	t.Cleanup(func() {
		DB = oldDB
	})
	return testDB
}
