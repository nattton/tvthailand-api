package bot

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Bot struct {
	Db *sql.DB
}

func NewBot(db *sql.DB) *Bot {
	b := new(Bot)
	b.Db = db
	return b
}

func (b *Bot) CheckKrobkruakao(start int) {
	krobkruakaos := Krobkruakaos(start)
	for _, kr := range krobkruakaos {
		fmt.Printf("%s - %s\nShort Url : %s, Date : %s\n", kr.Title, kr.Url, kr.ShortUrl, kr.Date)
		b.checkBotKrobkruakaoExistingAndAddBot(kr)
	}
}

func (b *Bot) checkBotKrobkruakaoExistingAndAddBot(video *Krobkruakao) {
	rows, err := b.Db.Query("SELECT id from bot_videos WHERE url = ?", video.Url)
	if err != nil {
		fmt.Println(err)
	} else {
		defer rows.Close()

		if !rows.Next() {
			row2s, err := b.Db.Query("SELECT id from episodes WHERE video LIKE ? ", "%"+video.ShortUrl+"%")
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

func (b *Bot) insertBotKrobkruakao(video *Krobkruakao) {
	_, err := b.Db.Exec("INSERT INTO bot_videos (channel_id, title, url, video_id, video_type, published_at, status) VALUES ('UCirZPTc9IoKM_DsA9aKbc4g', ?, ?, ?, 'url', NOW(), ?)", video.Title+" | "+video.Date, video.Url, video.ShortUrl, video.Status)
	if err != nil {
		panic(err)
	}
	fmt.Println("Insert Bot Krobkruakao ### ", video.Title, video.Url, video.ShortUrl, video.Date, "###")
}
