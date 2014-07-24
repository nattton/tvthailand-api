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
)

type EnryptResult struct {
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

	var results []*EnryptResult

	if idType == "" || id == 0 {
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
		results = append(results, &EnryptResult{showId, listId, videoId})
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

func OtvHandler(r render.Render) {
	var results []*OtvShowListItem
	newmap := map[string]interface{}{
		"message": "",
		"results": results,
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
		message = "Update modified date complete."
	default:
		message = "Please Select Process"
	}

	newmap := map[string]interface{}{
		"message": message,
		"results": results,
	}

	r.HTML(200, "admin/otv", newmap)

}
