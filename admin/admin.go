package admin

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AdminHandler struct {
	Db *sql.DB
}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	len_path := len(paths)
	if len_path > 1 {
		switch c := paths[1]; c {
		case "encryptvideo":
			h.encryptShow_ServeHTTP(w, r)
			return
		}
	}

	http.Error(w, `Not found`, http.StatusNotFound)
}

func (h *AdminHandler) encryptShow_ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	showListId, _ := strconv.Atoi(r.URL.Query().Get("listid"))

	if id == 0 && showListId == 0 {
		http.Error(w, `Not found id or listid`, http.StatusNotFound)
		return
	}

	var rows *sql.Rows

	if id > 0 {
		rows, _ = h.Db.Query("SELECT program_id, programlist_id, programlist_youtube FROM tv_programlist WHERE program_id = ? ORDER BY programlist_id DESC", id)
	} else {
		rows, _ = h.Db.Query("SELECT program_id, programlist_id, programlist_youtube FROM tv_programlist WHERE programlist_id = ? ORDER BY programlist_id DESC", showListId)
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

		h.updateVideoEncrypt(listId, videoId)

		result := fmt.Sprintf("Update Encrypt %d %d %s\n", showId, listId, videoId)
		log.Println(result)
		fmt.Fprintf(w, result)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}

func (h *AdminHandler) updateVideoEncrypt(listId int, videoId string) {
	encryptId := EncryptVideo(videoId)
	_, err := h.Db.Exec("UPDATE tv_programlist SET programlist_youtube_encrypt = ? WHERE programlist_id = ?", encryptId, listId)
	if err != nil {
		panic(err)
	}
}

var re = strings.NewReplacer(
	"+", "-",
	"=", ",",
	"a", "!",
	"b", "@",
	"c", "#",
	"d", "$",
	"e", "%",
	"f", "^",
	"g", "&",
	"h", "*",
	"i", "(",
	"j", ")",
	"k", "{",
	"l", "}",
	"m", "[",
	"n", "]",
	"o", ":",
	"p", ";",
	"q", "<",
	"r", ">",
	"s", "?",
)

func EncryptVideo(videoId string) string {
	encrypt := base64.StdEncoding.EncodeToString([]byte(videoId))
	return re.Replace(encrypt)
}
