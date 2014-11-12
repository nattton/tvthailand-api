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
	BotLimit    int
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
		youtubeVideos := y.getVideoByUser(youtubeUser.Username, youtubeUser.BotLimit)
		for _, video := range youtubeVideos {
			log.Println(video.Username, video.Title, video.VideoId)
			b.checkBotVideoExistingAndAddBot(video)
		}
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
	rows, err := b.Db.Query("SELECT id from tv_bot_videos WHERE video_id = ?", video.VideoId)
	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()

		if !rows.Next() {
			row2s, err := b.Db.Query("SELECT programlist_id from tv_programlist WHERE programlist_youtube LIKE ? ", "%"+video.VideoId+"%")
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
	_, err := b.Db.Exec("INSERT INTO tv_bot_videos (username, title, video_id, video_type, published, status) VALUES (?, ?, ?, 'youtube', ?, ?)", video.Username, video.Title, video.VideoId, video.Published, video.Status)
	if err != nil {
		panic(err)
	}
	log.Println("Insert Bot Video ### ", video.Username, video.Title, video.VideoId, video.Published, "###")
}
