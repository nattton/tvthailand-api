package api2

import (
	"database/sql"
	"log"
)

const showLimit = 20

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
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	Poster      string `json:"poster"`
	Detail      string `json:"detail"`
	LastEpname  string `json:"last_epname"`
	ViewCount   int    `json:"view_count"`
	Rating      string `json:"rating"`
	VoteCount   int    `json:"vote_count"`
	IsOtv       string `json:"is_otv"`
	OtvID       string `json:"otv_id"`
	OtvApiName  string `json:"otv_api_name"`
}

func GetCategoryShow(db *sql.DB, id string, start int) []*Show {
	var shows []*Show
	var rows *sql.Rows
	var err error
	var sqlQuery = "SELECT program_id id, program_title title, program_time description, program_thumbnail thumbnail, rating, is_otv,     otv_id, otv_api_name FROM tv_program "
	switch id {
	case "recents":
		rows, err = db.Query(sqlQuery+"ORDER BY update_date DESC LIMIT ?, ?", start, showLimit)
	case "tophits":
		rows, err = db.Query(sqlQuery+"ORDER BY view_count DESC LIMIT ?, ?", start, showLimit)
	default:
		rows, err = db.Query(sqlQuery+"WHERE category_id = ? ORDER BY update_date DESC LIMIT ?, ?", id, start, showLimit)
	}
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			id           string
			title        string
			description  string
			thumbnail    string
			rating       string
			is_otv       string
			otv_id       string
			otv_api_name string
		)
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

func GetChannelShow(db *sql.DB, id string, start int) []*Show {
	var shows []*Show
	rows, err := db.Query("SELECT program_id id, program_title title, program_time description, program_thumbnail thumbnail, rating, is_otv, otv_id, otv_api_name FROM tv_program WHERE channel_id = ? ORDER BY update_date DESC LIMIT ?, ?", id, start, showLimit)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			id           string
			title        string
			description  string
			thumbnail    string
			rating       string
			is_otv       string
			otv_id       string
			otv_api_name string
		)
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

func GetShowInfo(db *sql.DB, showId int) ShowInfo {
	var (
		id          string
		title       string
		description string
		thumbnail   string
		poster      string
		detail      string
		lastEpname  string
		viewCount   int
		rating      string
		voteCount   int
		isOtv       string
		otvId       string
		otvApiName  string
	)
	err := db.QueryRow("SELECT program_id id, program_title title, program_thumbnail thumbnail, poster, program_time description, program_detail detail, last_epname, view_count, rating, 5000 as vote_count, is_otv, otv_id, otv_api_name FROM tv_program WHERE program_id = ?", showId).Scan(&id, &title, &description, &thumbnail, &poster, &detail, &lastEpname, &viewCount, &rating, &voteCount, &isOtv, &otvId, &otvApiName)
	if err != nil {
		log.Fatal(err)
	}

	showInfo := ShowInfo{id, title, description, thumbnailUrlTv + thumbnail, poster, detail, lastEpname, viewCount, rating, voteCount, isOtv, otvId, otvApiName}

	return showInfo
}
