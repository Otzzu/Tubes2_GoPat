package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Node struct {
	Name string
	Path []string
}

type Visited map[string]bool
type Cache map[string][]string

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func getLinks(url string, cache Cache, cacheMutex *sync.RWMutex) ([]string, error) {
	cacheMutex.RLock()
	if cachedLinks, ok := cache[url]; ok {
		cacheMutex.RUnlock()
		return cachedLinks, nil
	}
	cacheMutex.RUnlock()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	doc.Find("#bodyContent a[href]").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists && strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":") {
			completeLink := "https://en.wikipedia.org" + link
			if !contains(links, completeLink) {
				links = append(links, completeLink)
			}
		}
	})

	cacheMutex.Lock()
	cache[url] = links
	cacheMutex.Unlock()

	return links, nil
}

func worker(jobs <-chan Node, results chan<- Node, cache Cache, cacheMutex *sync.RWMutex, visited Visited, visitedMutex *sync.Mutex) {
	for node := range jobs {
		links, err := getLinks(node.Name, cache, cacheMutex)
		if err != nil {
			continue // Handle the error appropriately in production
		}

		for _, link := range links {
			visitedMutex.Lock()
			if !visited[link] {
				visited[link] = true
				visitedMutex.Unlock()
				newPath := append([]string{}, node.Path...)
				newPath = append(newPath, link)
				results <- Node{Name: link, Path: newPath}
			} else {
				visitedMutex.Unlock()
			}
		}
	}
}

func BFS(start, goal string, cache Cache, cacheMutex *sync.RWMutex) ([]string, error) {
	visited := make(Visited)
	visitedMutex := new(sync.Mutex)
	jobs := make(chan Node, 100)
	results := make(chan Node, 100)

	numWorkers := 10 // Set the number of workers
	var wg sync.WaitGroup
	var pathFound bool
	var finalPath []string

	// Start workers
	wg.Add(numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go func() {
			worker(jobs, results, cache, cacheMutex, visited, visitedMutex)
			wg.Done()
		}()
	}

	// Send the first job
	jobs <- Node{Name: start, Path: []string{start}}
	visited[start] = true

	go func() {
		for result := range results {
			if result.Name == goal {
				finalPath = result.Path
				pathFound = true
				close(jobs)
				break
			}
			jobs <- result
		}
	}()

	wg.Wait()
	close(results)
	if pathFound {
		return finalPath, nil
	}
	return nil, fmt.Errorf("no path found from %s to %s", start, goal)
}

func main() {
	start := "https://en.wikipedia.org/wiki/Computer_science"
	goal := "https://en.wikipedia.org/wiki/Artificial_intelligence"

	cache := make(Cache)
	cacheMutex := new(sync.RWMutex)

	startTime := time.Now()
	path, err := BFS(start, goal, cache, cacheMutex)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Path found:", path)
	}
	fmt.Println("Time taken:", time.Since(startTime))
}
