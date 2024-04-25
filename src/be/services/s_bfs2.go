package services

import (
	"container/list"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	excludedNamespaces = []string{"Category:", "Wikipedia:", "File:", "Help:", "Portal:", "Special:", "Talk:", "User_template:", "Template_talk:", "Mainpage:", "Main_Page"}
	wikiArticleRegex   = regexp.MustCompile(`^/wiki/([^#:\s]+)$`)
)

func DecodePercentEncodedString(encodedString string) (string) {
    decodedString, err := url.QueryUnescape(encodedString)
    if err != nil {
        return encodedString // return the error if the decoding fails
    }
    return decodedString
}



func checkPath(path []string) bool {

	for i := 0; i < len(path)-1; i++ {
		if !isValidTransition2(path[i], path[i+1]) {
			return false
		}
	}
	return true
}

func isValidTransition2(from, to string) bool {
	links, err := ScrapeWikipediaQuery(from)
	if err != nil {
		return false
	}

	for _, link := range links {
		if link == to {
			return true
		}
	}
	return false
}

// bfsStep helps to perform a single BFS step, updating the queue, visited, and paths.
func bfsStep(current string, queue *list.List, visited map[string]bool, paths map[string][][]string) error {
	neighbors, err := ScrapeWikipediaQuery(current)
	if err != nil {
		return err
	}

	for _, neighbor := range neighbors {
		if !(visited)[neighbor] {
			(visited)[neighbor] = true
			queue.PushBack(neighbor)
			// Create new paths to the neighbor by extending each current path
			var newPaths [][]string
			for _, path := range paths[current] {
				newPath := make([]string, len(path)+1)
				copy(newPath, path)
				newPath[len(newPath)-1] = neighbor
				newPaths = append(newPaths, newPath)
			}
			paths[neighbor] = newPaths
		}
	}

	return nil
}

// getPaths constructs paths from the source or target to the current page.
func getPaths(pages []string, visited map[string][]string) [][]string {
	var paths [][]string
	for _, page := range pages {
		// if page == "" {
		// 	return make([][]string, 0)
		// } else {
		currentPaths := getPaths(visited[page], visited)
		for _, currentPath := range currentPaths {
			newPath := make([]string, 0, len(currentPath)+1)
			copy(newPath, currentPath)
			newPath = append(newPath, page)
			fmt.Println(newPath)
			paths = append(paths, newPath)

		}
		// }
	}
	return paths
}



// fetchOutgoingLinks fetches outgoing links from given Wikipedia page URLs.
// It extracts links from specific elements in parallel.
func fetchOutgoingLinks(pageURLs map[string][]string) (map[string][]string, error) {
	result := make(map[string][]string)
	mutex := &sync.Mutex{} // Safe update of the result map
	errors := make(chan error, len(pageURLs))
	client := customHTTPClient() // Use the optimized HTTP client

	var wg sync.WaitGroup
	const batchSize = 10 // Set a batch size

	pageURLsSlice := mapToSlice(pageURLs) // Convert map keys to slice for batching
	for i := 0; i < len(pageURLsSlice); i += batchSize {
		wg.Add(1)
		batch := pageURLsSlice[i:min(i+batchSize, len(pageURLsSlice))]
		go func(batch []string) {
			defer wg.Done()
			for _, url := range batch {
				links, err := scrapeOutgoingLinks(client, url)
				if err != nil {
					errors <- err
					return
				}
				mutex.Lock()
				result[url] = links
				mutex.Unlock()
			}
		}(batch)
	}

	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		return nil, <-errors
	}

	return result, nil
}

func mapToSlice(m map[string][]string) []string {
	var slice []string
	for key := range m {
		slice = append(slice, key)
	}
	return slice
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// scrapeOutgoingLinks scrapes a given Wikipedia page for outgoing links.
func scrapeOutgoingLinks(client *http.Client, url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("failed to fetch the page")
	}

	// Parse the page with goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	// Find the link tags and extract the URLs
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && wikiArticleRegex.MatchString(href) {
			exclude := false
			for _, namespace := range excludedNamespaces {
				if strings.Contains(href, namespace) {
					exclude = true
					break
				}
			}
			if !exclude {
				fullLink := fmt.Sprintf("https://en.wikipedia.org%s", href)
				links = append(links, fullLink)
			}
		}
	})

	return links, nil
}

func customHTTPClient() *http.Client {
	netTransport := &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 30 * time.Second,
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &http.Client{
		Timeout:   time.Second * 20,
		Transport: netTransport,
	}
}
