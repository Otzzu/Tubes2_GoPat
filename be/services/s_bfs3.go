package services

import (
	"container/list"
	"context"
	"fmt"

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
		"Special:", "Talk:", "User_template:", "Template_talk:", "Mainpage:", "Main_Page",
	}
	maxConcurrency = 20
	cache          sync.Map
)

func isExcluded(link string) bool {
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
	// if val, exist := cache.Load(url); exist {
	// 	return val.([]string), nil
	// }

	result := make([]string, 0, 50)

	c := colly.NewCollector(
		colly.AllowedDomains("wikipedia.org", "en.wikipedia.org"),
	)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"

	c.SetRequestTimeout(30 * time.Second)

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
	maxAttempts := 2
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

	// cache.Store(url, result)

	return result, nil
}

func b(urls *list.List, goal string, visited *sync.Map, sem chan struct{}, wg *sync.WaitGroup) [][]string {
	var mu sync.Mutex
	var allPath [][]string

	size := urls.Len()

	fmt.Println("size: ", size)

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

func cd(urls *list.List, goal string, visited *sync.Map, sem chan struct{}, wg *sync.WaitGroup) [][]string {
	var mu sync.Mutex
	var allPath [][]string
	found := false

	size := urls.Len()

	fmt.Println("size: ", size)

	for i := 0; i < size; i++ {
		path := urls.Remove(urls.Front()).([]string)
		last := path[len(path)-1]
		sem <- struct{}{}
		wg.Add(1)
		go func(url string, goal string) {
			defer wg.Done()
			defer func() { <-sem }()
			// fmt.Println(runtime.NumGoroutine())
			if found {
				return
			}

			res, _ := ScrapeWikipediaLinks(url)
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
					found = true
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

func AsyncBFS(start, goal string) [][]string {
	var visited sync.Map
	semp := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	queue := list.New()
	queue.PushBack([]string{start})
	visited.Store(start, true)

	i := 0
	for queue.Len() > 0 {
		i++
		fmt.Println("DEPTH: ", i)
		r := b(queue, goal, &visited, semp, &wg)
		if r != nil {
			return r
		}
	}

	return nil

}

func AsyncBFS2(start, goal string) []string {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var visited sync.Map
	queue := make(chan []string, 20)
	found := make(chan []string)
	var wg sync.WaitGroup

	// Start multiple worker goroutines
	numWorkers := 10 // You can set this according to your CPU cores or specific needs
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range queue {
				fmt.Println(path)
				currentURL := path[len(path)-1]
				if currentURL == goal {
					select {
					case found <- path:
						cancel() // Cancel all workers
						return
					case <-ctx.Done():
						return
					}
				}
				res, err := ScrapeWikipediaLinks(currentURL)
				if err != nil {
					fmt.Println("Error scraping:", err)
					continue
				}
				for _, u := range res {
					if _, loaded := visited.LoadOrStore(u, true); !loaded {
						newPath := append([]string(nil), append(path, u)...)
						select {
						case queue <- newPath:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}()
	}

	// Seed the initial URL
	queue <- []string{start}
	visited.Store(start, true)

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(queue)
	}()

	select {
	case res := <-found:
		return res
	case <-ctx.Done():
		return nil
	}
}

func AsyncBFS4(start, goal string) []string {
	var parent sync.Map
	var visited sync.Map
	var mutex sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, maxConcurrency)

	queue := []string{start}
	visited.Store(start, true)

	result := make(chan []string)

	done := make(chan bool)

	go func() {
		for len(queue) > 0 {

			local := make([]string, len(queue))
			copy(local, queue)

			// fmt.Println(len(queue))
			// fmt.Println(len(local))
			queue = make([]string, 0)
			for _, url := range local {
				// fmt.Println(url)
				sem <- struct{}{}
				wg.Add(1)
				go func(url string) {
					defer wg.Done()
					defer func() { <-sem }()

					select {
					case <-result:
						return
					default:

						res, _ := ScrapeWikipediaLinks(url)

						for _, link := range res {
							if link == goal {
								resultPath := []string{url, goal}
								found := false

								fmt.Println("FOUND")
								for !found {
									before, _ := parent.Load(url)
									resultPath = append([]string{before.(string)}, resultPath...)

									if before == start {
										found = true
										result <- resultPath
										done <- true
									}
								}
							} else {
								if _, exist := visited.LoadOrStore(link, true); !exist {
									parent.Store(link, url)
									select {
									case <-result:
										return
									default:
										mutex.Lock()
										queue = append(queue, link)
										mutex.Unlock()
									}

								}
							}
						}
					}
				}(url)
			}
			wg.Wait()

			select {
			case <-result:
				return
			default:
			}

		}

		done <- true

	}()

	select {
	case resultPath := <-result:
		return resultPath
	case <-done:
		return nil
	}

}

func AsyncBFS3(start, goal string) []string {
	var parent sync.Map
	var visited sync.Map
	var wg sync.WaitGroup
	var mutex sync.Mutex

	sem := make(chan struct{}, maxConcurrency)
	queue := []string{start}
	visited.Store(start, true)

	var resultPath []string
	var found bool

	for len(queue) > 0 {
		local := make([]string, len(queue))
		copy(local, queue)
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
			}(url)
		}
		wg.Wait()
	}

	if found {
		return resultPath
	}
	return nil

}

func AsyncBFS5(start, goal string) []string {
	var parent sync.Map
	var visited sync.Map
	var wg sync.WaitGroup
	var mutex sync.Mutex

	var resultPath []string
	var found bool

	queue := []string{start}
	visited.Store(start, true)

	for len(queue) > 0 {
		local := queue
		queue = []string{}

		batchsize := (len(local) / maxConcurrency) + 1
		fmt.Println(batchsize)

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
		return resultPath
	}

	return nil
}

func AsyncBFS6(start, goal string) []string {
	var parent sync.Map
	var visited sync.Map
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// sem := make(chan struct{}, maxConcurrency)
	queue := []string{start}
	visited.Store(start, true)

	var resultPath []string
	var found bool

	c := colly.NewCollector(
		colly.AllowedDomains("wikipedia.org", "en.wikipedia.org"),
	)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"

	c.SetRequestTimeout(30 * time.Second)

	c.Limit(&colly.LimitRule{
		Parallelism: maxConcurrency + 5,
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		url := e.Request.URL.String()
		if combinedRegex.MatchString(href) {
			link := "https://en.wikipedia.org" + href
			if !isExcluded(link) {
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
				}
			}()
		}

		wg.Wait()
		if found {
			break
		}

	}

	if len(resultPath) > 0 {
		return resultPath
	}
	return nil
}
