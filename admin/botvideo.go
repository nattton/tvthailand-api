package admin

import (
	"database/sql"
	"log"
	"strconv"
	"math"
)

type BotVideo struct {
	Db *sql.DB
}

func NewBotVideo(db *sql.DB) *BotVideo {
	b := new(BotVideo)
	b.Db = db
	return b
}

type FormSearchBotUser struct {
	Username string
	Status   int
	Page 		int32
}

type BotUser struct {
	Username    string
	Description string
	IsSelected  bool
}

type BotStatus struct {
	Id         int32
	Name       string
	IsSelected bool
}

type BotVideos struct {
	Videos []*BotVideoRow
	CountRow int32
	CurrentPage int32
	MaxPage int32
}

type BotVideoRow struct {
	Id          int32
	Username    string
	Description string
	UserType    string
	Title       string
	VideoId     string
	Published   string
	Status      int
}

func (b *BotVideo) getBotStatuses(id int) []*BotStatus {
	botStatuses := []*BotStatus{}
	botStatuses = append(botStatuses, &BotStatus{0, "Waiting", (id == 0)})
	botStatuses = append(botStatuses, &BotStatus{1, "Updated", (id == 1)})
	botStatuses = append(botStatuses, &BotStatus{-1, "Rejected", (id == -1)})
	return botStatuses
}

func (b *BotVideo) getBotStatusId(status string) int {
	switch status {
	case "Rejected":
		return -1
	case "Updated":
		return 1
	default:
		return 0
	}
}

func (b *BotVideo) getBotUsers(selectUsername string) []*BotUser {
	botUsers := []*BotUser{}
	rows, err := b.Db.Query("SELECT DISTINCT v.username, u.description from tv_bot_videos v LEFT JOIN tv_youtube_users u ON (v.username = u.username) ORDER BY description")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			username    string
			description string
		)
		if err := rows.Scan(&username, &description); err != nil {
			log.Fatal(err)
		}
		isSelected := selectUsername == username
		botUser := &BotUser{username, description, isSelected}
		botUsers = append(botUsers, botUser)
	}

	return botUsers
}

func (b *BotVideo) getBotVideos(f *FormSearchBotUser) *BotVideos {
	var countRow int32 = 0
	var limitRow int32 = 60
	botVideos := []*BotVideoRow{}

	var rows *sql.Rows
	var err error

	if f.Username == "all" || f.Username == "" {
		err = b.Db.QueryRow("SELECT count(id) from tv_bot_videos WHERE status = ?", f.Status).Scan(&countRow)
	} else {
		err = b.Db.QueryRow("SELECT count(id) from tv_bot_videos WHERE status = ? AND username = ?", f.Status, f.Username).Scan(&countRow)
	}

	if f.Username == "all" || f.Username == "" {
		rows, err = b.Db.Query("SELECT v.id, v.username, u.description, u.user_type, v.title, video_id, DATE_ADD(published, INTERVAL 7 HOUR), status from tv_bot_videos v LEFT JOIN tv_youtube_users u ON (v.username = u.username) WHERE status = ? ORDER BY v.username, published DESC LIMIT ?, ?", f.Status, (f.Page * limitRow), limitRow)
	} else {
		rows, err = b.Db.Query("SELECT v.id, v.username, u.description, u.user_type, v.title, video_id, DATE_ADD(published, INTERVAL 7 HOUR), status from tv_bot_videos v LEFT JOIN tv_youtube_users u ON (v.username = u.username) WHERE status = ? AND v.username = ? ORDER BY published DESC LIMIT ?, ?", f.Status, f.Username, (f.Page * limitRow), limitRow)
	}

	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id          int32
			username    string
			description string
			userType    string
			title       string
			videoId     string
			published   string
			status      int
		)
		if err := rows.Scan(&id, &username, &description, &userType, &title, &videoId, &published, &status); err != nil {
			log.Fatal(err)
		}
		botVideo := &BotVideoRow{id, username, description, userType, title, videoId, published, status}
		botVideos = append(botVideos, botVideo)
	}

	return &BotVideos{botVideos, countRow, f.Page, int32(math.Ceil(float64(countRow / limitRow)))}
	// return botVideos
}

func (b *BotVideo) setBotVideoStatus(id int, status int) {
	_, err := b.Db.Exec("UPDATE tv_bot_videos SET status = ? WHERE id = ?", status, id)
	if err != nil {
		panic(err)
	}
}

func (b *BotVideo) setBotVideosStatus(videoIds []string, updateStatus string) {
	statusId := b.getBotStatusId(updateStatus)
	for _, videoId := range videoIds {
		id, _ := strconv.Atoi(videoId)
		b.setBotVideoStatus(id, statusId)
	}
}
