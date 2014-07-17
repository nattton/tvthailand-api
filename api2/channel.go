package api2

import (
	"log"
)

type Channels struct {
	Channels []*Channel `json:"channels"`
}

type Channel struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	URL         string `json:"url"`
	HasShow     string `json:"has_show"`
}

func (h *Api2Handler) GetAllChannel() []*Channel {
	var channels []*Channel
	rows, err := h.Db.Query("SELECT id, title, description, thumbnail, url, has_show FROM tv_channel ORDER BY `order`")
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
