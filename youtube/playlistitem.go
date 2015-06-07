package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (y *Youtube) GetVideoJsonByPlaylistID(playlistID string, botLimit int, pageToken string) (api YoutubePlaylist, err error) {
	apiURL := fmt.Sprintf(YoutubePlaylistItemsAPIURL, y.apiKey, playlistID, botLimit, pageToken)
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
	fmt.Println(string(body))
	err = json.Unmarshal(body, &api)
	if err != nil {
		fmt.Println("### Json Parser Error ", apiURL, " ###")
	}
	return
}

func (y *Youtube) GetVideoByPlaylistID(username string, playlistID string, botLimit int, pageToken string) (totalResults int, youtubeVideos []*YoutubeVideo, prevPageToken string, nextPageToken string) {
	api, err := y.GetVideoJsonByPlaylistID(playlistID, botLimit, pageToken)
	if err != nil {
		panic(err)
	}
	prevPageToken = api.PrevPageToken
	nextPageToken = api.NextPageToken
	if len(api.Items) > 0 {
		for _, item := range api.Items {
			youtubeVideo := &YoutubeVideo{username, "", item.Snippet.ResourceId.VideoId, item.Snippet.Title, item.Snippet.PublishedAt, 0}
			youtubeVideos = append(youtubeVideos, youtubeVideo)
		}
	}
	totalResults = api.PageInfo.TotalResults
	return
}
