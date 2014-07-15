package api2

import (
	"log"
)

type Categories struct {
	Categories []*Category `json:"categories"`
}

type Category struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
}

func (h *Api2Handler) GetAllCategory() []*Category {
	var categories = []*Category{
		&Category{"recents", "รายการล่าสุด", "", thumbnailUrlCat + "00_recently.png"},
		&Category{"tophits", "Top Hits", "", thumbnailUrlCat + "00_cate_tophits.png"},
	}
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
		cat := &Category{id, title, description, thumbnailUrlCat + thumbnail}
		categories = append(categories, cat)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return categories
}
