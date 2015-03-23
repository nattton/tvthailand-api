package bot

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
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
	BotLimit    int
}

func NewBot(db *sql.DB) *Bot {
	b := new(Bot)
	b.Db = db
	return b
}

func (b *Bot) CheckAllYoutubeUser() {
	youtubeUsers := b.getYoutubeRobotUsers()
	for _, youtubeUser := range youtubeUsers {
		log.Println(youtubeUser.Username)
		b.CheckYoutubeUser(youtubeUser.Username, 1, youtubeUser.BotLimit)
	}
}

func (b *Bot) CheckYoutubeUser(username string, start int, botLimit int) {
	b.CheckYoutubeUserAndKeyword(username, start, botLimit, "")
}

func (b *Bot) CheckYoutubeUserAndKeyword(username string, start int, botLimit int, keyword string) {
	y := NewYoutube()
	_, youtubeVideos := y.GetVideoByUserAndKeyword(username, start, botLimit, keyword)
	for _, video := range youtubeVideos {
		log.Println(video.Username, video.Title, video.VideoID)
		b.checkBotVideoExistingAndAddBot(video)
	}
}

func (b *Bot) CheckAllVideoInYoutubeUserAndKeyword(username string, keyword string) {
	y := NewYoutube()
	botLimit := 50
	var total int
	if keyword != "" {
		fmt.Println(keyword)
		total, _ = y.GetVideoByUserAndKeyword(username, 1, botLimit, keyword)
	} else {
		total, _ = y.GetVideoByUser(username, 1, botLimit)
	}

	totalLoop := (total / botLimit) + 1
	for i := 0; i < totalLoop; i++ {
		b.CheckYoutubeUserAndKeyword(username, (i*botLimit)+1, botLimit, keyword)
	}
}

func (b *Bot) getYoutubeRobotUsers() []*YoutubeUser {
	var youtubeUsers = []*YoutubeUser{}
	rows, err := b.Db.Query("SELECT username, description, user_type, bot_limit FROM tv_youtube_users WHERE bot = 1 ORDER BY username")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			username    string
			description string
			userType    string
			botLimit    int
		)
		if err := rows.Scan(&username, &description, &userType, &botLimit); err != nil {
			log.Fatal(err)
		}
		youtubeUser := &YoutubeUser{username, description, userType, botLimit}
		youtubeUsers = append(youtubeUsers, youtubeUser)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return youtubeUsers
}

func (b *Bot) checkBotVideoExistingAndAddBot(video *YoutubeVideo) {
	rows, err := b.Db.Query("SELECT id from tv_bot_videos WHERE video_id = ?", video.VideoID)
	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()

		if !rows.Next() {
			row2s, err := b.Db.Query("SELECT programlist_id from tv_programlist WHERE programlist_youtube LIKE ? ", "%"+video.VideoID+"%")
			if err != nil {
				log.Fatal(err)
			}
			defer row2s.Close()
			if row2s.Next() {
				video.Status = 1
			} else {
				video.Status = 0
			}
			b.insertBotVideo(video)
		}
	}
}

func (b *Bot) insertBotVideo(video *YoutubeVideo) {
	_, err := b.Db.Exec("INSERT INTO tv_bot_videos (username, title, video_id, video_type, published, status) VALUES (?, ?, ?, 'youtube', ?, ?)", video.Username, video.Title, video.VideoID, video.Published, video.Status)
	if err != nil {
		panic(err)
	}
	log.Println("Insert Bot Video ### ", video.Username, video.Title, video.VideoID, video.Published, "###")
}
