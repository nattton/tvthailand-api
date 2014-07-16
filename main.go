package main

import (
	"database/sql"
	"flag"
	"github.com/code-mobi/tvthailand-api/api2"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.String("port", "9000", "PORT")
	flag.Parse()

	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.Handle("/static/", http.FileServer(http.Dir("./")))
	http.Handle("/api2/", &api2.Api2Handler{Db: db})
	http.HandleFunc("/", HomeHandler)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		panic(err)
	}
}

type Topic struct {
	TopicID int
	Name    string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var t = template.Must(template.ParseFiles(
		"templates/_base.html",
		"templates/index.html",
	))

	// t, _ := template.ParseFiles("templates/home.html")

	results := []Topic{
		Topic{1, "Title1"},
		Topic{2, "Title2"},
	}
	v := map[string]interface{}{
		"title":   "TV Thailand",
		"results": results,
	}

	if err := t.Execute(w, v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
