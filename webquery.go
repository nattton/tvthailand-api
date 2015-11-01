package main

import (
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/PuerkitoBio/goquery"
)

const YoutubeAPIURL = "https://www.googleapis.com/youtube/v3/videos?key=%s&id=%s&fields=items(id,snippet(channelId,title,publishedAt),statistics)&part=snippet"

type WebQuery struct {
}

// Video
type Video struct {
	Index     int
	Title     string
	WebURL    string
	VideoURL  string
	VideoID   string
	VideoType int
	Date      string
}

func (w *WebQuery) Query(webURL string) (videos []*Video, err error) {
	videos = []*Video{}
	doc, err := goquery.NewDocument(webURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	u, err := url.Parse(webURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	v := u.Query()

	if strings.Contains(webURL, "www.youtube.com/watch") {
		videoID := v.Get("v")
		title, exists := doc.Find("#eow-title").Attr("title")
		// publishDate := doc.Find("#watch-uploader-info .watch-time-text").Text()
		if exists {
			video := &Video{
				Index:     0,
				Title:     title,
				WebURL:    webURL,
				VideoURL:  webURL,
				VideoID:   videoID,
				VideoType: 0,
				Date:      w.Today(),
			}
			videos = append(videos, video)
		}
	} else if strings.Contains(webURL, "http://www.series8-fc.com") {
		doc.Find("div.entry ul li a").Each(func(i int, s *goquery.Selection) {
			title := s.Text()
			webURL, exists := s.Attr("href")
			if exists {
				// fmt.Printf("%d : %s - %s\n", i, title, url)
				video := &Video{Index: i, Title: title, WebURL: webURL}
				videos = append(videos, video)
			}
		})
	}

	if strings.Contains(webURL, "http://www.kr.series8-fc.com") {

	}

	return
}

func (w *WebQuery) TitleToDate(title string) string {

	return ""
}

func (w *WebQuery) Today() string {
	t := time.Now().Local()
	return t.Format("2014-01-01")
}
