package admin

import (
	"database/sql"
	"log"
	"strings"
)

type ShowList struct {
	Db *sql.DB
}

type ShowListRow struct {
	ShowID    int
	ShowTitle string
	ID        int
	EP        int
	Title     string
	SrcType   string
	VideoID   string
}

func (s *ShowList) getData(showID int) []*ShowListRow {
	var rows *sql.Rows
	var err error
	showListRows := []*ShowListRow{}
	if showID == 0 {
		rows, err = s.Db.Query("SELECT s.id, s.title, ep.id, ep.ep, ep.title, ep.src_type, ep.video FROM episodes ep INNER JOIN shows s ON (pl.program_id = p.program_id) ORDER BY s.id DESC, ep.id ASC LIMIT 0, 200")
	} else {
		rows, err = s.Db.Query("SELECT s.id, s.title, ep.id, ep.ep, ep.title, ep.src_type, ep.video FROM episodes ep INNER JOIN shows s ON (pl.program_id = p.program_id)  WHERE program_id = ? ORDER BY s.id DESC, ep.id ASC LIMIT 0, 200", showID)
	}

	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			showID    int
			showTitle string
			listID    int
			ep        int
			title     string
			srcType   string
			videoID   string
		)
		if err := rows.Scan(&showID, &showTitle, &listID, &ep, &title, &srcType, &videoID); err != nil {
			log.Fatal(err)
		}
		showListRow := &ShowListRow{showID, showTitle, listID, ep, title, srcType, videoID}
		showListRow.VideoID = strings.Replace(showListRow.VideoID, ",", " ", -1)
		showListRows = append(showListRows, showListRow)
	}

	return showListRows
}

func (s *ShowList) getProgram() []*ShowListRow {
	showListRows := []*ShowListRow{}
	rows, err := s.Db.Query("SELECT DISTINCT s.id, s.title FROM episodes ep INNER JOIN shows s ON (ep.id = s.id) WHERE banned = 1 ORDER BY s.id DESC, ep.id ASC LIMIT 0, 200")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			showID    int
			showTitle string
		)
		if err := rows.Scan(&showID, &showTitle); err != nil {
			log.Fatal(err)
		}
		showListRow := &ShowListRow{showID, showTitle, 0, 0, "", "", ""}
		showListRows = append(showListRows, showListRow)
	}

	return showListRows
}
