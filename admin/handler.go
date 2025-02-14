package admin

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/code-mobi/tvthailand-api/data"
	"github.com/code-mobi/tvthailand-api/utils"
	"github.com/code-mobi/tvthailand-api/youtube"
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
	} else if idType == "empty" {

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
	case "empty":
		rows, _ = db.Query("SELECT show_id, id, video FROM episodes WHERE video_encrypt = '' ORDER BY episodes.id DESC")
	case "showid":
		rows, _ = db.Query("SELECT show_id, id, video FROM episodes WHERE show_id = ? ORDER BY episodes.id DESC", id)
	case "listid":
		rows, _ = db.Query("SELECT show_id, id, video FROM episodes WHERE id = ? ORDER BY episodes.id DESC", id)
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
	rows, err := db.Query("SELECT show_id, id, video FROM episodes WHERE banned = 0 AND src_type = 14 ORDER BY id DESC")
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

		_, err := db.Exec("UPDATE episodes SET video = ?, video_encrypt = ?, src_type = 11, mthai_video = ? WHERE id = ? ORDER BY program_id DESC", videoResult, encryptResult, youtubeKey, programlistID)
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
	case "findembed":
		results = otv.FindEmbed()
		message = "Find Embed complete."
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

func BotVideoHandler(r render.Render, req *http.Request) {
	dbg, err := utils.OpenGormDB()
	if err != nil {
		panic(err.Error())
	}
	defer dbg.Close()

	username := req.FormValue("channelId")
	q := req.FormValue("q")
	status, _ := strconv.Atoi(req.FormValue("status"))
	page, _ := strconv.Atoi(req.FormValue("page"))
	isOrderTitle := req.FormValue("isOrderTitle") == "true"
	formSearch := data.FormSearchBotUser{username, q, status, int32(page), isOrderTitle}

	botStatuses := data.GetBotVideoStatuses(formSearch.Status)
	botUsers := data.GetBotVideoUsers(dbg, formSearch.ChannelID)
	// botVideos := b.getBotVideos(formSearch)

	newmap := map[string]interface{}{
		"formSearch":  formSearch,
		"botStatuses": botStatuses,
		"botUsers":    botUsers,
		// "botVideos":   botVideos,
	}

	r.HTML(200, "admin/botvideo", newmap)
}

func BotVideoPostHandler(r render.Render, req *http.Request) {
	dbg, err := utils.OpenGormDB()
	if err != nil {
		panic(err.Error())
	}
	defer dbg.Close()

	if err := req.ParseForm(); err != nil {
		//handle error http.Error() for example
		return
	}

	botVideos := req.Form["bot_video[]"]
	updateStatus := req.FormValue("update_status")
	data.SetBotVideosStatus(dbg, botVideos, updateStatus)

	BotVideoHandler(r, req)
}

func BotVideoJSONHandler(r render.Render, req *http.Request) {
	dbg, err := utils.OpenGormDB()
	if err != nil {
		panic(err.Error())
	}
	defer dbg.Close()

	channelId := req.FormValue("channelId")
	q := req.FormValue("q")
	status, _ := strconv.Atoi(req.FormValue("status"))
	page, _ := strconv.Atoi(req.FormValue("page"))
	isOrderTitle := req.FormValue("isOrderTitle") == "true"
	formSearch := data.FormSearchBotUser{
		ChannelID:    channelId,
		Q:            q,
		Status:       status,
		Page:         int32(page),
		IsOrderTitle: isOrderTitle,
	}

	botVideos := data.GetBotVideos(dbg, formSearch)
	r.JSON(200, botVideos)
}

func YoutubeSearchChannelJSONHandler(db *sql.DB, r render.Render, req *http.Request, w http.ResponseWriter) {
	callback := req.FormValue("callback")
	channelID := req.FormValue("channelId")
	q := req.FormValue("q")
	maxResults, atoiErr := strconv.Atoi(req.FormValue("maxResults"))
	if atoiErr != nil {
		maxResults = 40
	}
	pageToken := req.FormValue("pageToken")

	y := youtube.NewYoutube()
	api, err := y.GetVideoJsonByChannelID(channelID, q, maxResults, pageToken)
	if err != nil {
		utils.JSONP(w, 404, callback, err)
	} else {
		utils.JSONP(w, 200, callback, api)
	}
}

func YoutubePlaylistItemJSONHandler(db *sql.DB, r render.Render, req *http.Request, w http.ResponseWriter) {
	callback := req.FormValue("callback")
	playlistID := req.FormValue("playlistId")
	maxResults, atoiErr := strconv.Atoi(req.FormValue("maxResults"))
	if atoiErr != nil {
		maxResults = 40
	}
	pageToken := req.FormValue("pageToken")

	y := youtube.NewYoutube()
	api, err := y.GetVideoJsonByPlaylistID(playlistID, maxResults, pageToken)
	if err != nil {
		utils.JSONP(w, 404, callback, err)
	} else {
		utils.JSONP(w, 200, callback, api)
	}
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

func KrobkruakaoHandler(r render.Render) {
	r.HTML(200, "admin/krobkruakao", map[string]interface{}{})
}

func KrobkruakaoJSONHandler(db *sql.DB, r render.Render, req *http.Request) {
	start, _ := strconv.Atoi(req.FormValue("start"))
	r.JSON(200, Krobkruakaos(start))
}
