package utils

import (
	"database/sql"
	"os"
)

func OpenDB() (*sql.DB, error) {
	return sql.Open("mysql", os.Getenv("DATABASE_URL"))
}
