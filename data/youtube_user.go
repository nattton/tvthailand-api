package data

import (
	_ "github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

type YoutubeUser struct {
	ChannelID   string `json:"channelId"`
	Username    string `json:"username"`
	UserType    string `json:"userType"`
	Description string `json:"description"`
	ProgramID   int    `json:"programId"`
	BotEnabled  bool   `json:"botEnabled"`
	BotLimit    int    `json:"botLimit"`
	Official    bool   `json:"isOfficial"`
}

func BotEnabledUsers(db *gorm.DB) (users []YoutubeUser, err error) {
	err = db.Where("bot_enabled = ?", true).Find(&users).Error
	return
}

func UserByChannelID(db *gorm.DB, channelID string) (user YoutubeUser, err error) {
	err = db.Where("channel_id = ?", channelID).First(&user).Error
	return
}
