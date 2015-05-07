package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const ChannelYoutubeAPIURL = "https://www.googleapis.com/youtube/v3/channels?key=%s&forUsername=%s&part=snippet"

type ChannelYoutubeAPI struct {
	Items []*ChannelItem `json:"items"`
}

type ChannelItem struct {
	ID string `json:"id"`
}

func (y *Youtube) GetChannelIDByUser(username string) string {
	apiURL := fmt.Sprintf(ChannelYoutubeAPIURL, y.apiKey, username)
	fmt.Println(apiURL)
	resp, err := http.Get(apiURL)
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
