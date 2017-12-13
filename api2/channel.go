package api2

import (
	"database/sql"
	"log"
	"net/http"
)

type Channels []*Channel

type Channel struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	URL         string `json:"url"`
	HasShow     string `json:"has_show"`
}

func GetChannel(db *sql.DB) Channels {
	channels := Channels{}
	rows, err := db.Query("SELECT id, title, description, thumbnail, url, has_show FROM tv_channel ORDER BY `order`")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			id          string
			title       string
			description string
			thumbnail   string
			url         string
			has_show    string
		)
		if err := rows.Scan(&id, &title, &description, &thumbnail, &url, &has_show); err != nil {
			log.Fatal(err)
		}
		ch := &Channel{id, title, description, thumbnailUrlCh + thumbnail, url, has_show}
		channels = append(channels, ch)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return channels
}

func TestLive(db *sql.DB) error {
	stmt := "SELECT id, title, url FROM channels WHERE url != '' AND is_online = true ORDER BY order_display"
	rows, err := db.Query(stmt)
	if err != nil {
		return err
	}

	defer rows.Close()

	channels := Channels{}
	for rows.Next() {
		c := &Channel{}
		err := rows.Scan(&c.ID, &c.Title, &c.URL)
		if err != nil {
			return err
		}
		channels = append(channels, c)
	}

	for _, c := range channels {
		resp, err := http.Get(c.URL)
		if err != nil {
			log.Printf("#### Error Reuqest Channel %s %s / %s\n", c.ID, c.Title, c.URL)
			continue
		}

		log.Printf("Reuqest Channel %s %s / %s Status %s\n", c.ID, c.Title, c.URL, resp.Status)

	}

	return nil
}
