package main

import (
	"testing"
)

// Youtube
func TestYoutubeVideo(t *testing.T) {
	url := "https://www.youtube.com/watch?v=BQohiu2bMSQ"
	title := "The Voice Thailand - อิมเมจ - ใจนักเลง - 16 Nov 2014"
	videoID := "BQohiu2bMSQ"
	videoType := 0
	date := "2014-11-16"

	v := &WebQuery{}
	videos, err := v.Query(url)
	if err != nil {
		t.Error(err)
	}

	if len(videos) != 1 {
		t.Errorf("videos should length = 1, It length = %d", len(videos))
	} else {
		video := videos[0]
		if video.Title != title {
			t.Errorf("Title should be %s, but is %s", title, video.Title)
		}
		if video.WebURL != url {
			t.Errorf("URL should be %s, but is %s", url, video.WebURL)
		}
		if video.VideoURL != url {
			t.Errorf("VideoURL should be %s, but is %s", url, video.VideoURL)
		}
		if video.VideoID != videoID {
			t.Errorf("VideoID should be %s, but is %s", videoID, video.VideoID)
		}
		if video.VideoType != videoType {
			t.Errorf("VideoType should be %d, but is %d", videoType, video.VideoType)
		}
		if video.Date != date {
			t.Errorf("Date should be %s, but is %s", date, video.Date)
		}
	}
}

// http://www.series8-fc show
func TestSeries8FC(t *testing.T) {
	v := &WebQuery{}
	videos, err := v.Query("http://www.series8-fc.com/“secret-door”")
	if err != nil {
		t.Error(err)
	}

	if len(videos) == 0 {
		t.Error("videos should length > 0")
	}
}

// http://www.kr.series8-fc.com video
// func TestKrSeries8FC(t *testing.T) {
// 	v := &WebQuery{}
// 	videos, err := v.Query("http://www.kr.series8-fc.com/“secret-door-e1”")
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	if len(videos) == 0 {
// 		t.Error("videos should length > 0")
// 	}
// }
