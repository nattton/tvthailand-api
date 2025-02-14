package admin

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
}

func ExampleKrobkruakao() {
	krobkruakaos := Krobkruakaos(0)
	for _, kr := range krobkruakaos {
		fmt.Printf("%s - %s\nShort Url : %s, Date : %s\n", kr.Title, kr.Url, kr.ShortUrl, kr.Date)
	}
}

func Krobkruakaos(start int) (krobkruakaos []Krobkruakao) {
	url := fmt.Sprintf("http://www.krobkruakao.com/video_update.php?&starting=%d", start)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".content-150-box .content-150-img a").Each(func(i int, s *goquery.Selection) {
		url, b := s.Attr("href")
		img := s.Find("img").Eq(0)
		title, _ := img.Attr("title")
		if b {
			shortUrl, date := FindGooUrl(url)
			kr := Krobkruakao{title, url, shortUrl, date}
			krobkruakaos = append(krobkruakaos, kr)
		}
	})
	return
}

func FindGooUrl(url string) (shortUrl string, date string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".copy_url span").Each(func(i int, s *goquery.Selection) {
		shortUrl = s.Text()
	})

	doc.Find(".top-space-8").Each(func(i int, s *goquery.Selection) {
		date = strings.Replace(s.Find("p").Eq(1).Text(), "วันที่ ", "", -1)
	})

	return
}
