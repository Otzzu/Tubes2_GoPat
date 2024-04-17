package services

import (
	"container/list"
	"fmt"
	"sync"
)

type SearchState struct {
	Queue   *list.List
	Visited map[string]bool
	Paths   [][]string
	Found   bool
	Depth   int
	Lock    sync.Mutex
}

// func BFS(start string, goal string) ([]string, error) {
// 	queue := list.New()
// 	queue.PushBack([]string{start})

// 	visited := make(map[string]bool)
// 	visited[start] = true

// 	for queue.Len() > 0 {
// 		path := queue.Remove(queue.Front()).([]string)
// 		lastNode := path[len(path)-1]

// 		if lastNode == goal {
// 			return path, nil
// 		}

// 		links, err := ScrapeWikipediaGoQuery(lastNode)
// 		if err != nil {
// 			return nil, err
// 		}

// 		for _, link := range links {
// 			if !visited[link] {
// 				visited[link] = true
// 				newPath := make([]string, len(path))
// 				copy(newPath, path)
// 				newPath = append(newPath, link)
// 				queue.PushBack(newPath)
// 			}
// 		}
// 	}
// 	return nil, fmt.Errorf("path not found")
// }

func BFS(start string, goal string, findAllPaths bool) ([][]string, error) {
	queue := list.New()
	queue.PushBack([]string{start})

	visited := make(map[string]bool)
	visited[start] = true

	var paths [][]string
	foundDepth := -1

	for queue.Len() > 0 {
		path := queue.Remove(queue.Front()).([]string)
		// fmt.Println(path)
		lastNode := path[len(path)-1]
		currentDepth := len(path) - 1

		if lastNode == goal {
			if foundDepth == -1 || currentDepth == foundDepth {
				paths = append(paths, path)
				foundDepth = currentDepth
				if !findAllPaths {
					return paths, nil
				}
				continue
			}
		}

		if foundDepth != -1 && currentDepth > foundDepth {
			continue
		}

		links, err := ScrapeWikipediaGoQuery(lastNode)
		if err != nil {
			return nil, err
		}

		for _, link := range links {
			if !visited[link] {
				visited[link] = true
				newPath := make([]string, len(path))
				copy(newPath, path)
				newPath = append(newPath, link)
				queue.PushBack(newPath)
			}
		}
	}

	if len(paths) > 0 {
		return paths, nil
	}
	return nil, fmt.Errorf("path not found")
}

type Node struct {
	name    string
	isBegin bool
}

type VisitedVal struct {
	depth   int
	isBegin bool
}

func BFS2(start string, goal string) ([]string, error) {
	queue := list.New()
	queue.PushBack([]Node{{start, true}})
	queue.PushBack([]Node{{goal, false}})

	visited := make(map[string]VisitedVal)
	visited[start] = VisitedVal{0, true}
	visited[goal] = VisitedVal{0, false}

	for queue.Len() > 0 {
		path := queue.Remove(queue.Front()).([]Node)
		lastNode := path[len(path)-1]

		fmt.Println(path)
		links, err := ScrapeWikipediaColly(lastNode.name)
		if err != nil {
			return nil, err
		}

		for _, link := range links {
			value, exist := visited[link]
			if exist && value.isBegin != lastNode.isBegin {

				answerPath, err := combine(path, queue, link, lastNode.isBegin, value.depth)
				if err != nil {
					visited[link] = VisitedVal{visited[lastNode.name].depth + 1, lastNode.isBegin}
					newPath := make([]Node, len(path))
					copy(newPath, path)
					newPath = append(newPath, Node{link, lastNode.isBegin})
					queue.PushBack(newPath)
					continue
				}

				return answerPath, nil
			}

			if !exist {
				visited[link] = VisitedVal{visited[lastNode.name].depth + 1, lastNode.isBegin}
				newPath := make([]Node, len(path))
				copy(newPath, path)
				newPath = append(newPath, Node{link, lastNode.isBegin})
				queue.PushBack(newPath)
			}
		}
	}
	return nil, fmt.Errorf("path not found")
}

func combine(arr []Node, source *list.List, value string, isBegin bool, depth int) ([]string, error) {

	arr = append(arr, Node{value, isBegin})

	for el := source.Front(); el != nil; el = el.Next() {
		nodes := el.Value.([]Node)
		lastNode := nodes[len(nodes)-1]


		if lastNode.name == value && lastNode.isBegin != isBegin && len(nodes) == depth+1 {
			if isBegin {

				found1 := true
				for ctr := len(nodes) - 1; ctr >= 1; ctr-- {
					currNode := nodes[ctr]
					links, err := ScrapeWikipediaColly(currNode.name)
					if err != nil {
						found1 = false
						break
					}

					found := false
					for _, link := range links {
						if link == nodes[ctr-1].name {
							found = true
							break
						}
					}

					if !found {
						found1 = false
						break
					}
				}

				if !found1 {
					continue
				}
			} else {

				found1 := true
				for ctr := len(arr) - 1; ctr >= 1; ctr-- {
					currNode := arr[ctr]
					links, err := ScrapeWikipediaColly(currNode.name)
					if err != nil {
						found1 = false
						break
					}

					found := false
					for _, link := range links {
						if link == arr[ctr-1].name {
							found = true
							break
						}
					}

					if !found {
						found1 = false
						break
					}
				}

				if !found1 {
					continue
				}
			}

			var length int = 0
			if isBegin {
				length = len(arr)
			} else {
				length = len(nodes)
			}

			newPath := make([]string, length)

			if isBegin {

				for i, val := range arr {
					newPath[i] = val.name
				}

				for i := len(nodes) - 2; i >= 0; i-- {
					newPath = append(newPath, nodes[i].name)
				}
			} else {
				for i, val := range nodes {
					newPath[i] = val.name
				}

				for i := len(arr) - 2; i >= 0; i-- {
					newPath = append(newPath, arr[i].name)
				}
			}

			return newPath, nil
		}
	}

	return nil, fmt.Errorf("Combiner path failed")
}
