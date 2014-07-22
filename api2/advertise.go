package api2

import (
	"log"
)

type Advertises struct {
	Advertises []*Advertise `json:"ads"`
}

type Advertise struct {
	Name     string `json:"name"`
	URL      string `json:"utl"`
	Time     int    `json:"time"`
	Interval int    `json:"interval"`
}

func (h *Api2Handler) GetAdvertise() []*Advertise {
	var advertises []*Advertise
	rows, err := h.Db.Query("SELECT name, url, time FROM tv_advertise WHERE platform = ?", h.Device)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			name string
			url  string
			time int
		)
		if err := rows.Scan(&name, &url, &time); err != nil {
			log.Fatal(err)
		}
		advertise := &Advertise{name, url, time, 10000}
		advertises = append(advertises, advertise)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return advertises
}
