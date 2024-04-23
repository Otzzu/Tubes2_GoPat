package services

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	// "sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

var LinkCache = make(map[string][]string)

func ScrapeWikipediaQuery(url string) ([]string, error) {

	if links, ok := LinkCache[url]; ok {
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
	html.Find("a[href^='/wiki/']:not([href*=':']):not([href*='Main_Page'])").Each(func(i int, s *goquery.Selection) {
		if link, exists := s.Attr("href"); exists {
			fullLink := "https://en.wikipedia.org" + link
			links = append(links, fullLink)
		}
	})

	LinkCache[url] = links

	return links, nil
}

func ScrapeWikipediaColly(url string) ([]string, error) {

	// if links, ok := LinkCache[url]; ok {
	// 	// fmt.Println("USE CACHED")
	// 	return links, nil
	// }

	c := colly.NewCollector(
		colly.AllowedDomains("wikipedia.org", "en.wikipedia.org"),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{DomainGlob: "*.wikipedia.org", Parallelism: 10})

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"

	var links []string

	var validWikiLink = regexp.MustCompile(`^/wiki/[^:]*$`)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if validWikiLink.MatchString(link) && !strings.Contains(link, "Main_Page") {
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

	// LinkCache[url] = links
	return links, nil
}

func LinksToMap(links []string) map[string]bool {
	linkMap := make(map[string]bool)
	for _, link := range links {
		linkMap[link] = true
	}
	return linkMap
}
