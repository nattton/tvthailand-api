package data

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/code-mobi/tvthailand-api/youtube"
	"github.com/jinzhu/gorm"
)

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

func AddBotVideoChannel(db *gorm.DB, wg *sync.WaitGroup, throttle chan int, user YoutubeUser, item youtube.Item) {
	defer wg.Done()
	status := 0
	publishedAt, err := time.Parse(time.RFC3339Nano, item.Snippet.PublishedAt)
	if err != nil {
		log.Fatal(err)
	}

	videoID := item.ContentDetails.Upload.VideoID
	botVideo, _ := GetBotVideoByVideoID(db, videoID)
	if botVideo.ID == 0 {
		if CheckExistingVideoInEpisode(db, videoID) {
			status = 1
		}
		botVideo = BotVideo{
			ChannelID:   user.ChannelID,
			Title:       item.Snippet.Title,
			VideoID:     videoID,
			VideoType:   "youtube",
			PublishedAt: publishedAt,
			Status:      status,
		}
		db.Create(&botVideo)
	}

	fmt.Println(botVideo.ChannelID, botVideo.Title, botVideo.PublishedAt)
	<-throttle
}

func AddBotVideoPlaylist(db *gorm.DB, wg *sync.WaitGroup, throttle chan int, pl YoutubePlaylist, item youtube.PlaylistItem) {
	defer wg.Done()
	// status := 0
	// publishedAt, err := time.Parse(time.RFC3339Nano, item.Snippet.PublishedAt)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // videoID := item.Snippet.ResourceID.VideoID
	// botVideo, _ := GetBotVideoByVideoID(db, videoID)
	// if botVideo.ID == 0 {
	// 	if CheckExistingVideoInEpisode(db, videoID) {
	// 		status = 1
	// 	}
	// 	botVideo = BotVideo{
	// 		ChannelID:   pl.ChannelID,
	// 		PlaylistID:  pl.PlaylistID,
	// 		Title:       item.Snippet.Title,
	// 		VideoID:     videoID,
	// 		VideoType:   "youtube",
	// 		PublishedAt: publishedAt,
	// 		Status:      status,
	// 	}
	// 	db.Create(&botVideo)
	// } else {
	// 	if botVideo.PlaylistID == "" {
	// 		botVideo.PlaylistID = item.ID
	// 		db.Save(&botVideo)
	// 	}
	// }

	// fmt.Println(botVideo.ChannelID, botVideo.Title, botVideo.PublishedAt)
	<-throttle
}

func CheckExistingVideoInEpisode(db *gorm.DB, videoID string) bool {
	count, _ := GetCountEpisodeByVideoID(db, videoID)
	return count > 0
}

type BotUser struct {
	ChannelID   string
	Description string
	IsSelected  bool
}

func GetBotVideoUsers(db *gorm.DB, selectUsername string) (botUsers []BotUser) {
	err := db.Table("youtube_users").
		Where("description != ? AND channel_id != ?", "", "").
		Select("channel_id, description").
		Order("description").
		Scan(&botUsers).Error
	if err != nil {
		panic(err)
	}
	for index := range botUsers {
		botUsers[index].IsSelected = (selectUsername == botUsers[index].ChannelID)
	}
	return
}

type FormSearchBotUser struct {
	ChannelID    string
	Q            string
	Status       int
	Page         int32
	IsOrderTitle bool
}

type BotVideos struct {
	Videos      []BotVideoRow `json:"videos"`
	CountRow    int32         `json:"countRow"`
	CurrentPage int32         `json:"currentPage"`
	MaxPage     int32         `json:"maxPage"`
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
	botVideos := []BotVideoRow{}
	dbQ := db.Table("bot_videos").
		Where("bot_videos.status = ? AND bot_videos.title LIKE ?", f.Status, "%"+f.Q+"%").
		Select("bot_videos.id, bot_videos.channel_id, youtube_users.description, youtube_users.program_id, youtube_users.user_type, bot_videos.title, video_id, video_type, DATE_ADD(bot_videos.published_at, INTERVAL 7 HOUR) published_at, bot_videos.status, youtube_playlists.title playlist_title, youtube_playlists.program_id playlist_program_id").
		Joins("LEFT JOIN youtube_users ON bot_videos.channel_id = youtube_users.channel_id LEFT JOIN youtube_playlists ON bot_videos.playlist_id = youtube_playlists.playlist_id").
		Order("youtube_users.official DESC, bot_videos.channel_id ASC")
	if f.ChannelID == "all" || f.ChannelID == "" {
		dbQ.Count(&countRow)
	} else {
		dbQ = dbQ.Where("bot_videos.channel_id = ?", f.ChannelID)
		dbQ.Count(&countRow)
	}

	if f.IsOrderTitle {
		dbQ = dbQ.Order("bot_videos.title ASC")
	} else {
		dbQ = dbQ.Order("bot_videos.published_at DESC")
	}

	err := dbQ.Offset(f.Page * limitRow).Limit(limitRow).Scan(&botVideos).Error

	for index := range botVideos {
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
		Videos:      botVideos,
		CountRow:    countRow,
		CurrentPage: f.Page,
		MaxPage:     int32(math.Ceil(float64(countRow / limitRow))),
	}
}

type BotStatus struct {
	ID         int32
	Name       string
	IsSelected bool
}

func GetBotVideoStatuses(id int) []BotStatus {
	botStatuses := []BotStatus{}
	botStatuses = append(botStatuses, BotStatus{0, "Waiting", (id == 0)})
	botStatuses = append(botStatuses, BotStatus{1, "Updated", (id == 1)})
	botStatuses = append(botStatuses, BotStatus{-1, "Rejected", (id == -1)})
	botStatuses = append(botStatuses, BotStatus{2, "Suspended", (id == 2)})
	return botStatuses
}

func GetBotStatusID(status string) int {
	switch status {
	case "Rejected":
		return -1
	case "Updated":
		return 1
	case "Suspended":
		return 2
	default:
		return 0
	}
}

func SetBotVideoStatus(db *gorm.DB, id []int, status int) {
	db.Model(BotVideo{}).Where("id in (?)", id).UpdateColumn("status", status)
}

func SetBotVideosStatus(db *gorm.DB, videoIDs []string, updateStatus string) {
	statusID := GetBotStatusID(updateStatus)
	var ids []int
	for _, videoID := range videoIDs {
		id, _ := strconv.Atoi(videoID)
		ids = append(ids, id)
	}
	SetBotVideoStatus(db, ids, statusID)
}

func DeleteBotVideoByChannel(db *gorm.DB, channelID string) {
	db.Where("channel_id = ?", channelID).Delete(BotVideo{})
}
