package youtube

import (
	"os"
)

type Youtube struct {
	apiKey string
}

func NewYoutube() *Youtube {
	y := new(Youtube)
	y.apiKey = os.Getenv("YOUTUBE_API_KEY")
	return y
}

type YoutubeVideo struct {
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
	ContentDetails ContentDetails `json:"contentDetails"`
	Snippet        Snippet        `json:"snippet"`
}

type ContentDetails struct {
	Upload     Upload     `json:"upload"`
	ResourceID ResourceID `json:"resourceId"`
}

type Upload struct {
	VideoID string `json:"videoId"`
}

type Snippet struct {
	Title       string `json:"title"`
	PublishedAt string `json:"publishedAt"`
	Type        string `json:"type"`
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
