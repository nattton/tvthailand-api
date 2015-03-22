package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/code-mobi/tvthailand-api/utils"
)

const YoutubeAPIURL = "https://gdata.youtube.com/feeds/api/videos?author=%s&orderby=published&start-index=%d&max-results=%d&v=2&alt=json&random=%s"

type Youtube struct {
}

func NewYoutube() *Youtube {
	y := new(Youtube)
	return y
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
	Entries     []*Entry     `json:"entry"`
	TotalResult TotalResults `json:"openSearch$totalResults"`
}

type Entry struct {
	Title      Title      `json:"title"`
	MediaGroup MediaGroup `json:"media$group"`
	Published  Published  `json:"published"`
}

type TotalResults struct {
	Value int `json:"$t"`
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

func (y *Youtube) GetVideoByUser(username string, start int, botLimit int) (totalResults int, youtubeVideos []*YoutubeVideo) {
	youtubeVideos = []*YoutubeVideo{}
	apiURL := fmt.Sprintf(YoutubeAPIURL, username, start, botLimit, utils.GetTimeStamp())
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
	totalResults = api.Feed.TotalResult.Value
	return
}
