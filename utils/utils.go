package utils

import (
	"database/sql"
	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"os"
	"time"
)

func OpenDB() (*sql.DB, error) {
	return sql.Open("mysql", os.Getenv("DATABASE_DSN"))
}

func OpenGormDB() (gorm.DB, error) {
	db, err := gorm.Open("mysql", os.Getenv("DATABASE_DSN"))
	db.LogMode(martini.Env == "development")
	return db, err
}

func GetTimeStamp() string {
	t := time.Now().Local()
	return t.Format("20060102150405")
}
