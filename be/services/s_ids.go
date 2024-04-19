package services

import (
	"container/list"
	"fmt"
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

func lds(start string, goal string, maxDepth int) ([][]string, error) {
	stack := list.New()
	stack.PushBack([]string{start})

	visited := make(map[string]bool)
	visited[start] = true

	var paths [][]string
	foundDepth := -1

	for stack.Len() > 0 {
		path := stack.Remove(stack.Back()).([]string)
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
			links, err := ScrapeWikipediaQuery(lastNode)
			if err != nil {
				return nil, err
			}

			for _, link := range links {
				if !visited[link] {
					visited[link] = true
					newPath := make([]string, len(path))
					copy(newPath, path)
					newPath = append(newPath, link)
					stack.PushBack(newPath)
				}
			}
		}
	}

	if len(paths) > 0 {
		return paths, nil
	}
	return nil, fmt.Errorf("path not found")
}

// func IDS(start string, goal string, maxDepth int, ) ([]string, error){
//     for i:= 0; i < maxDepth; i++ {
//         path, err := lds(start, goal, i)
//         if err == nil {
//             return path, nil 
//         }
//     }

// 	return nil, fmt.Errorf("path not found in max depth %d", maxDepth)
// }

func IDS(start string, goal string, maxDepth int, ) ([][]string, error){
    for i:= 0; i < maxDepth; i++ {
        path, err := lds(start, goal, i)
        if err == nil {
            return path, nil 
        }
    }

	return nil, fmt.Errorf("path not found in max depth %d", maxDepth)
}