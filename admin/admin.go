package admin

import (
	"database/sql"
	"encoding/base64"
	"strings"
)

const limitRow int32 = 40

func updateVideoEncrypt(db *sql.DB, listID int, videoID string) {
	encryptID := EncryptVideo(videoID)
	_, err := db.Exec("UPDATE episodes SET video_encrypt = ? WHERE id = ?", encryptID, listID)
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
