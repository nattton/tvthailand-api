package bot

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Bot struct {
	Db *sql.DB
}

type YoutubeUsers struct {
	Users []*YoutubeUser
}

type YoutubeUser struct {
	Username    string
	Description string
	UserType    string
}

func NewBot(db *sql.DB) *Bot {
	b := new(Bot)
	b.Db = db
	return b
}

func (b *Bot) CheckYoutubeUser() {
	youtubeUsers := b.getYoutubeRobotUsers()
	y := &Youtube{}
	for _, youtubeUser := range youtubeUsers {
		log.Println(youtubeUser.Username)
		youtubeVideos := y.getVideoByUser(youtubeUser.Username)
		for _, video := range youtubeVideos {
			log.Println(video.Username, video.Title, video.VideoId)
			if b.checkBotVideoExisting(video) {
				b.insertBotVideo(video)
			}
		}
	}
}

func (b *Bot) getYoutubeRobotUsers() []*YoutubeUser {
	var youtubeUsers = []*YoutubeUser{}
	rows, err := b.Db.Query("SELECT username, description, user_type FROM tv_youtube_users WHERE bot = 1 ORDER BY username LIMIT 0,5")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			username    string
			description string
			userType    string
		)
		if err := rows.Scan(&username, &description, &userType); err != nil {
			log.Fatal(err)
		}
		youtubeUser := &YoutubeUser{username, description, userType}
		youtubeUsers = append(youtubeUsers, youtubeUser)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return youtubeUsers
}

func (b *Bot) checkBotVideoExisting(video *YoutubeVideo) bool {
	return true
}

func (b *Bot) insertBotVideo(video *YoutubeVideo) {
	_, err := b.Db.Exec("INSERT INTO tv_bot_videos (username, title, video_id, video_type, is_updated) VALUES (?, ?, ?, 'youtube', 0)", video.Username, video.Title, video.VideoId)
	if err != nil {
		panic(err)
	}
}
