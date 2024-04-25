package services

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly/v2"
	// "sync"
)

// func lds(start string, goal string, maxDepth int) ([]string, error) {
// 	stack := list.New()
// 	stack.PushBack([]string{start})

// 	visited := make(map[string]bool)
// 	visited[start] = true

// 	for stack.Len() > 0 {
// 		path := stack.Remove(stack.Back()).([]string)
// 		lastNode := path[len(path)-1]
// 		currentDepth := len(path) - 1

// 		if lastNode == goal {
// 			return path, nil
// 		}

// 		if currentDepth < maxDepth {
// 			links, err := ScrapeWikipediaGoQuery(lastNode)
// 			if err != nil {
// 				return nil, err
// 			}

// 			for _, link := range links {
// 				if !visited[link] {
// 					visited[link] = true
// 					newPath := make([]string, len(path))
// 					copy(newPath, path)
// 					newPath = append(newPath, link)
// 					stack.PushBack(newPath)
// 				}
// 			}
// 		}
// 	}
// 	return nil, fmt.Errorf("path not found")
// }

func ldsMulti(start string, goal string, maxDepth int, cached *map[string][]string, count *uint32) ([][]string, error) {
    var stack [][]string
    stack = append(stack, []string{start}) 

    visited := make(map[string]bool)
    visited[start] = true

    var paths [][]string
    foundDepth := -1

    for len(stack) > 0 {
        n := len(stack) - 1
        path := stack[n]
        stack = stack[:n]

        lastNode := path[len(path)-1]
        currentDepth := len(path) - 1

        if lastNode == goal {
            if foundDepth == -1 || currentDepth == foundDepth {
                paths = append(paths, path)
                foundDepth = currentDepth
                continue
            }
        }

        if foundDepth != -1 && currentDepth >= foundDepth {
            continue
        }

        if currentDepth < maxDepth {
            var links []string
            var err error

            if cachedLinks, ok := (*cached)[lastNode]; ok {
                links = cachedLinks
            } else {
                links, err = ScrapeWikipediaLinks(lastNode)
                if err != nil {
                    return nil, err
                }
                (*cached)[lastNode] = links
                (*count)++  
            }

            for _, link := range links {
                if !visited[link] {
                    visited[link] = true
                    newPath := append([]string(nil), path...) 
                    newPath = append(newPath, link)
                    stack = append(stack, newPath) 
                }
            }
        }
    }

    if len(paths) > 0 {
        return paths, nil
    }
    return nil, fmt.Errorf("path not found")
}
func lds(start string, goal string, maxDepth int, cached *map[string][]string, count *uint32) ([]string) {
    var stack [][]string
    stack = append(stack, []string{start}) // Using slice as stack

    visited := make(map[string]bool)
    visited[start] = true

    for len(stack) > 0 {
        // Pop from the stack
        n := len(stack) - 1
        path := stack[n]
        stack = stack[:n]

        lastNode := path[len(path)-1]
        currentDepth := len(path) - 1

        if lastNode == goal {
            return path// Return the path immediately when the goal is found
        }

        if currentDepth < maxDepth {
            var links []string
            // var err error

            if cachedLinks, ok := (*cached)[lastNode]; ok {
                links = cachedLinks
            } else {
                links, _ = ScrapeWikipediaLinks(lastNode)
                
                (*cached)[lastNode] = links
                (*count)++  // Increment the counter for each scrape
            }

            for _, link := range links {
                if !visited[link] {
                    visited[link] = true
                    newPath := append([]string(nil), path...) // Make a copy of the path
                    newPath = append(newPath, link)
                    stack = append(stack, newPath) // Push to the stack
                }
            }
        }
    }

    return nil
}



func IDS(start string, goal string, maxDepth int) ([]string, int){
    cached := make(map[string][]string)

	var countChecked uint32
	i := 0
	for  {
		if (i >= maxDepth){
			break
		}

        path  := lds(start, goal, i, &cached, &countChecked)
        if path != nil {
            return path, int(countChecked)
        }
		i++

    }

	return nil, int(countChecked)
}

type SearchState struct {
	TargetURL string
	Found     bool
	FoundPath []string
	Lock      sync.Mutex
	LinkCache sync.Map // Cache for links found on each page using sync.Map
}

// visitPageWithResult starts from a URL and searches up to a given depth. It collects paths to the target.
func visitPageWithResult(c *colly.Collector, url string, depth int, currentPath []string, state *SearchState) {
	if depth <= 0 || state.Found {
		return
	}
	
	// Append the current URL to the path
	newPath := append(currentPath, url)

	// Check if this URL is the target
	if url == state.TargetURL {
		state.Lock.Lock()
		if !state.Found { // Ensure that the path is set only once
			state.Found = true
			state.FoundPath = newPath
		}
		state.Lock.Unlock()
		return
	}
	
	fmt.Println("DEPTH LOKAL :", depth)
	// Check the cache first
	if links, foundInCache := state.LinkCache.Load(url); foundInCache {
		// Use cached links
		fmt.Print("CA ")
		for _, link := range links.([]string) {
			if !contains(newPath, link) && !state.Found {
				visitPageWithResult(c.Clone(), link, depth-1, newPath, state)
			}
		}
	} else {
		// If not found in cache, setup the collector to find and process links on the page
		fg := 0
		var discoveredLinks []string
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			fg++
			fmt.Println(fg)
			href := e.Attr("href")
			if combinedRegex.MatchString(href) {
				link := "https://en.wikipedia.org" + href
				if !isExcluded(link) {
					discoveredLinks = append(discoveredLinks, link)
					if !state.Found {
						// fmt.Print("Y ")

						visitPageWithResult(c.Clone(), link, depth-1, newPath, state)
					}
				}
			}

		})

		c.Visit(url)

		// Cache the discovered links
		state.LinkCache.Store(url, discoveredLinks)
	}
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func IDS2(start, goal string) []string {
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	c.Limit(&colly.LimitRule{
		Parallelism: maxConcurrency + 10,
	})

	state := &SearchState{
		TargetURL: goal,
	}

	depth := 1
	for {
		fmt.Println("DEPTH GEDE :", depth)
		visitPageWithResult(c, start, depth, []string{}, state)
		depth++

		if state.Found {
			fmt.Println(state.FoundPath)
			return state.FoundPath
		} else {

		}

		if depth >= 10 {
			fmt.Println("No Path found in depth >= 10")
			break
		}
	}

	return nil

}

