package admin

import (
	"database/sql"
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const MTHAI_URL_FMT = "http://video.mthai.com/cool/player/%s.html"

type EncryptResult struct {
	ShowId  int
	ListId  int
	VideoId string
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
		MthaiParseUrl(db, r)
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
			showId  int
			listId  int
			videoId string
		)
		if err := rows.Scan(&showId, &listId, &videoId); err != nil {
			log.Fatal(err)
		}
		updateVideoEncrypt(db, listId, videoId)
		results = append(results, &EncryptResult{showId, listId, videoId})
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

func MthaiParseUrl(db *sql.DB, r render.Render) {
	var results []*EncryptResult
	rows, err := db.Query("SELECT program_id, programlist_id, programlist_youtube youtubeKey FROM tv_programlist WHERE programlist_banned = 0 AND programlist_src_type = 14 ORDER BY programlist_id DESC")
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			programId     int
			programlistId int
			youtubeKey    string
		)

		if err := rows.Scan(&programId, &programlistId, &youtubeKey); err != nil {
			log.Fatal(err)
		}

		fmt.Println(programId, programlistId, youtubeKey)
		videoIds := strings.Split(youtubeKey, ",")

		var videoUrls []string
		for _, videoId := range videoIds {
			videoUrl := fmt.Sprintf(MTHAI_URL_FMT, videoId)
			videoUrls = append(videoUrls, videoUrl)
		}
		videoResult := strings.Join(videoUrls, ",")
		fmt.Println("videoResult :", videoResult)
		encryptResult := EncryptVideo(videoResult)
		fmt.Println("encryptResult :", encryptResult)

		_, err := db.Exec("UPDATE tv_programlist SET programlist_youtube = ?, programlist_youtube_encrypt = ?, programlist_src_type = 11, mthai_video = ? WHERE programlist_id = ? ORDER BY program_id DESC", videoResult, encryptResult, youtubeKey, programlistId)
		if err != nil {
			panic(err)
		}
		results = append(results, &EncryptResult{programId, programlistId, youtubeKey})
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
