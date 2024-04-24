package services

import (
	"be/models"
	"container/list"
	"fmt"
)


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
				queue.PushBack(newPath)
			}
		}
	}

	if len(paths) > 0 {
		return paths, nil
	}
	return nil, fmt.Errorf("path not found")
}

func BFS2(start string, goal string) ([]string, error) {
	LinkCache = make(map[string][]string)
	queue := list.New()
	queue.PushBack([]models.Node{{Name: start, IsBegin: true}})
	queue.PushBack([]models.Node{{Name: goal, IsBegin: false}})

	visited := make(map[string]models.VisitedVal)
	visited[start] = models.VisitedVal{Depth: 0, IsBegin: true, Path: make([]string, 0, 5)}
	visited[goal] = models.VisitedVal{Depth: 0, IsBegin: false, Path: make([]string, 0, 5)}

	for queue.Len() > 0 {
		path := queue.Remove(queue.Front()).([]models.Node)
		lastNode := path[len(path)-1]

		links, err := ScrapeWikipediaColly(lastNode.Name)
		if err != nil {
			return nil, err
		}

		for _, link := range links {
			value, exist := visited[link]
			if exist && value.IsBegin != lastNode.IsBegin {
				answerPath, err := combine(path, value.Path, link, lastNode.IsBegin)
				if err != nil {
					visited[link] = models.VisitedVal{Depth: visited[lastNode.Name].Depth + 1, IsBegin: lastNode.IsBegin, Path: append(visited[lastNode.Name].Path, lastNode.Name)}
					newPath := make([]models.Node, len(path)+1)
					copy(newPath, path)
					newPath[len(path)] = models.Node{Name: link, IsBegin: lastNode.IsBegin}
					queue.PushBack(newPath)
					continue
				}

				return answerPath, nil
			}

			if !exist {
				visited[link] = models.VisitedVal{Depth: visited[lastNode.Name].Depth + 1, IsBegin: lastNode.IsBegin, Path: append(visited[lastNode.Name].Path, lastNode.Name)}
				newPath := make([]models.Node, len(path)+1)
				copy(newPath, path)
				newPath[len(path)] = models.Node{Name: link, IsBegin: lastNode.IsBegin}
				queue.PushBack(newPath)
			}
		}
	}
	return nil, fmt.Errorf("path not found")
}

