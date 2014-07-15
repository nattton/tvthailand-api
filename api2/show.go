package api2

import (
	"database/sql"
	"log"
)

type Shows struct {
	Shows []*Show `json:"programs"`
}

type Show struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description"`
	Rating      string `json:"rating"`
	IsOtv       string `json:"is_otv"`
	OtvID       string `json:"otv_id"`
	OtvApiName  string `json:"otv_api_name"`
}

type ShowInfo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Thumbnail   string `json:"thumbnail"`
	Poster      string `json:"poster"`
	Description string `json:"description"`
	Detail      string `json:"detail"`
	LastEpname  string `json:"last_epname"`
	ViewCount   int    `json:"view_count"`
	Rating      string `json:"rating"`
	VoteCount   int    `json:"vote_count"`
	IsOtv       string `json:"is_otv"`
	OtvID       string `json:"otv_id"`
	OtvApiName  string `json:"otv_api_name"`
}

func (h *Api2Handler) GetShowByCategoryID(id string, start int, limit int) []*Show {
	var shows []*Show
	var rows *sql.Rows
	var err error
	var sqlQuery = "SELECT program_id id, program_title title, program_time description, program_thumbnail thumbnail, rating, is_otv,     otv_id, otv_api_name FROM tv_program "
	switch id {
	case "recents":
		rows, err = h.Db.Query(sqlQuery+"ORDER BY update_date DESC LIMIT ?, ?", start, limit)
	case "tophits":
		rows, err = h.Db.Query(sqlQuery+"ORDER BY view_count DESC LIMIT ?, ?", start, limit)
	default:
		rows, err = h.Db.Query(sqlQuery+"WHERE category_id = ? ORDER BY update_date DESC LIMIT ?, ?", id, start, limit)
	}
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var id string
		var title string
		var description string
		var thumbnail string
		var rating string
		var is_otv string
		var otv_id string
		var otv_api_name string
		if err := rows.Scan(&id, &title, &description, &thumbnail, &rating, &is_otv, &otv_id, &otv_api_name); err != nil {
			log.Fatal(err)
		}
		show := &Show{id, title, thumbnailUrlTv + thumbnail, description, rating, is_otv, otv_id, otv_api_name}
		shows = append(shows, show)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return shows
}

func (h *Api2Handler) GetShowByChannelID(id string, start int, limit int) []*Show {
	var shows []*Show
	rows, err := h.Db.Query("SELECT program_id id, program_title title, program_time description, program_thumbnail thumbnail, rating, is_otv, otv_id, otv_api_name FROM tv_program WHERE channel_id = ? ORDER BY update_date DESC LIMIT ?, ?", id, start, limit)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var id string
		var title string
		var description string
		var thumbnail string
		var rating string
		var is_otv string
		var otv_id string
		var otv_api_name string
		if err := rows.Scan(&id, &title, &description, &thumbnail, &rating, &is_otv, &otv_id, &otv_api_name); err != nil {
			log.Fatal(err)
		}
		show := &Show{id, title, thumbnailUrlTv + thumbnail, description, rating, is_otv, otv_id, otv_api_name}
		shows = append(shows, show)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return shows
}
