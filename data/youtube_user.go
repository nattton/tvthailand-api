package data

import (
	"fmt"
	_ "github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand-api/youtube"
	"log"
	"sync"
	"time"
)

type YoutubeUser struct {
	Username    string `json:"username" gorm:"primary_key"`
	ChannelID   string `json:"channelId"`
	Description string `json:"description"`
	UserType    string `json:"userType"`
	ProgramID   int    `json:"programId"`
	BotEnabled  bool   `json:"botEnabled"`
	BotLimit    int    `json:"botLimit"`
	Official    bool   `json:"isOfficial"`

	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func UpdateUserChannel(db *gorm.DB) {
	var users []YoutubeUser
	db.Where("user_type = ? AND channel_id = ?", "user", "").Find(&users)
	for _, user := range users {
		y := youtube.NewYoutube()
		user.ChannelID = y.ChannelIDByUser(user.Username)
		fmt.Printf("Username, %s, ChannelID : %s\n", user.Username, user.ChannelID)
		db.Save(&user)
	}
}

func CheckActiveUser(db *gorm.DB) {
	var users []YoutubeUser
	db.Order("updated_at ASC").Find(&users)
	for _, user := range users {
		y := youtube.NewYoutube()
		if !y.ChannelIsActive(user.ChannelID) {
			fmt.Printf("Username: %s, ChannelID: %s, Description: %s\n", user.Username, user.ChannelID, user.Description)
			DeleteBotVideoByChannel(db, user.ChannelID)
			db.Delete(&user)
		} else {
			db.Model(&user).UpdateColumns(YoutubeUser{UpdatedAt: time.Now()})
		}
	}
}

func BotEnabledUsers(db *gorm.DB) (users []YoutubeUser, err error) {
	err = db.Where("bot_enabled = ?", true).Find(&users).Error
	return
}

func UserByChannelID(db *gorm.DB, channelID string) (user YoutubeUser, err error) {
	err = db.Where("channel_id = ?", channelID).First(&user).Error
	return
}

func RunBotChannel(db *gorm.DB, channelId string, query string) {
	var user YoutubeUser
	err := db.Where("channel_id = ?", channelId).First(&user).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.Description, user.ChannelID)
	user.RunBot(db, true, query, "")
}

func RunBotChannels(db *gorm.DB) {
	users, _ := BotEnabledUsers(db)
	for _, user := range users {
		fmt.Println(user.Description, user.ChannelID)
		user.RunBot(db, false, "", "")
	}
}

func (user YoutubeUser) RunBot(db *gorm.DB, continuous bool, query string, nextToken string) {
	var wg sync.WaitGroup
	limitRow := user.BotLimit
	if continuous {
		limitRow = 40
	}
	y := youtube.NewYoutube()
	youtube, err := y.GetVideoJsonByChannelID(user.ChannelID, query, limitRow, nextToken)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range youtube.Items {
		throttle <- 1
		wg.Add(1)
		go AddBotVideoChannel(db, &wg, throttle, user, item)
	}
	wg.Wait()

	if continuous {
		user.RunBot(db, continuous, query, youtube.NextPageToken)
	}
}
