package api2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dropbox/godropbox/memcache"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
)

const cacheExpireTime = uint32(600)

func setCached(client memcache.ClientShard, key string, value []byte) {
	item := memcache.Item{
		Key:        key,
		Value:      value,
		Flags:      uint32(123),
		Expiration: cacheExpireTime,
	}
	resp := client.Set(&item)
	if resp.Status() != 0 {
		log.Println(resp.Error())
	}
}

func AdvertiseListHandler(db *sql.DB, client memcache.ClientShard, params martini.Params, req *http.Request) (int, string) {
	device := req.URL.Query().Get("device")
	key := fmt.Sprintf("Api2/AdvertiseListHandler/%s", device)

	res := client.Get(key)
	if res.Status() == 0 {
		return http.StatusOK, string(res.Value())
	}

	advertises := GetAdvertise(db, device)
	result := &Advertises{advertises}
	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Oops!!!"
	}

	setCached(client, key, b)
	return http.StatusOK, string(b)
}

func SectionListHandler(db *sql.DB, client memcache.ClientShard, params martini.Params) (int, string) {
	key := "Api2/SectionListHandler"

	res := client.Get(key)
	if res.Status() == 0 {
		return http.StatusOK, string(res.Value())
	}

	categories := GetCategory(db)
	channels := GetChannel(db)
	radios := GetRadio(db)
	result := &Section{categories, channels, radios}
	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Oops!!!"
	}

	setCached(client, key, b)
	return http.StatusOK, string(b)
}

func CategoryListHandler(db *sql.DB, client memcache.ClientShard, params martini.Params) (int, string) {
	key := "Api2/CategoryListHandler"

	res := client.Get(key)
	if res.Status() == 0 {
		return http.StatusOK, string(res.Value())
	}

	categories := GetCategory(db)
	result := &Categories{categories}
	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Oops!!!"
	}

	setCached(client, key, b)
	return http.StatusOK, string(b)
}

func CategoryShowHandler(db *sql.DB, client memcache.ClientShard, params martini.Params) (int, string) {
	id := params["id"]
	start, _ := strconv.Atoi(params["start"])

	key := fmt.Sprintf("Api2/CategoryShowHandler/%s/%d", id, start)

	res := client.Get(key)
	if res.Status() == 0 {
		return http.StatusOK, string(res.Value())
	}

	shows := GetCategoryShow(db, id, start)
	result := &Shows{shows}
	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Oops!!!"
	}

	setCached(client, key, b)
	return http.StatusOK, string(b)
}

func ChannelListHandler(db *sql.DB, client memcache.ClientShard, params martini.Params) (int, string) {
	key := "Api2/ChannelListHandler"

	res := client.Get(key)
	if res.Status() == 0 {
		return http.StatusOK, string(res.Value())
	}

	channels := GetChannel(db)
	result := channels
	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Oops!!!"
	}

	setCached(client, key, b)
	return http.StatusOK, string(b)
}

func ChannelShowHandler(db *sql.DB, client memcache.ClientShard, params martini.Params) (int, string) {
	id := params["id"]
	start, _ := strconv.Atoi(params["start"])

	key := fmt.Sprintf("Api2/ChannelShowHandler/%s/%d", id, start)

	res := client.Get(key)
	if res.Status() == 0 {
		return http.StatusOK, string(res.Value())
	}

	shows := GetChannelShow(db, id, start)
	result := &Shows{shows}
	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Oops!!!"
	}

	setCached(client, key, b)
	return http.StatusOK, string(b)
}

func RadioListHandler(db *sql.DB, client memcache.ClientShard, params martini.Params) (int, string) {
	key := "Api2/RadioListHandler"

	res := client.Get(key)
	if res.Status() == 0 {
		return http.StatusOK, string(res.Value())
	}

	radios := GetRadio(db)
	result := &Radios{radios}
	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Oops!!!"
	}

	setCached(client, key, b)
	return http.StatusOK, string(b)
}

func EpisodeListHandler(db *sql.DB, client memcache.ClientShard, params martini.Params) (int, string) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return http.StatusInternalServerError, "Oops!!!"
	}
	start, _ := strconv.Atoi(params["start"])

	key := fmt.Sprintf("Api2/EpisodeList/%d/%d", id, start)

	res := client.Get(key)
	if res.Status() == 0 {
		return http.StatusOK, string(res.Value())
	}
	showInfo := GetShowInfo(db, id)
	episodes := GetEpisode(db, id, start)
	result := &Episodes{200, showInfo, episodes}

	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Oops!!!"
	}
	setCached(client, key, b)
	return http.StatusOK, string(b)
}
