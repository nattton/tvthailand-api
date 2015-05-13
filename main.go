package main

import (
	"flag"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"

	"github.com/code-mobi/tvthailand-api/admin"
	"github.com/code-mobi/tvthailand-api/api2"
	"github.com/code-mobi/tvthailand-api/bot"
	"github.com/code-mobi/tvthailand-api/utils"
	"github.com/dropbox/godropbox/memcache"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/render"
)

func main() {
	port := flag.String("port", "9000", "PORT")
	command := flag.String("command", "", "COMMAND")
	user := flag.String("user", "", "USER")
	channel := flag.String("channel", "", "CHANNEL")
	keyword := flag.String("keyword", "", "KEYWORD")
	flag.Parse()

	db, err := utils.OpenDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if *command == "botrun" {
		fmt.Println(*command)
		b := bot.NewBot(db)
		if *user == "" {
			b.CheckRobotChannel()
		} else {
			b.CheckYoutubeUser(*user, 1, 30)
		}

	} else if *command == "findchannel" {
		fmt.Println(*command)
		b := bot.NewBot(db)
		if *user == "" {
			b.FindChannel()
		}

	} else if *command == "findvideochannel" {
		fmt.Println(*command)
		b := bot.NewBot(db)
		if *user == "" || *channel == "" {
			fmt.Println("Must have -user=... -channel=...")
		} else {
			b.CheckVideoInChannel(*user, *channel)
		}

	} else if *command == "botrunvideo" {
		fmt.Println(*command)
		b := bot.NewBot(db)
		if *user == "" {
			fmt.Println("Must have -user=...")
		} else {
			b.CheckAllVideoInYoutubeUserAndKeyword(*user, *keyword)
		}

	} else {
		conn, err := net.Dial("tcp", os.Getenv("MEMCACHED_HOST"))
		if err != nil {
			fmt.Println(err)
		}
		client := memcache.NewRawClient(0, conn)

		m := martini.Classic()
		m.Map(db)
		m.Map(client)
		m.Use(render.Renderer(render.Options{
			Directory:  "templates",
			Layout:     "layout",
			Delims:     render.Delims{"{[{", "}]}"},
			Charset:    "UTF-8",
			IndentJSON: true,
		}))

		m.Get("/", func(r render.Render) {
			r.JSON(200, map[string]interface{}{"hello": "world สวัสดี"})
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
			r.Get("/otv", admin.OtvHandler)
			r.Post("/otv", admin.OtvProcessHandler)
			r.Get("/botvideo", admin.BotVideoHandler)
			r.Post("/botvideo", admin.BotVideoPostHandler)
			r.Get("/botvideo.json", admin.BotVideoJSONHandler)
			r.Get("/show.json", admin.ShowJSONHandler)
			r.Get("/youtube", admin.YoutubeHandler)
			r.Get("/youtube.search.channel", admin.YoutubeSearchChannelJSONHandler)
			// r.Get("/youtube.json", admin.YoutubeJSONHandler)
			r.Get("/showlist", admin.ShowListHandler)

		})

		m.Get("/flush", func() string {
			client.Flush(1)
			return "Flush!!!"
		})

		if err := http.ListenAndServe(":"+*port, m); err != nil {
			panic(err)
		}
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
