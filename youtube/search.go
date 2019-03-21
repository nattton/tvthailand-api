package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/facebookgo/httpcontrol"
)

// const YoutubeSearchAPIURL = "https://www.googleapis.com/youtube/v3/search?key=%s&channelId=%s&q=%s&part=snippet&fields=prevPageToken,nextPageToken,pageInfo,items(id(videoId),snippet(title,publishedAt,channelTitle))&order=date&maxResults=%d&pageToken=%s"
const YoutubeSearchAPIURL = "https://www.googleapis.com/youtube/v3/activities?key=%s&channelId=%s&q=%s&part=snippet,contentDetails&fields=items(contentDetails,kind,snippet(channelId,description,publishedAt,title,type)),nextPageToken,pageInfo,prevPageToken&maxResults=%d&pageToken=%s"

func (y *Youtube) GetVideoJsonByChannelID(channelID string, query string, botLimit int, pageToken string) (api YoutubeAPI, err error) {
	apiURL := fmt.Sprintf(YoutubeSearchAPIURL, y.apiKey, channelID, url.QueryEscape(query), 50, pageToken)
	fmt.Println(apiURL)
	client := &http.Client{
		Transport: &httpcontrol.Transport{
			RequestTimeout: time.Minute,
			MaxTries:       3,
		},
	}
	resp, err := client.Get(apiURL)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal(body, &api)
	if err != nil {
		fmt.Println("### Json Parser Error ", apiURL, " ###")
	}
	return
}

func (y *Youtube) GetVideoByChannelID(channelID string, q string, botLimit int, pageToken string) (totalResults int, youtubeVideos []*YoutubeVideo, prevPageToken string, nextPageToken string) {
	api, err := y.GetVideoJsonByChannelID(channelID, q, botLimit, pageToken)
	if err != nil {
		panic(err)
	}
	prevPageToken = api.PrevPageToken
	nextPageToken = api.NextPageToken
	if len(api.Items) > 0 {
		for _, item := range api.Items {
			if item.Snippet.Type == "upload" {
				youtubeVideo := &YoutubeVideo{channelID, item.ContentDetails.Upload.VideoID, item.Snippet.Title, item.Snippet.PublishedAt, 0}
				youtubeVideos = append(youtubeVideos, youtubeVideo)
			}
		}
	}
	totalResults = api.PageInfo.TotalResults
	return
}
