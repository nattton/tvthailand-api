package admin

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/code-mobi/tvthailand-api/bot"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/render"
)

const MThaiPlayerURL = "http://video.mthai.com/cool/player/%s.html"

type EncryptResult struct {
	ShowID  int
	ListID  int
	VideoID string
}

func EncryptHandler(r render.Render) {
	var results []string
	emptymap := map[string]interface{}{
		"showid":  "",
		"message": "",
		"results": results,
	}
	r.HTML(200, "admin/encrypt", emptymap)
}

func EncryptUpdateHandler(db *sql.DB, params martini.Params, req *http.Request, r render.Render) {
	idType := req.FormValue("idType")
	id, _ := strconv.Atoi(req.FormValue("id"))

	var results []*EncryptResult

	if idType == "mthaiparseurl" {
		MthaiParseURL(db, r)
		return
	} else if idType == "" || id == 0 {
		emptymap := map[string]interface{}{
			"showid":  "",
			"message": "*Fill the form",
			"results": results,
		}
		r.HTML(200, "admin/encrypt", emptymap)
		return
	}

	var rows *sql.Rows

	switch idType {
	case "showid":
		rows, _ = db.Query("SELECT program_id, programlist_id, programlist_youtube FROM tv_programlist WHERE program_id = ? ORDER BY programlist_id DESC", id)
	case "listid":
		rows, _ = db.Query("SELECT program_id, programlist_id, programlist_youtube FROM tv_programlist WHERE programlist_id = ? ORDER BY programlist_id DESC", id)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			showID  int
			listID  int
			videoID string
		)
		if err := rows.Scan(&showID, &listID, &videoID); err != nil {
			log.Fatal(err)
		}
		updateVideoEncrypt(db, listID, videoID)
		results = append(results, &EncryptResult{showID, listID, videoID})
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	newmap := map[string]interface{}{
		"showid":  id,
		"message": fmt.Sprintf("Update Encrypt %s %d", idType, id),
		"results": results,
	}

	r.HTML(200, "admin/encrypt", newmap)
}

func MthaiParseURL(db *sql.DB, r render.Render) {
	var results []*EncryptResult
	rows, err := db.Query("SELECT program_id, programlist_id, programlist_youtube youtubeKey FROM tv_programlist WHERE programlist_banned = 0 AND programlist_src_type = 14 ORDER BY programlist_id DESC")
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			programID     int
			programlistID int
			youtubeKey    string
		)

		if err := rows.Scan(&programID, &programlistID, &youtubeKey); err != nil {
			log.Fatal(err)
		}

		fmt.Println(programID, programlistID, youtubeKey)
		videoIDs := strings.Split(youtubeKey, ",")

		var videoURLs []string
		for _, videoID := range videoIDs {
			videoURL := fmt.Sprintf(MThaiPlayerURL, videoID)
			videoURLs = append(videoURLs, videoURL)
		}
		videoResult := strings.Join(videoURLs, ",")
		fmt.Println("videoResult :", videoResult)
		encryptResult := EncryptVideo(videoResult)
		fmt.Println("encryptResult :", encryptResult)

		_, err := db.Exec("UPDATE tv_programlist SET programlist_youtube = ?, programlist_youtube_encrypt = ?, programlist_src_type = 11, mthai_video = ? WHERE programlist_id = ? ORDER BY program_id DESC", videoResult, encryptResult, youtubeKey, programlistID)
		if err != nil {
			panic(err)
		}
		results = append(results, &EncryptResult{programID, programlistID, youtubeKey})
	}

	emptymap := map[string]interface{}{
		"showid":  "",
		"message": "",
		"results": results,
	}
	r.HTML(200, "admin/encrypt", emptymap)
}

func OtvHandler(r render.Render) {
	var results []*OtvShowListItem
	newmap := map[string]interface{}{
		"processTypes": OtvProcessOption(),
		"message":      "",
		"results":      results,
	}
	r.HTML(200, "admin/otv", newmap)
}

