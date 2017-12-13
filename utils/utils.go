package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
)

func OpenDB() (*sql.DB, error) {
	return sql.Open("mysql", os.Getenv("DATABASE_DSN"))
}

func OpenGormDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", os.Getenv("DATABASE_DSN"))
	db.LogMode(martini.Env == "development")
	return db, err
}

func GetTimeStamp() string {
	t := time.Now().Local()
	return t.Format("20060102150405")
}

func JSONP(writer http.ResponseWriter, code int, callback string, v interface{}) {
	writer.WriteHeader(code)
	jsonBytes, _ := json.Marshal(v)
	if callback != "" {
		fmt.Fprintf(writer, "%s(%s)", callback, jsonBytes)
	} else {
		writer.Write(jsonBytes)
	}
}
