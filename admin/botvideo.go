package admin

import (
	"database/sql"
	"log"
)

type BotVideo struct {
	Db *sql.DB
}

func NewBotVideo(db *sql.DB) *BotVideo {
	b := new(BotVideo)
	b.Db = db
	return b
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

type BotVideoRow struct {
	Id        int32
	Username  string
	Title     string
	VideoId   string
	Published string
	Status    int
}

func (b *BotVideo) getBotStatuses(id int) []*BotStatus {
	botStatuses := []*BotStatus{}
	botStatuses = append(botStatuses, &BotStatus{0, "Waiting", (id == 0)})
	botStatuses = append(botStatuses, &BotStatus{1, "Updated", (id == 1)})
	botStatuses = append(botStatuses, &BotStatus{-1, "Rejected", (id == -1)})
	return botStatuses
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

func (b *BotVideo) getBotVideos(qUsername string) []*BotVideoRow {
	botVideos := []*BotVideoRow{}

	var rows *sql.Rows
	var err error

	if qUsername == "" {
		rows, err = b.Db.Query("SELECT id, username, title, video_id, published, status from tv_bot_videos ORDER BY username, published LIMIT 0, 50")
	} else {
		rows, err = b.Db.Query("SELECT id, username, title, video_id, published, status from tv_bot_videos WHERE username = ? ORDER BY published", qUsername)
	}

	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id        int32
			username  string
			title     string
			videoId   string
			published string
			status 		int
		)
		if err := rows.Scan(&id, &username, &title, &videoId, &published, &status); err != nil {
			log.Fatal(err)
		}
		botVideo := &BotVideoRow{id, username, title, videoId, published, status}
		botVideos = append(botVideos, botVideo)
	}

	return botVideos
}
