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
			PlaylistID:  item.ID,
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
	Username     string
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
	ID          int32     `json:"id"`
	Username    string    `json:"username"`
	Description string    `json:"description"`
	ProgramID   int64     `json:"programId"`
	UserType    string    `json:"userType"`
	Title       string    `json:"title"`
	VideoID     string    `json:"videoId"`
	VideoType   string    `json:"videoType"`
	PublishedAt time.Time `json:"publishedAt"`
	Status      int       `json:"status"`
}

func GetBotVideos(db *gorm.DB, f FormSearchBotUser) BotVideos {
	var countRow int32
	botVideos := []*BotVideoRow{}
	dbQ := db.Table("bot_videos").Where("status = ? AND title LIKE ?", f.Status, "%"+f.Q+"%").Select("bot_videos.id, bot_videos.username, youtube_users.description, youtube_users.program_id, youtube_users.user_type, bot_videos.title, video_id, video_type, DATE_ADD(bot_videos.published_at, INTERVAL 7 HOUR), bot_videos.status").Joins("LEFT JOIN youtube_users ON (bot_videos.channel_id = youtube_users.channel_id)")
	if f.Username == "all" || f.Username == "" {
		dbQ.Count(&countRow)
	} else {
		dbQ = dbQ.Where("username = ?", "%"+f.Q+"%")
		dbQ.Count(&countRow)
	}

	if f.IsOrderTitle {
		dbQ = dbQ.Order("bot_videos.title ASC")
	}

	err := dbQ.Limit(limitRow).Scan(&botVideos).Error

	if err != nil {
		panic(err)
	}

	return BotVideos{
		Videos:   botVideos,
		CountRow: countRow, CurrentPage: f.Page,
		MaxPage: int32(math.Ceil(float64(countRow / limitRow))),
	}
}
