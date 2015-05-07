package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const YoutubeAPIURL = "https://www.googleapis.com/youtube/v3/search?key=%s&channelId=%s&part=snippet&fields=pageInfo,items(id(videoId),snippet(title,publishedAt,channelTitle))&order=date&maxResults=%d"

type Youtube struct {
	apiKey string
}

func NewYoutube() *Youtube {
	y := new(Youtube)
	y.apiKey = os.Getenv("YOUTUBE_API_KEY")
	return y
}

type YoutubeVideo struct {
	Username  string
	ChannelID string
	VideoID   string
	Title     string
	Published string
	Status    int
}

type YoutubeAPI struct {
	PageInfo  PageInfo  `json:"pageInfo"`
	Items     []*Item   `json:"items"`
	ErrorInfo ErrorInfo `json:"error"`
}

type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type Item struct {
	ID      ItemID  `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type ItemID struct {
	VideoID string `json:"videoId"`
}

type Snippet struct {
	Title       string `json:"title"`
	PublishedAt string `json:"publishedAt"`
}

func (y *Youtube) GetVideoByChannelID(username string, channelID string, botLimit int) (totalResults int, youtubeVideos []*YoutubeVideo) {
	apiURL := fmt.Sprintf(YoutubeAPIURL, y.apiKey, channelID, botLimit)
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

	var api YoutubeAPI
	err = json.Unmarshal(body, &api)
	if err != nil {
		fmt.Println("### Json Parser Error", apiURL, "###")
	}

	if len(api.Items) > 0 {
		for _, item := range api.Items {
			youtubeVideo := &YoutubeVideo{username, channelID, item.ID.VideoID, item.Snippet.Title, item.Snippet.PublishedAt, 0}
			youtubeVideos = append(youtubeVideos, youtubeVideo)
		}
	}
	totalResults = api.PageInfo.TotalResults
	return
}

// func (y *Youtube) GetVideoByUser(username string, botLimit int) (totalResults int, youtubeVideos []*YoutubeVideo) {
// 	apiURL := fmt.Sprintf(YoutubeAPIURL, username, botLimit)
// 	fmt.Println(apiURL)
// 	resp, err := http.Get(apiURL)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	var api YoutubeAPI
// 	err = json.Unmarshal(body, &api)
// 	if err != nil {
// 		fmt.Println("### Json Parser Error", apiURL, "###")
// 	}
//
// 	if len(api.Items) > 0 {
// 		for _, item := range api.Items {
// 			youtubeVideo := &YoutubeVideo{username, entry.Title.Value, entry.MediaGroup.VideoID.Value, entry.Published.Value, 0}
// 			youtubeVideos = append(youtubeVideos, youtubeVideo)
// 		}
// 	}
// 	totalResults = api.Feed.TotalResult.Value
// 	return
// }

func (y *Youtube) GetVideoByUserAndKeyword(username string, start int, botLimit int, keyword string) (totalResults int, youtubeVideos []*YoutubeVideo) {
	// youtubeVideos = []*YoutubeVideo{}
	// if strings.Contains(keyword, " ") {
	// 	keyword = url.QueryEscape(keyword)
	// }
	// apiURL := fmt.Sprintf(YoutubeAPIURL, username, start, botLimit, keyword)
	// fmt.Println(apiURL)
	// resp, err := http.Get(apiURL)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// var api YoutubeAPI
	// err = json.Unmarshal(body, &api)
	// if err != nil {
	// 	fmt.Println("### Json Parser Error", apiURL, "###")
	// }
	//
	// if len(api.Items) > 0 {
	// 	for _, item := range api.Items {
	// 		youtubeVideo := &YoutubeVideo{username, entry.Title.Value, entry.MediaGroup.VideoID.Value, entry.Published.Value, 0}
	// 		youtubeVideos = append(youtubeVideos, youtubeVideo)
	// 	}
	// }
	// totalResults = api.Feed.TotalResult.Value
	return
}
