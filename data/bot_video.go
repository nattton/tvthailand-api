package data

import (
	"fmt"
	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand-api/youtube"
	"log"
	"math"
	"sync"
	"time"
)

const DateFMT = "2006-01-02T15:04:05"

type BotVideo struct {
	ID          int       `json:"id"`
	ChannelID   string    `json:"channelId"`
	PlaylistID  string    `json:"playlistId"`
	Title       string    `json:"title"`
	VideoID     string    `json:"videoId"`
	VideoType   string    `json:"videoType"`
	PublishedAt time.Time `json:"publishedAt"`
	Status      int       `json:"status"`
}

type BotVideoDistinct struct {
	Username string
}

func MigrateUsernameToChannelID(db *gorm.DB) {
	var botVideos []BotVideoDistinct
	db.Model(BotVideo{}).Where("channel_id = ?", "").Select("DISTINCT(username)").Order("published DESC").Scan(&botVideos)
	for _, botVideo := range botVideos {
		var youtubeUser YoutubeUser
		err := db.Where("username = ?", botVideo.Username).First(&youtubeUser).Error
		if err != nil {
			panic(err)
		} else {
			db.Model(BotVideo{}).Where("username = ?", botVideo.Username).Updates(BotVideo{ChannelID: youtubeUser.ChannelID})
		}
	}
}

func GetBotVideoByVideoID(db *gorm.DB, videoID string) (botVideo BotVideo, err error) {
	err = db.Where("video_id = ?", videoID).First(&botVideo).Error
	return
}

func AddBotVideoPlaylist(db *gorm.DB, wg *sync.WaitGroup, throttle chan int, pl YoutubePlaylist, item youtube.PlaylistItem) {
	defer wg.Done()
	status := 0
	publishedAt, err := time.Parse(time.RFC3339Nano, item.Snippet.PublishedAt)
	if err != nil {
		log.Fatal(err)
	}

	episode, _ := GetEpisodeByVideoID(db, item.Snippet.ResourceID.VideoID)
	if episode.ID > 0 {
		status = 1
	}

	botVideo, _ := GetBotVideoByVideoID(db, item.Snippet.ResourceID.VideoID)
	if botVideo.ID == 0 {
		botVideo = BotVideo{
			ChannelID:   pl.ChannelID,
			PlaylistID:  pl.PlaylistID,
			Title:       item.Snippet.Title,
			VideoID:     item.Snippet.ResourceID.VideoID,
			VideoType:   "youtube",
			PublishedAt: publishedAt,
			Status:      status,
		}
		db.Create(&botVideo)
	} else {
		if botVideo.PlaylistID == "" {
			botVideo.PlaylistID = item.ID
			db.Save(&botVideo)
		}
	}

	fmt.Println(botVideo.ChannelID, botVideo.Title, botVideo.PublishedAt)
	<-throttle
}

type FormSearchBotUser struct {
	ChannelID    string
	Q            string
	Status       int
	Page         int32
	IsOrderTitle bool
}

type BotVideos struct {
	Videos      []*BotVideoRow `json:"videos"`
	CountRow    int32          `json:"countRow"`
	CurrentPage int32          `json:"currentPage"`
	MaxPage     int32          `json:"maxPage"`
}

type BotVideoRow struct {
	ID                int32     `json:"id"`
	ChannelID         string    `json:"channelId"`
	Description       string    `json:"description"`
	ProgramID         int64     `json:"programId"`
	UserType          string    `json:"userType"`
	Title             string    `json:"title"`
	VideoID           string    `json:"videoId"`
	VideoType         string    `json:"videoType"`
	PublishedAt       time.Time `json:"publishedAt"`
	Status            int       `json:"status"`
	PlaylistTitle     string    `json:"-"`
	PlaylistProgramID int64     `json:"-"`
}

func GetBotVideos(db *gorm.DB, f FormSearchBotUser) BotVideos {
	var countRow int32
	botVideos := []*BotVideoRow{}
	dbQ := db.Table("bot_videos").Where("bot_videos.status = ? AND bot_videos.title LIKE ?", f.Status, "%"+f.Q+"%").Select("bot_videos.id, bot_videos.channel_id, youtube_users.description, youtube_users.program_id, youtube_users.user_type, bot_videos.title, video_id, video_type, DATE_ADD(bot_videos.published_at, INTERVAL 7 HOUR) published_at, bot_videos.status, youtube_playlists.title playlist_title, youtube_playlists.program_id playlist_program_id").Joins("LEFT JOIN youtube_users ON bot_videos.channel_id = youtube_users.channel_id LEFT JOIN youtube_playlists ON bot_videos.playlist_id = youtube_playlists.playlist_id").Order("youtube_users.official DESC, bot_videos.channel_id ASC")
	if f.ChannelID == "all" || f.ChannelID == "" {
		dbQ.Count(&countRow)
	} else {
		dbQ = dbQ.Where("bot_videos.channel_id = ?", f.ChannelID)
		dbQ.Count(&countRow)
	}

	if f.IsOrderTitle {
		dbQ = dbQ.Order("bot_videos.title ASC, published_at DESC")
	}

	err := dbQ.Offset(f.Page * limitRow).Limit(limitRow).Scan(&botVideos).Error

	for index, _ := range botVideos {
		if botVideos[index].PlaylistProgramID > 0 {
			botVideos[index].ProgramID = botVideos[index].PlaylistProgramID
		}
		if botVideos[index].PlaylistTitle != "" {
			botVideos[index].Title = fmt.Sprintf("%s | %s", botVideos[index].PlaylistTitle, botVideos[index].Title)
		}
	}

	if err != nil {
		panic(err)
	}

	return BotVideos{
		Videos:   botVideos,
		CountRow: countRow, CurrentPage: f.Page,
		MaxPage: int32(math.Ceil(float64(countRow / limitRow))),
	}
}