func OtvProcessHandler(db *sql.DB, r render.Render, req *http.Request) {
	processType := req.FormValue("processType")
	var message string
	var results []*OtvShowListItem
	otv := &Otv{Db: db}
	switch processType {
	case "modified":
		results = otv.UpdateModified()
		message = "Update modified date complete."
	case "existing":
		results = otv.CheckOtvExisting()
		message = "Check Otv Existing complete."
	default:
		message = "Please Select Process"
	}

	processOptions := OtvProcessOption()

	for i, process := range processOptions {
		processOptions[i].Checked = (process.Value == processType)
	}

	newmap := map[string]interface{}{
		"processTypes": processOptions,
		"message":      message,
		"results":      results,
	}

	r.HTML(200, "admin/otv", newmap)

}

func BotVideoHandler(db *sql.DB, r render.Render, req *http.Request) {
	username := req.FormValue("username")
	status, _ := strconv.Atoi(req.FormValue("status"))
	page, _ := strconv.Atoi(req.FormValue("page"))
	formSearch := &FormSearchBotUser{username, status, int32(page)}

	b := NewBotVideo(db)
	botStatuses := b.getBotStatuses(formSearch.Status)
	botUsers := b.getBotUsers(formSearch.Username)
	// botVideos := b.getBotVideos(formSearch)

	newmap := map[string]interface{}{
		"formSearch":  formSearch,
		"botStatuses": botStatuses,
		"botUsers":    botUsers,
		// "botVideos":   botVideos,
	}

	r.HTML(200, "admin/botvideo", newmap)
}

func BotVideoPostHandler(db *sql.DB, r render.Render, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		//handle error http.Error() for example
		return
	}

	log.Println(req.Form)

	b := NewBotVideo(db)
	botVideos := req.Form["bot_video[]"]
	updateStatus := req.FormValue("update_status")
	b.setBotVideosStatus(botVideos, updateStatus)

	BotVideoHandler(db, r, req)
}

func BotVideoJSONHandler(db *sql.DB, r render.Render, req *http.Request) {
	username := req.FormValue("username")
	status, _ := strconv.Atoi(req.FormValue("status"))
	page, _ := strconv.Atoi(req.FormValue("page"))
	formSearch := &FormSearchBotUser{username, status, int32(page)}

	b := NewBotVideo(db)
	botVideos := b.getBotVideos(formSearch)
	r.JSON(200, botVideos)
}

func YoutubeHandler(db *sql.DB, r render.Render, req *http.Request) {
	username := req.FormValue("username")
	status, _ := strconv.Atoi(req.FormValue("status"))
	page, _ := strconv.Atoi(req.FormValue("page"))
	formSearch := &FormSearchBotUser{username, status, int32(page)}

	b := NewBotVideo(db)
	botStatuses := b.getBotStatuses(formSearch.Status)
	botUsers := b.getBotUsers(formSearch.Username)
	// botVideos := b.getBotVideos(formSearch)

	newmap := map[string]interface{}{
		"formSearch":  formSearch,
		"botStatuses": botStatuses,
		"botUsers":    botUsers,
		// "botVideos":   botVideos,
	}

	r.HTML(200, "admin/youtube", newmap)
}

func YoutubeJSONHandler(db *sql.DB, r render.Render, req *http.Request) {
	username := req.FormValue("username")
	y := bot.NewYoutube()
	_, youtubeVideos := y.GetVideoByUser(username, 1, 40)
	r.JSON(200, youtubeVideos)
}

func ShowJSONHandler(db *sql.DB, r render.Render, req *http.Request) {

}

func ShowListHandler(db *sql.DB, r render.Render, req *http.Request) {
	showID, _ := strconv.Atoi(req.FormValue("show_id"))
	showList := &ShowList{db}
	shows := showList.getProgram()
	results := showList.getData(showID)
	newmap := map[string]interface{}{
		"shows":   shows,
		"results": results,
	}

	r.HTML(200, "admin/showlist", newmap)
}
