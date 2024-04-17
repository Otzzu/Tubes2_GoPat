package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

var linkCache = make(map[string][]string)

func ScrapeWikipediaGoQuery(url string) ([]string, error) {

	if links, ok := linkCache[url]; ok {
		return links, nil
	}

	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	defer res.Body.Close()

	html, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	html.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists && strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":") && !strings.Contains(link, "Main_Page") {
			fullLink := "https://en.wikipedia.org" + link
			links = append(links, fullLink)
		}
	})
	linkCache[url] = links

	return links, nil
}

func ScrapeWikipediaColly(url string) ([]string, error) {

	if links, ok := linkCache[url]; ok {
		fmt.Println("USE CACHED")
		return links, nil
	}

	c := colly.NewCollector(
		colly.AllowedDomains("wikipedia.org", "en.wikipedia.org"),
	)

	var links []string

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":") && !strings.Contains(link, "Main_Page") {
			fullLink := "https://en.wikipedia.org" + link
			links = append(links, fullLink)
		}
	})

	// c.OnError(func(_ *colly.Response, err error) {
	// 	fmt.Println("Scrape Error:", err)
	// })

	err := c.Visit(url)
	if err != nil {
		fmt.Println("Scrape Error:", err)

		return nil, err
	}

	c.Wait()

	linkCache[url] = links
	return links, nil
}
