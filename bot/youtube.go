package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const YOUTUBE_API_URL = "https://gdata.youtube.com/feeds/api/videos?author=%s&orderby=published&start-index=1&max-results=10&v=2&alt=json"

type Youtube struct {
}

type YoutubeVideo struct {
	Username  string
	Title     string
	VideoId   string
	Published string
}

type YoutubeAPI struct {
	Feed Feed `json:"feed"`
}

type Feed struct {
	Entries []*Entry `json:"entry"`
}

type Entry struct {
	Title      Title      `json:"title"`
	MediaGroup MediaGroup `json:"media$group"`
	Published  Published  `json:"published"`
}

type Title struct {
	Value string `json:"$t"`
}

type MediaGroup struct {
	VideoId VideoId `json:"yt$videoid"`
}
type VideoId struct {
	Value string `json:"$t"`
}

type Published struct {
	Value string `json:"$t"`
}

func (y *Youtube) getVideoByUser(username string) []*YoutubeVideo {
	youtubeVideos := []*YoutubeVideo{}
	apiUrl := fmt.Sprintf(YOUTUBE_API_URL, username)
	resp, err := http.Get(apiUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var api YoutubeAPI
	err = json.Unmarshal(body, &api)
	if err != nil {
		fmt.Println("### Json Parser Error", apiUrl, "###")
	}

	if len(api.Feed.Entries) > 0 {
		for _, entry := range api.Feed.Entries {
			youtubeVideo := &YoutubeVideo{username, entry.Title.Value, entry.MediaGroup.VideoId.Value, entry.Published.Value}
			youtubeVideos = append(youtubeVideos, youtubeVideo)
		}
	}
	return youtubeVideos
}
