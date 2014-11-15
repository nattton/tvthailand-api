package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const YOUTUBE_API_URL = "https://gdata.youtube.com/feeds/api/videos?author=%s&orderby=published&start-index=1&max-results=%d&v=2&alt=json"

type Youtube struct {
}

type YoutubeVideo struct {
	Username  string
	Title     string
	VideoID   string
	Published string
	Status    int
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
	VideoID VideoID `json:"yt$videoid"`
}
type VideoID struct {
	Value string `json:"$t"`
}

type Published struct {
	Value string `json:"$t"`
}

func (y *Youtube) getVideoByUser(username string, botLimit int) []*YoutubeVideo {
	youtubeVideos := []*YoutubeVideo{}
	apiURL := fmt.Sprintf(YOUTUBE_API_URL, username, botLimit)
	resp, err := http.Get(apiURL)
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
		fmt.Println("### Json Parser Error", apiURL, "###")
	}

	if len(api.Feed.Entries) > 0 {
		for _, entry := range api.Feed.Entries {
			youtubeVideo := &YoutubeVideo{username, entry.Title.Value, entry.MediaGroup.VideoID.Value, entry.Published.Value, 0}
			youtubeVideos = append(youtubeVideos, youtubeVideo)
		}
	}
	return youtubeVideos
}
