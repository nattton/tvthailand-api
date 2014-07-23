package api2

import (
	"database/sql"
	"log"
)

type Radios struct {
	Radios []*Radio `json:"radios"`
}

type Radio struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	URL         string `json:"url"`
	Category    string `json:"category"`
}

func GetRadio(db *sql.DB) []*Radio {
	var radios []*Radio
	rows, err := db.Query("SELECT id, title, description, thumbnail, url, category FROM tv_radio WHERE online = 1 ORDER BY `order`")
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
			category    string
		)
		if err := rows.Scan(&id, &title, &description, &thumbnail, &url, &category); err != nil {
			log.Fatal(err)
		}
		radio := &Radio{id, title, description, thumbnailUrlCat + thumbnail, url, category}
		radios = append(radios, radio)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return radios
}
