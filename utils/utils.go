package utils

import (
	"database/sql"
	"os"
	"time"
)

func OpenDB() (*sql.DB, error) {
	return sql.Open("mysql", os.Getenv("DATABASE_URL"))
}

func GetTimeStamp() string {
	t := time.Now().Local()
	return t.Format("20060102150405")
}
