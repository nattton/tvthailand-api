package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/facebookgo/httpcontrol"
)

func (y *Youtube) GetVideoJsonByChannelID(channelID string, q string, botLimit int, pageToken string) (api YoutubeAPI, err error) {
	apiURL := fmt.Sprintf(YoutubeSearchAPIURL, y.apiKey, channelID, q, botLimit, pageToken)
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
	err = json.Unmarshal(body, &api)
	if err != nil {
		fmt.Println("### Json Parser Error ", apiURL, " ###")
	}
	return
}

func (y *Youtube) GetVideoByChannelID(username string, channelID string, q string, botLimit int, pageToken string) (totalResults int, youtubeVideos []*YoutubeVideo, prevPageToken string, nextPageToken string) {
	api, err := y.GetVideoJsonByChannelID(channelID, q, botLimit, pageToken)
	if err != nil {
		panic(err)
	}
	prevPageToken = api.PrevPageToken
	nextPageToken = api.NextPageToken
	if len(api.Items) > 0 {
		for _, item := range api.Items {
			youtubeVideo := &YoutubeVideo{username, channelID, item.ID.VideoID, item.Snippet.Title, item.Snippet.PublishedAt, 0}
			youtubeVideos = append(youtubeVideos, youtubeVideo)
		}
	}
	totalResults = api.PageInfo.TotalResults
	return
}
