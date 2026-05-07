// Package db provides database access and helpers
package db

func isPostgresDialect() bool {
	return DB.Name() == "postgres"
}
