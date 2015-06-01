package bot

import (
	"database/sql"
	"fmt"

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
	ChannelID   string
	Description string
	UserType    string
	BotLimit    int
}

func NewBot(db *sql.DB) *Bot {
	b := new(Bot)
	b.Db = db
	return b
}

func (b *Bot) CheckRobotChannel() {
	youtubeUsers := b.getYoutubeRobotChannels()
	for _, youtubeUser := range youtubeUsers {
		fmt.Println(youtubeUser.Username)
		y := NewYoutube()
		_, youtubeVideos, _, _ := y.GetVideoByChannelID(youtubeUser.Username, youtubeUser.ChannelID, youtubeUser.BotLimit, "")
		for _, video := range youtubeVideos {
			fmt.Println(video.Username, video.Title, video.VideoID)
			b.checkBotVideoExistingAndAddBot(video)
		}
	}
}

func (b *Bot) CheckAllYoutubeUser() {
	youtubeUsers := b.getYoutubeRobotUsers()
	for _, youtubeUser := range youtubeUsers {
		fmt.Println(youtubeUser.Username)
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
		fmt.Println(video.Username, video.Title, video.VideoID)
		b.checkBotVideoExistingAndAddBot(video)
	}
}

func (b *Bot) CheckVideoInChannel(username string, channelID string) {
	y := NewYoutube()
	botLimit := 50

	nextPageToken := ""
	for {
		_, youtubeVideos, _, nextToken := y.GetVideoByChannelID(username, channelID, botLimit, nextPageToken)
		nextPageToken = nextToken
		for _, video := range youtubeVideos {
			fmt.Println(video.Username, video.Title, video.VideoID)
			b.checkBotVideoExistingAndAddBot(video)
		}
		if nextPageToken == "" {
			break
		}
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
		// total, _ = y.GetVideoByUser(username, 1, botLimit)
	}

	totalLoop := (total / botLimit) + 1
	for i := 0; i < totalLoop; i++ {
		b.CheckYoutubeUserAndKeyword(username, (i*botLimit)+1, botLimit, keyword)
	}
}

func (b *Bot) getYoutubeRobotChannels() (youtubeUsers []YoutubeUser) {
	rows, err := b.Db.Query("SELECT username, channel_id, description, user_type, bot_limit FROM tv_youtube_users WHERE channel_id != '' AND bot = 1 ORDER BY official DESC, username ASC")
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()
	for rows.Next() {
		youtubeUser := YoutubeUser{}
		if err := rows.Scan(&youtubeUser.Username, &youtubeUser.ChannelID, &youtubeUser.Description, &youtubeUser.UserType, &youtubeUser.BotLimit); err != nil {
			fmt.Println(err)
		}
		youtubeUsers = append(youtubeUsers, youtubeUser)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}
	return youtubeUsers
}

func (b *Bot) getYoutubeRobotUsers() []*YoutubeUser {
	q := "SELECT username, channel_id, description, user_type, bot_limit FROM tv_youtube_users WHERE bot = 1 ORDER BY username"
	return b.queryYoutubeUsers(q)
}

func (b *Bot) getEmptyChannel() []*YoutubeUser {
	q := "SELECT username, channel_id, description, user_type, bot_limit FROM tv_youtube_users WHERE bot = 1 AND user_type = 'user' AND channel_id = ''"
	return b.queryYoutubeUsers(q)
}

func (b *Bot) queryYoutubeUsers(q string) []*YoutubeUser {
	var youtubeUsers = []*YoutubeUser{}
	rows, err := b.Db.Query(q)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			username    string
			channelID   string
			description string
			userType    string
			botLimit    int
		)
		if err := rows.Scan(&username, &channelID, &description, &userType, &botLimit); err != nil {
			fmt.Println(err)
		}
		youtubeUser := &YoutubeUser{username, channelID, description, userType, botLimit}
		youtubeUsers = append(youtubeUsers, youtubeUser)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}
	return youtubeUsers
}

func (b *Bot) FindChannel() {
	y := NewYoutube()
	youtubeUsers := b.getEmptyChannel()
	for _, youtubeUser := range youtubeUsers {
		channelID := y.GetChannelIDByUser(youtubeUser.Username)
		fmt.Printf("Username, %s, ChannelID : %s\n", youtubeUser.Username, channelID)
		_, err := b.Db.Exec("UPDATE tv_youtube_users SET channel_id = ? WHERE username = ?", channelID, youtubeUser.Username)
		if err != nil {
			panic(err)
		}
	}
}

func (b *Bot) checkBotVideoExistingAndAddBot(video *YoutubeVideo) {
	rows, err := b.Db.Query("SELECT id from tv_bot_videos WHERE video_id = ?", video.VideoID)
	if err != nil {
		fmt.Println(err)
	} else {
		defer rows.Close()

		if !rows.Next() {
			row2s, err := b.Db.Query("SELECT programlist_id from tv_programlist WHERE programlist_youtube LIKE ? ", "%"+video.VideoID+"%")
			if err != nil {
				fmt.Println(err)
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

func (b *Bot) CheckKrobkruakao(start int) {
	krobkruakaos := Krobkruakaos(start)
	for _, kr := range krobkruakaos {
		fmt.Printf("%s - %s\nShort Url : %s, Date : %s\n", kr.Title, kr.Url, kr.ShortUrl, kr.Date)
		b.checkBotKrobkruakaoExistingAndAddBot(kr)
	}
}

func (b *Bot) checkBotKrobkruakaoExistingAndAddBot(video *Krobkruakao) {
	rows, err := b.Db.Query("SELECT id from tv_bot_videos WHERE url = ?", video.Url)
	if err != nil {
		fmt.Println(err)
	} else {
		defer rows.Close()

		if !rows.Next() {
			row2s, err := b.Db.Query("SELECT programlist_id from tv_programlist WHERE programlist_youtube LIKE ? ", "%"+video.ShortUrl+"%")
			if err != nil {
				fmt.Println(err)
			}
			defer row2s.Close()
			if row2s.Next() {
				video.Status = 1
			} else {
				video.Status = 0
			}
			b.insertBotKrobkruakao(video)
		}
	}
}

func (b *Bot) insertBotVideo(video *YoutubeVideo) {
	_, err := b.Db.Exec("INSERT INTO tv_bot_videos (username, title, video_id, video_type, published, status) VALUES (?, ?, ?, 'youtube', ?, ?)", video.Username, video.Title, video.VideoID, video.Published, video.Status)
	if err != nil {
		panic(err)
	}
	fmt.Println("Insert Bot Video ### ", video.Username, video.Title, video.VideoID, video.Published, "###")
}

func (b *Bot) insertBotKrobkruakao(video *Krobkruakao) {
	_, err := b.Db.Exec("INSERT INTO tv_bot_videos (username, title, url, video_id, video_type, published, status) VALUES ('krobkruakao3', ?, ?, ?, 'url', NOW(), ?)", video.Title+" | "+video.Date, video.Url, video.ShortUrl, video.Status)
	if err != nil {
		panic(err)
	}
	fmt.Println("Insert Bot Krobkruakao ### ", video.Title, video.Url, video.ShortUrl, video.Date, "###")
}
