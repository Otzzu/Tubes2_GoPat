package services

import (
	"be/repository"
	"container/list"
	"fmt"
	"net/url"
	"sync/atomic"

	"strings"
	"sync"
	"time"

	"regexp"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

var (
	combinedRegex = regexp.MustCompile(`^/wiki/([^#:\s]+)$`)

	excludedNamespaces2 = []string{
		"Category:", "Wikipedia:", "File:", "Help:", "Portal:",
		"Special:", "Talk:", "User_template:", "Template_talk:", "Mainpage:",
	}
	maxConcurrency = 20
)

func isExcluded(link string) bool {
	if link == "/wiki/Main_Page" {
		return true
	}

	for _, ns := range excludedNamespaces2 {
		if regexp.MustCompile(`^` + regexp.QuoteMeta(ns)).MatchString(link) {
			return true
		}
	}
	return false
}

func ScrapeMultipleWikipediaLinks(urls []string, cache *sync.Map) ([]string, error) {
	// Create a collector
	c := colly.NewCollector(
		colly.Async(true),
	)

	// Create a queue with 2 consumer threads
	q, err := queue.New(
		100, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: len(urls)}, // Use default queue storage
	)
	if err != nil {
		return nil, err
	}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Check if the link matches the combined regex
		if combinedRegex.MatchString(link) {
			fullLink := "https://en.wikipedia.org" + link
			urlKey := e.Request.URL.String()
			// Safely append the link to the slice for the URL
			cache.Store(urlKey, func(value interface{}) interface{} {
				if value == nil {
					return []string{fullLink}
				}
				return append(value.([]string), fullLink)
			})
		}
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Add all URLs to the queue
	for _, url := range urls {
		q.AddURL(url)
	}

	// Consume URLs
	q.Run(c)

	// Wait until threads are finished
	c.Wait()
	return []string{}, nil
}

func DecodePercentEncodedString(encodedString string) string {
	decodedString, err := url.QueryUnescape(encodedString)
	if err != nil {
		return encodedString // return the error if the decoding fails
	}
	return decodedString
}

func ScrapeWikipediaLinks(url string) ([]string, error) {
	if exist, err := repository.GetChildrenByParent(url); exist != nil {
		if err != nil {
			return nil, err
		}
		return exist, nil
	}

	result := make([]string, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("wikipedia.org", "en.wikipedia.org"),
	)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"

	c.SetRequestTimeout(15 * time.Second)

	c.Limit(&colly.LimitRule{
		Parallelism: 1,
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if !isVisible(e) {
			return
		}

		link := e.Attr("href")
		if combinedRegex.MatchString(link) {
			fullLink := "https://en.wikipedia.org" + link
			if !isExcluded(link) {
				result = append(result, fullLink)
			}
		}
	})

	var attempt int
	maxAttempts := 1
	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == 0 { // Network error, possibly a timeout
			attempt++
			if attempt < maxAttempts {
				fmt.Println("Retrying:", r.Request.URL)
				r.Request.Retry()
			} else {
				fmt.Println("Request failed after retries:", r.Request.URL, "\nError:", err)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
	})

	c.Visit(url)

	if len(result) > 0 {

		err := repository.SaveArticleWithChildren(url, result)
		if err != nil {
			fmt.Println(err)
		}
	}

	return result, nil
}

func CompareArrays(arr1, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	countMap := make(map[string]int)

	for _, item := range arr1 {
		countMap[item]++
	}

	for _, item := range arr2 {
		if countMap[item] == 0 {
			return false
		}
		countMap[item]--
	}

	for _, count := range countMap {
		if count != 0 {
			return false
		}
	}

	return true
}

func isVisible(e *colly.HTMLElement) bool {
	class := e.Attr("class")
	class = strings.ReplaceAll(class, " ", "")
	if strings.Contains(class, "nowraplinks") {
		return false
	}

	// Check parent elements for visibility
	for parent := e.DOM.Parent(); parent.Length() != 0; parent = parent.Parent() {

		parentClass, found := parent.Attr("class")
		parentClass = strings.ReplaceAll(parentClass, " ", "")
		if found && strings.Contains(parentClass, "nowraplinks") {
			return false
		}
	}
	return true
}

