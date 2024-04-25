package services

import (
	"container/list"
	"fmt"
	"sync/atomic"

	// "runtime"
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
	maxConcurrency = 40
	cache          sync.Map
)

func isExcluded(link string) bool {
	if link == "https://en.wikipedia.org/wiki/Main_Page" {
		return false
	}

	for _, ns := range excludedNamespaces2 {
		if regexp.MustCompile(`^` + regexp.QuoteMeta(ns)).MatchString(link) {
			return true
		}
	}
	return false
}

func ScrapeMultipleWikipediaLinks(urls []string, cache *sync.Map) ([]string, error) {
	// results := &sync.Map{}

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

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL.String())
	// 	results.Store(r.URL.String(), []string{})
	// })

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

func ScrapeWikipediaLinks(url string) ([]string, error) {
	if val, exist := cache.Load(url); exist {
		if len(val.([]string)) > 0 {

			return val.([]string), nil
		}
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
		link := e.Attr("href")
		// fmt.Println("tes")
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
		// fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println(r.URL)
	})

	c.Visit(url)

	if (len(result) > 0){

		cache.Store(url, result)
	}

	return result, nil
}

func b(urls *list.List, goal string, visited *sync.Map, sem chan struct{}, wg *sync.WaitGroup, count *uint32) [][]string {
	var mu sync.Mutex
	var allPath [][]string

	size := urls.Len()

	// fmt.Println("size: ", size)

	for i := 0; i < size; i++ {
		path := urls.Remove(urls.Front()).([]string)
		last := path[len(path)-1]
		sem <- struct{}{}
		wg.Add(1)
		go func(url string, goal string) {
			defer wg.Done()
			defer func() { <-sem }()
			// fmt.Println(runtime.NumGoroutine())

			res, _ := ScrapeWikipediaLinks(url)
			atomic.AddUint32(count, 1)
			for _, u := range res {
				// fmt.Println(u)
				// fmt.Println(goal, "goal")
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


	queue := list.New()
	queue.PushBack([]string{start})
	visited.Store(start, true)

	i := 0
	for queue.Len() > 0 {
		i++
		fmt.Println("DEPTH: ", i)
		r := b(queue, goal, &visited, semp, &wg, &countChecked)
		if r != nil {
			return r, int(countChecked)
		}
	}

	return nil, int(countChecked)

}

func reconstructPath(start, goal string, parent *sync.Map) []string {
    var path []string
    for at := goal; at != start; {
        path = append([]string{at}, path...)
        p, _ := parent.Load(at)
        at = p.(string)
    }
	path = append([]string{start}, path...)

    return path
}

func AsyncBFS7(start, goal string) ([]string, int) {
    visited := &sync.Map{}
    parent := &sync.Map{}
    queue := []string{start}
    visited.Store(start, true)

    var wg sync.WaitGroup
    var mutex sync.Mutex
    var found uint32
    var count int32

    for len(queue) > 0 && atomic.LoadUint32(&found) == 0 {
        currentBatch := make([]string, len(queue))
        copy(currentBatch, queue)
        queue = []string{}

        for _, url := range currentBatch {
            wg.Add(1)
            go func(url string) {
                defer wg.Done()
                if atomic.LoadUint32(&found) == 1 {
                    return
                }
                links, err := ScrapeWikipediaLinks(url)
                if err != nil {
                    fmt.Println("Error scraping:", err)
                    return
                }
                atomic.AddInt32(&count, 1)

                for _, link := range links {
                    if link == goal {
                        atomic.StoreUint32(&found, 1)
                        parent.Store(goal, url)
                        return
                    }

                    if _, loaded := visited.LoadOrStore(link, true); !loaded {
                        parent.Store(link, url)
                        mutex.Lock()
                        queue = append(queue, link)
                        mutex.Unlock()
                    }
                }
            }(url)
        }
        wg.Wait()
    }

    if atomic.LoadUint32(&found) == 1 {
        return reconstructPath(start, goal, parent), int(count)
    }
    return nil, int(count)
}

func AsyncBFS3(start, goal string) ([]string, int) {
	var parent sync.Map
	var visited sync.Map
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var count uint32

	sem := make(chan struct{}, maxConcurrency)
	queue := []string{start}
	visited.Store(start, true)

	var resultPath []string
	var found bool

	for len(queue) > 0 {
		local := queue
		queue = []string{}

		for _, url := range local {
			sem <- struct{}{}
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				defer func() { <-sem }()

				if found {
					return
				}

				res, _ := ScrapeWikipediaLinks(url)
				atomic.AddUint32(&count, 1)
				for _, link := range res {
					if found {
						return
					}

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
						parent.Store(link, url)
						if !found {
							mutex.Lock()
							queue = append(queue, link)
							mutex.Unlock()
						} else {
							return
						}
					}
				}
			}(url)
		}
		wg.Wait()
	}

	if found {
		return resultPath, int(count)
	}
	return nil, int(count)

}

func AsyncBFS5(start, goal string) ([]string, int) {
	var parent sync.Map
	var visited sync.Map
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var countChecked uint32

	var resultPath []string
	var found bool

	queue := []string{start}
	visited.Store(start, true)

	for len(queue) > 0 {
		local := queue
		queue = []string{}

		batchsize := (len(local) / maxConcurrency) + 1
		// fmt.Println(batchsize)

		for i := 0; i < len(local); i += batchsize {
			end := i + batchsize
			if end > len(local) {
				end = len(local)
			}
			fmt.Println(i)
			wg.Add(1)
			go func(links []string) {
				defer wg.Done()

				for _, url := range links {
					res, _ := ScrapeWikipediaLinks(url)
					atomic.AddUint32(&countChecked, 1)
					for _, link := range res {
						if found {
							return
						}

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
							parent.Store(link, url)
							mutex.Lock()
							if !found {
								queue = append(queue, link)
							} else {
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

func AsyncBFS6(start, goal string) ([]string, int) {
	var parent sync.Map
	var visited sync.Map
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var countChecked uint32 = 1

	// sem := make(chan struct{}, maxConcurrency)
	queue := []string{start}
	visited.Store(start, true)


	var resultPath []string
	var found bool

	c := colly.NewCollector(
		colly.AllowedDomains("wikipedia.org", "en.wikipedia.org"),
	)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"

	c.SetRequestTimeout(15 * time.Second)

	c.Limit(&colly.LimitRule{
		Parallelism: 1,
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if found {
			
			return
		}
		href := e.Attr("href")
		url := e.Request.URL.String()
		if combinedRegex.MatchString(href) {
			link := "https://en.wikipedia.org" + href
			if !isExcluded(link) {
				

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
	})

	for len(queue) > 0 {
		queueUse := queue
		queue = make([]string, 0)
		fmt.Println(len(queueUse))
		for j := 0; j < maxConcurrency; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					if found {
						return
					}

					mutex.Lock()
					if len(queueUse) <= 0 {
						mutex.Unlock()

						return
					}
					url := queueUse[0]
					queueUse = queueUse[1:]

					mutex.Unlock()
					c.Visit(url)
					atomic.AddUint32(&countChecked, 1)
				}
			}()
		}

		wg.Wait()
		if found {
			break
		}

	}

	if len(resultPath) > 0 {
		return resultPath, int(countChecked)
	}
	return nil, int(countChecked)
}
