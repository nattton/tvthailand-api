package api2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const thumbnailUrlCat = "http://thumbnail.instardara.com/category/"
const thumbnailUrlCh = "http://thumbnail.instardara.com/channel/"
const thumbnailUrlRadio = "http://thumbnail.instardara.com/radio/"
const thumbnailUrlTv = "http://thumbnail.instardara.com/tv/"
const thumbnailUrlPoster = "http://thumbnail.instardara.com/poster/"

type Api2Handler struct {
	Db     *sql.DB
	Device string
}

func (h *Api2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.URL)
	w.Header().Set("Content-Type", "application/json")

	h.Device = r.URL.Query().Get("device")
	paths := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	len_path := len(paths)
	if len_path > 1 {
		switch c := paths[1]; c {
		case "advertise":
			return
		case "section":
			h.GetSection_ServeHTTP(w, r)
			return
		case "category":
			if len_path == 2 {
				h.GetCategories_ServeHTTP(w, r)
			} else {
				start := 0
				if len_path > 3 {
					start, _ = strconv.Atoi(paths[3])
				}
				h.GetShowByCategoryID_ServeHTTP(w, r, paths[2], start)
			}

			return
		case "channel":
			if len_path == 2 {
				h.GetChannels_ServeHTTP(w, r)
			} else {
				start := 0
				if len_path > 3 {
					start, _ = strconv.Atoi(paths[3])
				}
				id := paths[2]
				h.GetShowByChannelID_ServeHTTP(w, r, id, start)
			}
			return

		case "radio":
			h.GetRadios_ServeHTTP(w, r)
			return
		case "whatsnew":
			start := 0
			if len_path > 2 {
				start, _ = strconv.Atoi(paths[2])
			}
			h.GetShowByCategoryID_ServeHTTP(w, r, "recents", start)
			return
		case "episode":
			start := 0
			if len_path > 3 {
				start, _ = strconv.Atoi(paths[3])
			}
			id := paths[2]
			h.GetEpisode_ServeHTTP(w, r, id, start)
			return
		}
	}

	http.Error(w, `{"code":404, "message":"Not found"}`, http.StatusNotFound)
}

func (h *Api2Handler) GetSection_ServeHTTP(w http.ResponseWriter, r *http.Request) {
	categories := h.GetAllCategory()
	channels := h.GetAllChannel()
	radios := h.GetAllRadio()
	res := &Section{categories, channels, radios}
	resJ, _ := json.Marshal(res)
	fmt.Fprintf(w, string(resJ))
}

func (h *Api2Handler) GetCategories_ServeHTTP(w http.ResponseWriter, r *http.Request) {
	categories := h.GetAllCategory()
	res := &Categories{categories}
	resJ, _ := json.Marshal(res)
	fmt.Fprintf(w, string(resJ))
}

func (h *Api2Handler) GetChannels_ServeHTTP(w http.ResponseWriter, r *http.Request) {
	channels := h.GetAllChannel()
	res := &Channels{channels}

	resJ, _ := json.Marshal(res)
	fmt.Fprintf(w, string(resJ))
}

func (h *Api2Handler) GetRadios_ServeHTTP(w http.ResponseWriter, r *http.Request) {
	radios := h.GetAllRadio()
	res := &Radios{radios}

	resJ, _ := json.Marshal(res)
	fmt.Fprintf(w, string(resJ))
}

func (h *Api2Handler) GetShowByCategoryID_ServeHTTP(w http.ResponseWriter, r *http.Request, id string, start int) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	shows := h.GetShowByCategoryID(id, start, limit)
	res := &Shows{shows}

	resJ, _ := json.Marshal(res)
	fmt.Fprintf(w, string(resJ))
}

func (h *Api2Handler) GetShowByChannelID_ServeHTTP(w http.ResponseWriter, r *http.Request, id string, start int) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	shows := h.GetShowByChannelID(id, start, limit)
	res := &Shows{shows}

	resJ, _ := json.Marshal(res)
	fmt.Fprintf(w, string(resJ))
}

func (h *Api2Handler) GetEpisode_ServeHTTP(w http.ResponseWriter, r *http.Request, id string, start int) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	showInfo := h.GetShowInfo(id)
	episodes := h.GetEpisode(id, start, limit)
	res := &Episodes{200, showInfo, episodes}
	res.Info = showInfo
	resJ, _ := json.Marshal(res)
	fmt.Fprintf(w, string(resJ))
}
