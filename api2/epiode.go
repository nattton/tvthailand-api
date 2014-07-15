package api2

import (
	"log"
)

type Episodes struct {
	Code     int        `json:"code"`
	Info     ShowInfo   `json:"info"`
	Episodes []*Episode `json:"episodes"`
}

type Episode struct {
	ID           string `json:"id"`
	EP           string `json:"ep"`
	Title        string `json:"title"`
	VideoEncrypt string `json:"video_encrypt"`
	SrcType      string `json:"src_type"`
	Date         string `json:"date"`
	ViewCount    int    `json:"view_count"`
	Parts        string `json:"parts"`
	Pwd          string `json:"pwd"`
}

func (h *Api2Handler) GetEpisode(id string, start int) []*Episode {
	var episodes []*Episode
	rows, err := h.Db.Query("SELECT programlist_id id, programlist_ep ep, programlist_epname title,  programlist_youtube_encrypt video_encrypt, programlist_src_type src_type, programlist_date date, programlist_count view_count, parts, programlist_password pwd FROM tv_programlist WHERE programlist_banned = 0 AND program_id = ? ORDER BY `ep`, `id` DESC LIMIT ?, ?", id, start, 20)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			id           string
			ep           string
			title        string
			videoEncrypt string
			srcType      string
			date         string
			viewCount    int
			parts        string
			pwd          string
		)
		if err := rows.Scan(&id, &ep, &title, &videoEncrypt, &srcType, &date, &viewCount, &parts, &pwd); err != nil {
			log.Fatal(err)
		}
		episode := &Episode{id, ep, title, videoEncrypt, srcType, date, viewCount, parts, pwd}
		episodes = append(episodes, episode)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return episodes
}
