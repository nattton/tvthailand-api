package data

import (
	"fmt"
	_ "github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand-api/youtube"
	"log"
	"sync"
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

func RunBotPlaylist(db *gorm.DB, playlistId string) {
	var playlist YoutubePlaylist
	err := db.Where("playlist_id = ?", playlistId).First(&playlist).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(playlist.Title, playlist.PlaylistID)
	playlist.RunBot(db, true, "")
}

func RunBotPlaylists(db *gorm.DB) {
	playlists, _ := BotEnabledPlaylists(db)
	for _, playlist := range playlists {
		fmt.Println(playlist.Title, playlist.PlaylistID)
		playlist.RunBot(db, false, "")
	}
}

func (p YoutubePlaylist) RunBot(db *gorm.DB, continuous bool, nextToken string) {
	var wg sync.WaitGroup
	y := youtube.NewYoutube()
	youtubePlaylist, err := y.GetVideoJsonByPlaylistID(p.PlaylistID, p.BotLimit, nextToken)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range youtubePlaylist.Items {
		throttle <- 1
		wg.Add(1)
		go AddBotVideoPlaylist(db, &wg, throttle, p, item)
	}
	wg.Wait()

	if continuous {
		p.RunBot(db, continuous, youtubePlaylist.NextPageToken)
	}
}
