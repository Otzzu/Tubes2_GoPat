package services

import (
	"container/list"
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
    excludedNamespaces = []string{"Category:", "Wikipedia:", "File:", "Help:", "Portal:", "Special:", "Talk:", "User_template:", "Template_talk:", "Mainpage:"}
    wikiArticleRegex   = regexp.MustCompile(`^/wiki/([^#:\s]+)$`)
)

// BiDirectionalBFS finds all shortest paths between start and goal using bidirectional BFS.
func BfsTes(start, goal string) ([][]string, error) {
	if start == goal {
		return [][]string{{start}}, nil
	}

	// Initialize structures for BFS from both directions
	fromStart := make(map[string][][]string)
	fromGoal := make(map[string][][]string)
	visitedFromStart := make(map[string]bool)
	visitedFromGoal := make(map[string]bool)
	queueFromStart := list.New()
	queueFromGoal := list.New()
	queueFromStart.PushBack(start) // Initialize the start queue with the start node
	queueFromGoal.PushBack(goal)

	fromStart[start] = [][]string{{start}}
	fromGoal[goal] = [][]string{{goal}}
	visitedFromStart[start] = true
	visitedFromGoal[goal] = true

	// Search until both queues are empty
	for queueFromStart.Len() > 0 || queueFromGoal.Len() > 0 {
		if queueFromStart.Len() > 0 {
			current := queueFromStart.Remove(queueFromStart.Front()).(string)
			err := bfsStep(current, queueFromStart, visitedFromStart, fromStart)
			if err != nil {
				return nil, err
			}
		}

		if queueFromGoal.Len() > 0 {
			current := queueFromGoal.Remove(queueFromGoal.Front()).(string)
			err := bfsStep(current, queueFromGoal, visitedFromGoal, fromGoal)
			if err != nil {
				return nil, err
			}
		}

		// Check for intersections
		for node := range visitedFromStart {
			if visitedFromGoal[node] {
				// Combine paths through the intersection node
				var combinedPaths [][]string
				fmt.Println("START ", fromStart[node])
				fmt.Println("GOAL ", fromGoal[node])
				for _, path1 := range fromStart[node] {
					for _, path2 := range fromGoal[node] {
						// Reverse the path from the goal to the meeting point to make it go from meeting point to goal
						reversedPath2 := make([]string, len(path2))
						for i, v := range path2 {
							reversedPath2[len(path2)-1-i] = v
						}

						if checkPath(reversedPath2) {

							// Avoid repeating the intersection node and combine the paths
							fullPath := make([]string, len(path1)+len(reversedPath2)-1)
							copy(fullPath, path1)
							copy(fullPath[len(path1):], reversedPath2[1:]) // Skip the meeting point in reversedPath2

							combinedPaths = append(combinedPaths, fullPath)
						}
					}
				}
				return combinedPaths, nil
			}
		}
	}

	return nil, fmt.Errorf("Path not found")
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
func getPaths(pages []string, visited map[string][]string)[][]string {
	var paths [][]string
	for _, page := range pages {
		// if page == "" {
		// 	return make([][]string, 0)
		// } else {
			currentPaths := getPaths(visited[page], visited)
			for _, currentPath := range currentPaths {
				newPath := make([]string, 0, len(currentPath) + 1)
				copy(newPath, currentPath)
				newPath = append(newPath, page)
				fmt.Println(newPath)
				paths = append(paths, newPath)

			}
		// }
	}
	return paths
}

// BreadthFirstSearch performs a bi-directional BFS and returns all shortest paths between the source and target URLs.
func BreadthFirstSearch(sourcePageURL, targetPageURL string) ([][]string, error) {
	if sourcePageURL == targetPageURL {
		return [][]string{{sourcePageURL}}, nil
	}

	visitedForward := make(map[string][]string)
	visitedBackward := make(map[string][]string)
	unvisitedForward := map[string][]string{sourcePageURL: {""}}
	unvisitedBackward := map[string][]string{targetPageURL: {""}}
	var paths [][]string

	for len(paths) == 0 && len(unvisitedForward) != 0 && len(unvisitedBackward) != 0 {
		if len(unvisitedForward) < len(unvisitedBackward){
		// if len(unvisitedForward) > 0 {
			outgoingLinks, err := fetchOutgoingLinks(unvisitedForward)
			// fmt.Println(outgoingLinks)
			if err != nil {
				return nil, fmt.Errorf("error fetching forward links: %v", err)
			}
			for source := range unvisitedForward {
				copy(visitedForward[source], unvisitedForward[source])
			}

			for k := range unvisitedForward {
				delete(unvisitedForward, k)
			}
			

			for sourceURL, targetURLs := range outgoingLinks {
				for _, targetURL := range targetURLs {
					_, found1 := visitedForward[targetURL]
					_, found2:= unvisitedForward[targetURL]

					if !found1 && !found2 {
						unvisitedForward[targetURL] = []string{sourceURL}
					} else if found2 {
						unvisitedForward[targetURL] = append(unvisitedForward[targetURL], sourceURL)
					}
				}
			}

			fmt.Println("foward", unvisitedForward)
			fmt.Println("foward", visitedForward)
		} else {
			outgoingLinks, err := fetchOutgoingLinks(unvisitedBackward)
			// fmt.Println(outgoingLinks)
			if err != nil {
				return nil, fmt.Errorf("error fetching forward links: %v", err)
			}
			for source := range unvisitedBackward {
				copy(visitedForward[source], unvisitedForward[source])

			}

			for k := range unvisitedBackward {
				delete(unvisitedBackward, k)
			}
			

			for sourceURL, targetURLs := range outgoingLinks {
				for _, targetURL := range targetURLs {
					_, found1 := visitedBackward[targetURL]
					_, found2:= unvisitedBackward[targetURL]

					if !found1 && !found2 {
						unvisitedBackward[targetURL] = []string{sourceURL}
					} else if found2 {
						unvisitedBackward[targetURL] = append(unvisitedForward[targetURL], sourceURL)
					}
				}
			}
		}

		// Check for intersection and build paths
		for pageURL := range unvisitedForward {
			if _, ok := unvisitedBackward[pageURL]; ok {
				pathsFromSource := getPaths(unvisitedForward[pageURL], visitedForward)
				pathsFromTarget := getPaths(unvisitedForward[pageURL], visitedBackward)

				fmt.Println(pathsFromSource)
				fmt.Println(pathsFromTarget)

				for _, pathFromSource := range pathsFromSource {
					for _, pathFromTarget := range pathsFromTarget {
						currentPath := append(pathFromSource, append([]string{pageURL}, reverseSlice(pathFromTarget)...)...)
						paths = append(paths, currentPath)
					}
				}
				return paths, nil
			}
		}
	}

	if len(paths) == 0 {
		return nil, errors.New("no path found")
	}
	return paths, nil
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
        MaxIdleConns:        100,
        IdleConnTimeout:     30 * time.Second,
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