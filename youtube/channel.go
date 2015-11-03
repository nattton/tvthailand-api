package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/facebookgo/httpcontrol"
)

const YoutubeChannelForUsernameAPIURL = "https://www.googleapis.com/youtube/v3/channels?key=%s&forUsername=%s&part=snippet"
const YoutubeChannelInfoAPIURL = "https://www.googleapis.com/youtube/v3/channels?key=%s&id=%s&part=snippet"

type ChannelYoutubeAPI struct {
	PageInfo PageInfo       `json:"pageInfo"`
	Items    []*ChannelItem `json:"items"`
}

type ChannelItem struct {
	ID string `json:"id"`
}

func (y *Youtube) ChannelIDByUser(username string) string {
	apiURL := fmt.Sprintf(YoutubeChannelForUsernameAPIURL, y.apiKey, username)
	api, _ := y.GetChannelByAPI(apiURL)
	if len(api.Items) > 0 {
		for _, item := range api.Items {
			return item.ID
		}
	}
	return ""
}

func (y *Youtube) ChannelIsActive(channelID string) bool {
	apiURL := fmt.Sprintf(YoutubeChannelInfoAPIURL, y.apiKey, channelID)
	api, err := y.GetChannelByAPI(apiURL)
	if err != nil {
		return true
	}
	return api.PageInfo.TotalResults > 0
}

func (u *Youtube) GetChannelByAPI(apiURL string) (api ChannelYoutubeAPI, err error) {
	fmt.Println(apiURL)
	client := &http.Client{
		Transport: &httpcontrol.Transport{
			RequestTimeout: time.Minute,
			MaxTries:       3,
		},
	}
	resp, err := client.Get(apiURL)
	if err != nil {
		fmt.Println("### Json Parser Error", apiURL, "###", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("### Json Parser Error", apiURL, "###", err)
		return
	}

	err = json.Unmarshal(body, &api)
	if err != nil {
		fmt.Println("### Json Parser Error", apiURL, "###")
		return
	}
	return
}
