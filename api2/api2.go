package api2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dropbox/godropbox/memcache"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const cacheExpireTime = uint32(600)

const thumbnailUrlCat = "http://thumbnail.instardara.com/category/"
const thumbnailUrlCh = "http://thumbnail.instardara.com/channel/"
const thumbnailUrlRadio = "http://thumbnail.instardara.com/radio/"
const thumbnailUrlTv = "http://thumbnail.instardara.com/tv/"
const thumbnailUrlPoster = "http://thumbnail.instardara.com/poster/"

type Api2Handler struct {
	Db             *sql.DB
	MemcacheClient memcache.ClientShard
	Device         string
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
			h.GetAdvertise_ServeHTTP(w, r)
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

func (h *Api2Handler) setCached(key string, value []byte) {
	item := memcache.Item{
		Key:        key,
		Value:      value,
		Flags:      uint32(123),
		Expiration: cacheExpireTime,
	}
	resp := h.MemcacheClient.Set(&item)
	if resp.Status() != 0 {
		log.Println(resp.Error())
	}
}

func (h *Api2Handler) GetAdvertise_ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := fmt.Sprintf("Api2Handler/Advertise/%s", h.Device)

	res := h.MemcacheClient.Get(key)
	if res.Status() == 0 {
		fmt.Fprintf(w, string(res.Value()))
		return
	}

	advertises := h.GetAdvertise()
	result := &Advertises{advertises}

	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		http.Error(w, "Oops!!!", http.StatusInternalServerError)
	}

	h.setCached(key, b)
	fmt.Fprintf(w, string(b))
}

func (h *Api2Handler) GetSection_ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := "Api2Handler/Section"

	res := h.MemcacheClient.Get(key)
	if res.Status() == 0 {
		fmt.Fprintf(w, string(res.Value()))
		return
	}

	categories := h.GetAllCategory()
	channels := h.GetAllChannel()
	radios := h.GetAllRadio()
	result := &Section{categories, channels, radios}

	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		http.Error(w, "Oops!!!", http.StatusInternalServerError)
	}

	h.setCached(key, b)
	fmt.Fprintf(w, string(b))
}

func (h *Api2Handler) GetCategories_ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := "Api2Handler/Categories"

	res := h.MemcacheClient.Get(key)
	if res.Status() == 0 {
		fmt.Fprintf(w, string(res.Value()))
		return
	}

	categories := h.GetAllCategory()
	result := &Categories{categories}
	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		http.Error(w, "Oops!!!", http.StatusInternalServerError)
	}

	h.setCached(key, b)
	fmt.Fprintf(w, string(b))
}

func (h *Api2Handler) GetChannels_ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := "Api2Handler/Channels"

	res := h.MemcacheClient.Get(key)
	if res.Status() == 0 {
		fmt.Fprintf(w, string(res.Value()))
		return
	}

	channels := h.GetAllChannel()
	result := &Channels{channels}

	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		http.Error(w, "Oops!!!", http.StatusInternalServerError)
	}

	h.setCached(key, b)
	fmt.Fprintf(w, string(b))
}

func (h *Api2Handler) GetRadios_ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := "Api2Handler/Radio"

	res := h.MemcacheClient.Get(key)
	if res.Status() == 0 {
		fmt.Fprintf(w, string(res.Value()))
		return
	}
	radios := h.GetAllRadio()
	result := &Radios{radios}

	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		http.Error(w, "Oops!!!", http.StatusInternalServerError)
	}

	h.setCached(key, b)
	fmt.Fprintf(w, string(b))
}

func (h *Api2Handler) GetShowByCategoryID_ServeHTTP(w http.ResponseWriter, r *http.Request, id string, start int) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	key := fmt.Sprintf("Api2Handler/ShowByCategoryID/%s/%d/%d", id, start, limit)

	res := h.MemcacheClient.Get(key)
	if res.Status() == 0 {
		fmt.Fprintf(w, string(res.Value()))
		return
	}

	shows := h.GetShowByCategoryID(id, start, limit)
	result := &Shows{shows}

	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		http.Error(w, "Oops!!!", http.StatusInternalServerError)
	}

	h.setCached(key, b)
	fmt.Fprintf(w, string(b))
}

func (h *Api2Handler) GetShowByChannelID_ServeHTTP(w http.ResponseWriter, r *http.Request, id string, start int) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	key := fmt.Sprintf("Api2Handler/ShowByChannelID/%s/%d/%d", id, start, limit)

	res := h.MemcacheClient.Get(key)
	if res.Status() == 0 {
		fmt.Fprintf(w, string(res.Value()))
		return
	}

	shows := h.GetShowByChannelID(id, start, limit)
	result := &Shows{shows}

	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		http.Error(w, "Oops!!!", http.StatusInternalServerError)
	}

	h.setCached(key, b)
	fmt.Fprintf(w, string(b))
}

func (h *Api2Handler) GetEpisode_ServeHTTP(w http.ResponseWriter, r *http.Request, id string, start int) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	key := fmt.Sprintf("Api2Handler/Episode/%s/%d/%d", id, start, limit)

	res := h.MemcacheClient.Get(key)
	if res.Status() == 0 {
		fmt.Fprintf(w, string(res.Value()))
		return
	}

	showInfo := h.GetShowInfo(id)
	episodes := h.GetEpisode(id, start, limit)
	result := &Episodes{200, showInfo, episodes}

	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		http.Error(w, "Oops!!!", http.StatusInternalServerError)
	}

	h.setCached(key, b)
	fmt.Fprintf(w, string(b))
}
