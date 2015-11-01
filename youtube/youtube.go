package youtube

import (
	"os"
)

const YoutubeSearchAPIURL = "https://www.googleapis.com/youtube/v3/search?key=%s&channelId=%s&q=%s&part=snippet&fields=prevPageToken,nextPageToken,pageInfo,items(id(videoId),snippet(title,publishedAt,channelTitle))&order=date&maxResults=%d&pageToken=%s"
const YoutubePlaylistItemsAPIURL = "https://www.googleapis.com/youtube/v3/playlistItems?key=%s&playlistId=%s&part=snippet&fields=prevPageToken,nextPageToken,pageInfo,items(id,snippet(title,publishedAt,channelTitle,resourceId(videoId)))&order=date&maxResults=%d&pageToken=%s"

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
	PrevPageToken string    `json:"prevPageToken"`
	NextPageToken string    `json:"nextPageToken"`
	PageInfo      PageInfo  `json:"pageInfo"`
	Items         []Item    `json:"items"`
	ErrorInfo     ErrorInfo `json:"error"`
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
	Title       string     `json:"title"`
	PublishedAt string     `json:"publishedAt"`
	ResourceID  ResourceID `json:"resourceId"`
}

type YoutubePlaylist struct {
	PrevPageToken string         `json:"prevPageToken"`
	NextPageToken string         `json:"nextPageToken"`
	PageInfo      PageInfo       `json:"pageInfo"`
	Items         []PlaylistItem `json:"items"`
	ErrorInfo     ErrorInfo      `json:"error"`
}

type PlaylistItem struct {
	ID      string  `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type ResourceID struct {
	VideoID string `json:"videoId"`
}
