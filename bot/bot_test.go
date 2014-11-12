package bot

import (
	"github.com/code-mobi/tvthailand-api/utils"
	"testing"
)

func TestGetYoutubeRobotUsers(t *testing.T) {
	db, _ := utils.OpenDB()
	b := NewBot(db)
	youtubeUsers := b.getYoutubeRobotUsers()

	if youtubeUsers == nil {
		t.Error("It should not be nil")
	}

	if len(youtubeUsers) == 0 {
		t.Error("It should length > 0")
	}
}
