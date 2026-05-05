package db

func isPostgresDialect() bool {
	return DB.Dialector.Name() == "postgres"
}
