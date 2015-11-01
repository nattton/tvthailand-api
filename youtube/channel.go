package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/facebookgo/httpcontrol"
)

const YoutubeChannelAPIURL = "https://www.googleapis.com/youtube/v3/channels?key=%s&forUsername=%s&part=snippet"

type ChannelYoutubeAPI struct {
	Items []*ChannelItem `json:"items"`
}

type ChannelItem struct {
	ID string `json:"id"`
}

func (y *Youtube) GetChannelIDByUser(username string) string {
	apiURL := fmt.Sprintf(YoutubeChannelAPIURL, y.apiKey, username)
	fmt.Println(apiURL)
	client := &http.Client{
		Transport: &httpcontrol.Transport{
			RequestTimeout: time.Minute,
			MaxTries:       3,
		},
	}
	resp, err := client.Get(apiURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var api ChannelYoutubeAPI
	err = json.Unmarshal(body, &api)
	if err != nil {
		fmt.Println("### Json Parser Error", apiURL, "###")
	}

	if len(api.Items) > 0 {
		for _, item := range api.Items {
			return item.ID
		}
	}
	return ""
}
