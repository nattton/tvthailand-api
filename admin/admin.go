package admin

import (
	"database/sql"
	"encoding/base64"
	"strings"
)

const limitRow int32 = 40

func updateVideoEncrypt(db *sql.DB, listID int, videoID string) {
	encryptID := EncryptVideo(videoID)
	_, err := db.Exec("UPDATE tv_programlist SET programlist_youtube_encrypt = ? WHERE programlist_id = ?", encryptID, listID)
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

func EncryptVideo(videoID string) string {
	encrypt := base64.StdEncoding.EncodeToString([]byte(videoID))
	return re.Replace(encrypt)
}
