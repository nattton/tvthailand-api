package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/facebookgo/httpcontrol"
)

const YoutubePlaylistItemsAPIURL = "https://www.googleapis.com/youtube/v3/playlistItems?key=%s&playlistId=%s&part=snippet&fields=prevPageToken,nextPageToken,pageInfo,items(id,snippet(title,publishedAt,channelTitle,resourceId(videoId)))&maxResults=%d&pageToken=%s"

func (y *Youtube) GetVideoJsonByPlaylistID(playlistID string, botLimit int, pageToken string) (api YoutubePlaylist, err error) {
	apiURL := fmt.Sprintf(YoutubePlaylistItemsAPIURL, y.apiKey, playlistID, botLimit, pageToken)
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

func (y *Youtube) GetVideoByPlaylistID(channelId string, playlistID string, botLimit int, pageToken string) (totalResults int, youtubeVideos []*YoutubeVideo, prevPageToken string, nextPageToken string) {
	api, err := y.GetVideoJsonByPlaylistID(playlistID, botLimit, pageToken)
	if err != nil {
		panic(err)
	}
	prevPageToken = api.PrevPageToken
	nextPageToken = api.NextPageToken
	if len(api.Items) > 0 {
		for _, item := range api.Items {
			youtubeVideo := &YoutubeVideo{channelId, item.Snippet.ResourceID.VideoID, item.Snippet.Title, item.Snippet.PublishedAt, 0}
			youtubeVideos = append(youtubeVideos, youtubeVideo)
		}
	}
	totalResults = api.PageInfo.TotalResults
	return
}
