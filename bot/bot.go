package bot

import (
	"database/sql"
	"fmt"
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
	youtubeUsers := getYoutubeRobotUsers(b.Db)
	for _, youtubeUser := range youtubeUsers {
		fmt.Println(youtubeUser.Username)
	}
}

func getYoutubeRobotUsers(db *sql.DB) []*YoutubeUser {
	var youtubeUsers = []*YoutubeUser{}
	rows, err := db.Query("SELECT username, description, user_type FROM tv_youtube_users WHERE bot = 1 ORDER BY username")
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
