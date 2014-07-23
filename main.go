package main

import (
	"database/sql"
	"flag"
	// "fmt"
	"github.com/code-mobi/tvthailand-api/admin"
	"github.com/code-mobi/tvthailand-api/api2"
	"github.com/dropbox/godropbox/memcache"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"html/template"
	"log"
	"net"
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

	conn, err := net.Dial("tcp", os.Getenv("MEMCACHED_HOST"))
	if err != nil {
		log.Fatal(err)
	}
	client := memcache.NewRawClient(0, conn)

	m := martini.Classic()
	m.Map(db)
	m.Map(client)
	m.Use(render.Renderer(render.Options{
  	Layout: "layout",
	}))

	m.Get("/", func() string {
		return "Hello world!"
	})

	m.Group("/api2", func(r martini.Router) {
		r.Get("/advertise", api2.AdvertiseListHandler)
		r.Get("/section", api2.SectionListHandler)
		r.Get("/category", api2.CategoryListHandler)
		r.Get("/category/:id", api2.CategoryShowHandler)
		r.Get("/category/:id/:start", api2.CategoryShowHandler)
		r.Get("/channel", api2.ChannelListHandler)
		r.Get("/channel/:id", api2.ChannelShowHandler)
		r.Get("/channel/:id/:start", api2.ChannelShowHandler)
		r.Get("/radio", api2.RadioListHandler)
		r.Get("/episode/:id", api2.EpisodeListHandler)
		r.Get("/episode/:id/:start", api2.EpisodeListHandler)
	})

	m.Group("/admin", func(r martini.Router) {
		r.Get("/encrypt", admin.EncryptHandler)
		r.Post("/encrypt", admin.EncryptUpdateHandler)
	})

	if err := http.ListenAndServe(":"+*port, m); err != nil {
		panic(err)
	}
}

// func main() {
// 	port := flag.String("port", "9000", "PORT")
// 	flag.Parse()
//
// 	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
//
// 	conn, _ := net.Dial("tcp", os.Getenv("MEMCACHED_HOST"))
// 	client := memcache.NewRawClient(0, conn)
//
// 	http.Handle("/static/", http.FileServer(http.Dir("./")))
// 	http.Handle("/api2/", &api2.Api2Handler{Db: db, MemcacheClient: client})
// 	http.Handle("/admin/", &admin.AdminHandler{Db: db})
// 	http.HandleFunc("/flush", func(w http.ResponseWriter, r *http.Request) {
// 		client.Flush(1)
// 		fmt.Fprintf(w, "Flush")
// 	})
// 	http.HandleFunc("/", HomeHandler)
// 	if err := http.ListenAndServe(":"+*port, nil); err != nil {
// 		panic(err)
// 	}
// }

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
