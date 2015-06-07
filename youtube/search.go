package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (y *Youtube) GetVideoJsonByChannelID(channelID string, botLimit int, pageToken string) (api YoutubeAPI, err error) {
	apiURL := fmt.Sprintf(YoutubeSearchAPIURL, y.apiKey, channelID, botLimit, pageToken)
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
	err = json.Unmarshal(body, &api)
	if err != nil {
		fmt.Println("### Json Parser Error ", apiURL, " ###")
	}
	return
}

func (y *Youtube) GetVideoByChannelID(username string, channelID string, botLimit int, pageToken string) (totalResults int, youtubeVideos []*YoutubeVideo, prevPageToken string, nextPageToken string) {
	api, err := y.GetVideoJsonByChannelID(channelID, botLimit, pageToken)
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