func helperMulti(urls *list.List, goal string, visited *sync.Map, sem chan struct{}, wg *sync.WaitGroup, count *uint32) [][]string {
	var mu sync.Mutex
	var allPath [][]string

	size := urls.Len()

	for i := 0; i < size; i++ {
		path := urls.Remove(urls.Front()).([]string)
		last := path[len(path)-1]
		sem <- struct{}{}
		wg.Add(1)
		go func(url string, goal string) {
			defer wg.Done()
			defer func() { <-sem }()

			res, _ := ScrapeWikipediaLinks(url)
			for _, u := range res {
				if u == goal {
					newPath := make([]string, len(path))
					copy(newPath, path)
					newPath = append(newPath, u)
					mu.Lock()
					allPath = append(allPath, newPath)
					fmt.Println(allPath)
					mu.Unlock()

				} else {

					if _, exist := visited.LoadOrStore(u, true); !exist {
						newPath := make([]string, len(path))
						copy(newPath, path)
						newPath = append(newPath, u)
						atomic.AddUint32(count, 1)
						mu.Lock()
						urls.PushBack(newPath)
						mu.Unlock()
					}
				}

			}
		}(last, goal)
	}

	wg.Wait()

	if len(allPath) > 0 {
		fmt.Println(len(allPath))
		return allPath
	} else {
		return nil
	}

}

func AsyncBFSMulti(start, goal string) ([][]string, int) {
	var visited sync.Map
	semp := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	var countChecked uint32 = 1

	if start == goal {
		return [][]string{{start}}, 1
	}

	queue := list.New()
	queue.PushBack([]string{start})
	visited.Store(start, true)

	i := 0
	for queue.Len() > 0 {
		i++
		fmt.Println("DEPTH: ", i)
		r := helperMulti(queue, goal, &visited, semp, &wg, &countChecked)
		if r != nil {
			return r, int(countChecked)
		}
	}

	return nil, int(countChecked)
}

func AsyncBFS(start, goal string) ([]string, int) {
	var parent sync.Map
	var visited sync.Map
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var countChecked uint32

	var resultPath []string
	var found bool

	queue := []string{start}
	visited.Store(start, true)

	if start == goal {
		return []string{start}, 1
	}

	for len(queue) > 0 && !found {
		local := queue
		queue = []string{}

		batchsize := (len(local) / maxConcurrency) + 1
		fmt.Println(batchsize)

		for i := 0; i < len(local); i += batchsize {
			end := i + batchsize
			if end > len(local) {
				end = len(local)
			}
			wg.Add(1)
			go func(links []string) {
				defer wg.Done()

				for _, url := range links {
					res, _ := ScrapeWikipediaLinks(url)
					for _, link := range res {
						mutex.Lock()
						if found {
							mutex.Unlock()
							return
						}
						mutex.Unlock()

						if link == goal {
							path := []string{goal}
							for at := url; at != start; {
								path = append([]string{at}, path...)
								tes, _ := parent.Load(at)
								at = tes.(string)
							}
							path = append([]string{start}, path...)

							mutex.Lock()
							if !found {
								resultPath = path
								found = true
								fmt.Println(resultPath)
							}
							mutex.Unlock()
							return
						}

						if _, exist := visited.LoadOrStore(link, true); !exist {
							atomic.AddUint32(&countChecked, 1)

							parent.Store(link, url)
							mutex.Lock()
							if !found {
								queue = append(queue, link)
							} else {
								mutex.Unlock()

								return
							}
							mutex.Unlock()
						}
					}
				}

			}(local[i:end])
		}
		wg.Wait()
	}

	if found {
		return resultPath, int(countChecked)
	}

	return nil, int(countChecked)
}
