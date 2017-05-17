package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Krobkruakao struct {
	Title    string
	Url      string
	ShortUrl string
	Date     string
	Status   int
}

func Krobkruakaos(start int) (krobkruakaos []*Krobkruakao) {
	url := fmt.Sprintf("http://www.krobkruakao.com/video_update.php?&starting=%d", start)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".content-150-box .content-150-img a").Each(func(i int, s *goquery.Selection) {
		url, b := s.Attr("href")
		if b {
			title, shortUrl, date := FindGooUrl(url)
			kr := &Krobkruakao{title, url, shortUrl, date, 0}
			if strings.Contains(shortUrl, "goo.gl") {
				krobkruakaos = append(krobkruakaos, kr)
			}
		}
	})
	return
}

func FindGooUrl(url string) (title string, shortUrl string, date string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
    
    doc.Find(".top-space-8 .Drak12normal").Each(func(i int, s *goquery.Selection) {
		title = s.Text()
	})
    
    doc.Find(".top-space-8 .tag-h").Each(func(i int, s *goquery.Selection) {
        title = fmt.Sprintf("%s | %s", title, s.Text())
	})

	doc.Find(".copy_url span").Each(func(i int, s *goquery.Selection) {
		shortUrl = s.Text()
	})

	doc.Find(".top-space-8").Each(func(i int, s *goquery.Selection) {
		date = strings.Replace(s.Find("p").Eq(1).Text(), "วันที่ ", "", -1)
	})

	return
}
