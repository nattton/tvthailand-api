package admin

import (
	"database/sql"
	"encoding/base64"
	"strings"
)

func updateVideoEncrypt(db *sql.DB, listId int, videoId string) {
	encryptId := EncryptVideo(videoId)
	_, err := db.Exec("UPDATE tv_programlist SET programlist_youtube_encrypt = ? WHERE programlist_id = ?", encryptId, listId)
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
