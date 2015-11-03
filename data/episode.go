package data

import (
	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"time"
)

type Episode struct {
	ID        int    `gorm:"primary_key"`
	HashID    string `json:"-"`
	ShowID    int    `json:"-"`
	Ep        int    `json:"-"`
	Title     string
	Video     string `json:"-"`
	SrcType   int
	Date      time.Time `json:"-"`
	ViewCount int       `json:"-"`
	Parts     string    `json:"-"`
	Password  string    `json:"-"`
	Thumbnail string

	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`

	IsURL bool `sql:"-" json:"-"`
}

func GetCountEpisodeByVideoID(db *gorm.DB, videoID string) (count int, err error) {
	err = db.Model(Episode{}).Unscoped().Where("video LIKE ?", "%"+videoID+"%").Count(&count).Error
	return
}
