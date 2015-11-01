package main

import (
	"flag"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"

	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/dropbox/godropbox/memcache"
	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/go-martini/martini"
	_ "github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/martini-contrib/auth"
	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/martini-contrib/render"
	"github.com/code-mobi/tvthailand-api/admin"
	"github.com/code-mobi/tvthailand-api/api2"
	"github.com/code-mobi/tvthailand-api/bot"
	"github.com/code-mobi/tvthailand-api/data"
	"github.com/code-mobi/tvthailand-api/utils"
)

type CmdParam struct {
	Command string
	User    string
	Channel string
	Q       string
	Start   int
	Stop    int
}

func main() {
	command := flag.String("command", "", "COMMAND = botrun | findchannel | findvideochannel(user, channel)")
	user := flag.String("user", "", "USER")
	channel := flag.String("channel", "", "CHANNEL")
	q := flag.String("q", "", "QUERY")
	start := flag.Int("start", 0, "START")
	stop := flag.Int("stop", 0, "STOP")
	flag.Parse()

	if *command != "" {
		cmdParam := &CmdParam{
			*command,
			*user,
			*channel,
			*q,
			*start,
			*stop,
		}
		processCommand(cmdParam)
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	db, err := utils.OpenDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	dbg, err := utils.OpenGormDB()
	if err != nil {
		panic(err.Error())
	}
	defer dbg.Close()

	conn, err := net.Dial("tcp", os.Getenv("MEMCACHED_SERVER"))
	if err != nil {
		panic(err.Error())
	}
	client := memcache.NewRawClient(0, conn)

	m := martini.Classic()
	m.Map(db)
	m.Map(dbg)
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

	authAdmin := auth.BasicFunc(func(username, password string) bool {
		return auth.SecureCompare(username, "saly") && auth.SecureCompare(password, "admin888")
	})

	m.Group("/admin", func(r martini.Router) {
		r.Get("/encrypt", admin.EncryptHandler)
		r.Post("/encrypt", admin.EncryptUpdateHandler)
		r.Get("/otv", authAdmin, admin.OtvHandler)
		r.Post("/otv", authAdmin, admin.OtvProcessHandler)
		r.Get("/botvideo", authAdmin, admin.BotVideoHandler)
		r.Post("/botvideo", authAdmin, admin.BotVideoPostHandler)
		r.Get("/botvideo.json", admin.BotVideoJSONHandler)
		r.Get("/show.json", admin.ShowJSONHandler)
		r.Get("/youtube", admin.YoutubeHandler)
		r.Get("/youtube.search.channel", admin.YoutubeSearchChannelJSONHandler)
		r.Get("/youtube.playlistItems", admin.YoutubePlaylistItemJSONHandler)
		r.Get("/showlist", admin.ShowListHandler)
		r.Get("/krobkruakao", admin.KrobkruakaoHandler)
		r.Get("/krobkruakao.json", admin.KrobkruakaoJSONHandler)

	})

	m.Get("/flush", func() string {
		client.Flush(1)
		return "Flush!!!"
	})

	m.RunOnAddr(":" + port)
}

func processCommand(cmd *CmdParam) {
	db, err := utils.OpenDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	dbg, err := utils.OpenGormDB()
	if err != nil {
		panic(err.Error())
	}
	defer dbg.Close()

	fmt.Println(cmd.Command)
	switch cmd.Command {
	case "krobkruakao":
		admin.ExampleKrobkruakao()
	case "botrun":
		b := bot.NewBot(db)
		b.CheckRobotChannel()
	case "botkrobkruakao":
		b := bot.NewBot(db)
		b.CheckKrobkruakao(cmd.Start)
	case "findchannel":
		b := bot.NewBot(db)
		b.FindChannel()
	case "findvideochannel":
		b := bot.NewBot(db)
		if cmd.User == "" || cmd.Channel == "" {
			fmt.Println("Must have -user=... -channel=...")
		} else {
			b.CheckVideoInChannel(cmd.User, cmd.Channel, cmd.Q)
		}
	case "otvupdate":
		otv := &admin.Otv{Db: db}
		otv.UpdateModified()
	case "findconan":
		b := bot.NewBot(db)
		for i := cmd.Start; i < cmd.Stop; i++ {
			ep := fmt.Sprintf("EP%%20%d", i)
			b.CheckVideoInChannel("conanofficial", "UCmbpqlWIyoPEVUzU6iTf1OA", ep)
		}
	case "runbotpl":
		data.RunBotPlaylists(&dbg)
	case "migrate_botvideo":
		data.MigrateUsernameToChannelID(&dbg)
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
