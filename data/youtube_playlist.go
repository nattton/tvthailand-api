package data

import (
	"fmt"
	_ "github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand-api/youtube"
	"log"
)

type YoutubePlaylist struct {
	ChannelID  string
	Title      string
	ProgramID  int
	PlaylistID string
	BotEnabled bool
	BotLimit   int
}

func BotEnabledPlaylists(db *gorm.DB) (playlists []YoutubePlaylist, err error) {
	err = db.Where("bot_enabled = ?", true).Find(&playlists).Error
	return
}

func RunBotPlaylists(db *gorm.DB) {
	playlists, _ := BotEnabledPlaylists(db)
	for _, playlist := range playlists {
		fmt.Println(playlist.Title, playlist.PlaylistID)
		playlist.RunBot(db)
	}
}

func (p YoutubePlaylist) RunBot(db *gorm.DB) {
	y := youtube.NewYoutube()
	youtubePlaylist, err := y.GetVideoJsonByPlaylistID(p.PlaylistID, p.BotLimit, "")
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range youtubePlaylist.Items {
		AddBotVideoPlaylist(db, p, item)
	}
}
