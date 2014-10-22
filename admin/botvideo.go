package admin

import (
	"database/sql"
  "log"
)

type BotVideo struct {
  Db *sql.DB
}

func NewBotVideo(db *sql.DB) *BotVideo {
  b := new(BotVideo)
  b.Db = db
  return b
}

type BotVideoRow struct {
  Id int32
  Username  string
  Title     string
  VideoId   string
  Published string
}

func (b *BotVideo) getBotVideo() []*BotVideoRow {
  botVideos := []*BotVideoRow{}
  rows, err := b.Db.Query("SELECT id, username, title, video_id, published from tv_bot_videos ORDER BY username, published")
  if err != nil {
    panic(err)
  }
  defer rows.Close()
  for rows.Next() {
    var (
      id int32
      username string
      title string
      videoId string
      published string
      )
      if err := rows.Scan(&id, &username, &title, &videoId, &published); err != nil {
        log.Fatal(err)
      }
      botVideo := &BotVideoRow{id, username, title, videoId, published}
      botVideos = append(botVideos, botVideo)
  }

  return botVideos
}
