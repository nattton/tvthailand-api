package admin

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

type Krobkruakao struct {
	Title    string
	Url      string
	ShortUrl string
}

func ExampleKrobkruakao() {
	krobkruakaos := Krobkruakaos()
	for _, kr := range krobkruakaos {
		fmt.Printf("%s - %s\nShort Url : %s\n", kr.Title, kr.Url, kr.ShortUrl)
	}
}

func Krobkruakaos() (krobkruakaos []Krobkruakao) {
	doc, err := goquery.NewDocument("http://www.krobkruakao.com/%E0%B8%A3%E0%B8%B2%E0%B8%A2%E0%B8%81%E0%B8%B2%E0%B8%A3%E0%B8%82%E0%B9%88%E0%B8%B2%E0%B8%A7%E0%B8%A2%E0%B9%89%E0%B8%AD%E0%B8%99%E0%B8%AB%E0%B8%A5%E0%B8%B1%E0%B8%87/")
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".content-150-box .content-150-img a").Each(func(i int, s *goquery.Selection) {
		url, b := s.Attr("href")
		img := s.Find("img").Eq(0)
		title, _ := img.Attr("title")
		if b {
			shortUrl := FindGooUrl(url)
			kr := Krobkruakao{title, url, shortUrl}
			krobkruakaos = append(krobkruakaos, kr)
		}
	})
	return
}

func FindGooUrl(url string) (shortUrl string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".copy_url span").Each(func(i int, s *goquery.Selection) {
		shortUrl = s.Text()
		return
	})
	return
}