func BFS3(start string, goal string) ([]string, error) {
	queueFront := list.New()
	queueBack := list.New()
	queueFront.PushBack([]string{start})
	queueBack.PushBack([]string{goal})

	visitedFront := make(map[string]models.VisitedVal)
	visitedBack := make(map[string]models.VisitedVal)
	visitedFront[start] = models.VisitedVal{Depth: 0, Path: make([]string, 0, 5)}
	visitedBack[goal] = models.VisitedVal{Depth: 0, Path: make([]string, 0, 5)}

	for queueFront.Len() > 0 || queueBack.Len() > 0 {

		if queueFront.Len() > 0 {

			path := queueFront.Remove(queueFront.Front()).([]string)
			lastNode := path[len(path)-1]

			links, err := ScrapeWikipediaColly(lastNode)
			if err != nil {
				return nil, err
			}

			for _, link := range links {
				value, exist := visitedBack[link]
				if exist {
					answerPath, err := combine2(path, value.Path, link)
					if err != nil {
						if _, exist2 := visitedFront[link]; !exist2 {

							visitedFront[link] = models.VisitedVal{Depth: visitedFront[lastNode].Depth + 1, Path: append(visitedFront[lastNode].Path, lastNode)}
							newPath := make([]string, len(path)+1)
							copy(newPath, path)
							newPath[len(path)] = link
							queueFront.PushBack(newPath)
						}
						delete(visitedBack, link)
						continue
					}

					return answerPath, nil
				}

				if !exist {
					if _, exist2 := visitedFront[link]; !exist2 {

						visitedFront[link] = models.VisitedVal{Depth: visitedFront[lastNode].Depth + 1, Path: append(visitedFront[lastNode].Path, lastNode)}
						newPath := make([]string, len(path)+1)
						copy(newPath, path)
						newPath[len(path)] = link
						queueFront.PushBack(newPath)
					}
				}
			}
		}

		if queueBack.Len() > 0 {

			path := queueBack.Remove(queueBack.Front()).([]string)
			lastNode := path[len(path)-1]

			links, err := ScrapeWikipediaColly(lastNode)
			if err != nil {
				return nil, err
			}

			for _, link := range links {
				value, exist := visitedFront[link]
				if exist {
					answerPath, err := combine2(value.Path, path, link)
					if err == nil {

						return answerPath, nil
					}
				}

				if !exist {
					if _, exist2 := visitedBack[link]; !exist2 {

						visitedBack[link] = models.VisitedVal{Depth: visitedBack[lastNode].Depth + 1, Path: append(visitedBack[lastNode].Path, lastNode)}
						newPath := make([]string, len(path)+1)
						copy(newPath, path)
						newPath[len(path)] = link
						queueBack.PushBack(newPath)
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("path not found")
}

func BFS4(start string, goal string) ([]string, error) {
	queueFront := list.New()
	queueBack := list.New()
	queueFront.PushBack([]string{start})
	queueBack.PushBack([]string{goal})

	visitedFront := make(map[string]models.VisitedVal)
	visitedBack := make(map[string]models.VisitedVal)
	visitedFront[start] = models.VisitedVal{Depth: 0, Path: make([]string, 0, 5)}
	visitedBack[goal] = models.VisitedVal{Depth: 0, Path: make([]string, 0, 5)}

	for queueFront.Len() > 0 && queueBack.Len() > 0 {
		fmt.Println(len(visitedFront), "Len")
		if answerPath, err := processQueue(queueFront, visitedFront, visitedBack, true); err != nil || answerPath != nil {
			// fmt.Println(visitedFront["https://en.wikipedia.org/wiki/Hemp"], " front\n")
			return answerPath, err
		}
		if answerPath, err := processQueue(queueBack, visitedBack, visitedFront, false); err != nil || answerPath != nil {
			// fmt.Println(visitedFront["https://en.wikipedia.org/wiki/Hemp"], " back\n")
			// fmt.Println(visitedBack["https://en.wikipedia.org/wiki/Hemp"], " back\n")

			return answerPath, err
		}
	}
	return nil, fmt.Errorf("path not found")
}

func processQueue(queue *list.List, visitedThis, visitedOther map[string]models.VisitedVal, isFront bool) ([]string, error) {
	if queue.Len() == 0 {
		return nil, nil
	}

	path := queue.Remove(queue.Front()).([]string)
	lastNode := path[len(path)-1]

	links, err := ScrapeWikipediaQuery(lastNode)
	if err != nil {
		return nil, err
	}

	for _, link := range links {

		// if otherValue, exist := visitedOther[link]; exist {
		// 	if !isFront {
		// 		if validatePath(link, path) {
		// 			combinedPath := combinePaths(otherValue.Path, link, reverseSlice(path))
		// 			return combinedPath, nil
		// 		} else {
		// 			if _, exist := visitedThis[link]; !exist {

		// 				visitedThis[link] = models.VisitedVal{Depth: visitedThis[lastNode].Depth + 1, Path: append(visitedThis[lastNode].Path, lastNode)}
		// 				newPath := append(make([]string, 0, len(path)+1), path...)
		// 				newPath = append(newPath, link)
		// 				queue.PushBack(newPath)
		// 			}
		// 			delete(visitedOther, link)

		// 		}
		// 	} else {
		// 		if validatePath(link, otherValue.Path) {
		// 			combinedPath := combinePaths(path, link, reverseSlice(otherValue.Path))
		// 			return combinedPath, nil
		// 		} else {
		// 			if _, exist := visitedThis[link]; !exist {

		// 				visitedThis[link] = models.VisitedVal{Depth: visitedThis[lastNode].Depth + 1, Path: append(visitedThis[lastNode].Path, lastNode)}
		// 				newPath := append(make([]string, 0, len(path)+1), path...)
		// 				newPath = append(newPath, link)
		// 				queue.PushBack(newPath)
		// 			}
		// 			delete(visitedOther, link)

		// 		}

		// 	}
		// } else {

		// 	if _, exist := visitedThis[link]; !exist {

		// 		visitedThis[link] = models.VisitedVal{Depth: visitedThis[lastNode].Depth + 1, Path: append(visitedThis[lastNode].Path, lastNode)}
		// 		newPath := append(make([]string, 0, len(path)+1), path...)
		// 		newPath = append(newPath, link)
		// 		queue.PushBack(newPath)
		// 	}
		// }

		if isFront {
			value, exist := visitedOther[link]
			if exist {
				if validatePath(link, value.Path) {
					return combinePaths(path, link, value.Path), nil
				}

			} else {
				if _, exist2 := visitedThis[link]; !exist2 {
					visitedThis[link] = models.VisitedVal{Depth: visitedThis[lastNode].Depth + 1, Path: append(visitedThis[lastNode].Path, lastNode)}
					newPath := make([]string, 0, len(path)+1)
					copy(newPath, path)
					newPath = append(newPath, link)
					queue.PushBack(newPath)
				}
			}
		} else {
			value, exist := visitedOther[link]
			if exist {
				if validatePath(link, path) {
					return combinePaths(value.Path, link, path), nil
				}
			} else {
				if _, exist2 := visitedThis[link]; !exist2 {
					visitedThis[link] = models.VisitedVal{Depth: visitedThis[lastNode].Depth + 1, Path: append(visitedThis[lastNode].Path, lastNode)}
					newPath := make([]string, 0, len(path)+1)
					copy(newPath, path)
					newPath = append(newPath, link)
					queue.PushBack(newPath)
				}
			}
		}
	}
	return nil, nil
}

func validatePath(meetingPoint string, backPath []string) bool {
	fullPath := append(backPath, meetingPoint)

	for i := len(fullPath) - 1; i >= 1; i-- {
		if !isValidTransition(fullPath[i], fullPath[i-1]) {
			return false
		}
	}
	return true
}

func isValidTransition(from, to string) bool {
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

func combinePaths(frontPath []string, link string, backPath []string) []string {

	frontPath = append(frontPath, link)
	combinedPath := append(frontPath, backPath...)
	return combinedPath
}

func reverseSlice(s []string) []string {
	a := make([]string, len(s))
	for i, v := range s {
		a[len(s)-1-i] = v
	}
	return a
}

// func BFS4(start string, goal string) ([]string, error) {
// 	queue := list.New()
// 	queue.PushBack([]models.Node{{Name: start, IsBegin: true}})
// 	queue.PushBack([]models.Node{{Name: goal, IsBegin: false}})

// 	// foundDepth := -1

// 	visited := make(map[string]models.VisitedVal)
// 	visited[start] = models.VisitedVal{Depth: 0, IsBegin: true, Path: make([]string, 0, 5)}
// 	visited[goal] = models.VisitedVal{Depth: 0, IsBegin: false, Path: make([]string, 0, 5)}

// 	for queue.Len() > 0 {
// 		path := queue.Remove(queue.Front()).([]models.Node)
// 		lastNode := path[len(path)-1]

// 		links, err := ScrapeWikipediaColly(lastNode.Name)
// 		if err != nil {
// 			return nil, err
// 		}

// 		for _, link := range links {
// 			value, exist := visited[link]
// 			if exist && value.IsBegin != lastNode.IsBegin {
// 				// foundDepth = value.Depth + len(path) + 1
// 				answerPath, err := combine(path, value.Path, link, lastNode.IsBegin)
// 				if err != nil {
// 					visited[link] = models.VisitedVal{Depth: visited[lastNode.Name].Depth + 1, IsBegin: lastNode.IsBegin, Path: append(visited[lastNode.Name].Path, lastNode.Name)}
// 					newPath := make([]models.Node, len(path)+1)
// 					copy(newPath, path)
// 					newPath[len(path)] = models.Node{Name: link, IsBegin: lastNode.IsBegin}
// 					queue.PushBack(newPath)
// 					continue
// 				}

// 				return answerPath, nil
// 			}

// 			if !exist {
// 				visited[link] = models.VisitedVal{Depth: visited[lastNode.Name].Depth + 1, IsBegin: lastNode.IsBegin, Path: append(visited[lastNode.Name].Path, lastNode.Name)}
// 				newPath := make([]models.Node, len(path)+1)
// 				copy(newPath, path)
// 				newPath[len(path)] = models.Node{Name: link, IsBegin: lastNode.IsBegin}
// 				queue.PushBack(newPath)
// 			}
// 		}
// 	}
// 	return nil, fmt.Errorf("path not found")
// }

func combine(arr []models.Node, arr2 []string, value string, isBegin bool) ([]string, error) {

	arr = append(arr, models.Node{Name: value, IsBegin: isBegin})
	arr2 = append(arr2, value)

	if isBegin {

		for ctr := len(arr2) - 1; ctr >= 1; ctr-- {
			linkWiki := arr2[ctr]
			links, err := ScrapeWikipediaColly(linkWiki)
			if err != nil {
				return nil, fmt.Errorf("Combiner path failed")
			}

			found := false
			for _, link := range links {
				if link == arr2[ctr-1] {
					found = true
					break
				}
			}

			if !found {
				return nil, fmt.Errorf("Combiner path failed")
			}
		}

	} else {

		for ctr := len(arr) - 1; ctr >= 1; ctr-- {
			linkWiki := arr[ctr]
			links, err := ScrapeWikipediaColly(linkWiki.Name)
			if err != nil {
				return nil, fmt.Errorf("Combiner path failed")

			}

			found := false
			for _, link := range links {
				if link == arr[ctr-1].Name {
					found = true
					break
				}
			}

			if !found {
				return nil, fmt.Errorf("Combiner path failed")

			}
		}

	}

	var length int = len(arr) + len(arr2) - 1
	fmt.Println(len(arr), " ", len(arr2))

	newPath := make([]string, length)

	if isBegin {
		lenBefore := 0

		for _, val := range arr {
			newPath[lenBefore] = val.Name
			lenBefore++
		}

		for i := len(arr2) - 2; i >= 0; i-- {
			newPath[lenBefore] = arr2[i]
			lenBefore++
		}
	} else {
		lenBefore := 0

		for _, val := range arr2 {
			newPath[lenBefore] = val
			lenBefore++
		}

		for i := len(arr) - 2; i >= 0; i-- {
			newPath[lenBefore] = arr[i].Name
			lenBefore++
		}
	}

	return newPath, nil

}

func combine2(arr []string, arr2 []string, value string) ([]string, error) {

	arr = append(arr, value)
	arr2 = append(arr2, value)

	for ctr := len(arr2) - 1; ctr >= 1; ctr-- {
		linkWiki := arr2[ctr]
		links, err := ScrapeWikipediaColly(linkWiki)
		if err != nil {
			return nil, fmt.Errorf("Combiner path failed")
		}

		found := false
		for _, link := range links {
			if link == arr2[ctr-1] {
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("Combiner path failed")
		}
	}

	var length int = len(arr) + len(arr2) - 1
	fmt.Println(len(arr), " ", len(arr2))

	newPath := make([]string, length)

	lenBefore := 0

	for _, val := range arr {
		newPath[lenBefore] = val
		lenBefore++
	}

	for i := len(arr2) - 2; i >= 0; i-- {
		newPath[lenBefore] = arr2[i]
		lenBefore++
	}

	return newPath, nil

}
