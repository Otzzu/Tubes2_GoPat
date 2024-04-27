package services

import (
	"fmt"
)

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
            }

            for _, link := range links {
                if !visited[link] {
                    visited[link] = true
					(*count)++  // Increment the counter for each scrape
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

	if(start == goal){
		return []string{start}, 1
	}
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