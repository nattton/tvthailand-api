package api2

import (
	"log"
)

type Advertise struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
}

func (h *Api2Handler) GetAdvertise() []*Advertise {
	var advertises []*Advertise
	rows, err := h.Db.Query("SELECT id, title, description, thumbnail FROM tv_category ORDER BY `order`")
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
		)
		if err := rows.Scan(&id, &title, &description, &thumbnail); err != nil {
			log.Fatal(err)
		}
		advertise := &Advertise{id, title, description, thumbnailUrlCat + thumbnail}
		advertises = append(advertises, advertise)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return advertises
}
