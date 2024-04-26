package controllers

import (
	"be/models"
	"be/services"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func SearchBFS(ctx *gin.Context) {
	var input models.SearchBodyRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error Bind JSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error Bind JSON"})
		return
	}
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error Validator"})
		return
	}

	// fullLinkStart := "https://en.wikipedia.org/wiki/" + input.Start
	// fullLinkGoal := "https://en.wikipedia.org/wiki/" + input.Goal

	// paths, err := services.MainBFS(fullPathStart, fullPathGoal)
	// paths, _ := services.BFS2(fullLinkStart, fullLinkGoal)
	start := time.Now()
	paths, countChecked := services.AsyncBFS6(input.Start, input.Goal)
	duration := time.Since(start).Milliseconds()

	if paths == nil {
		// fmt.Println(err.Error())

		ctx.JSON(http.StatusNotFound, gin.H{"found": false, "message": "path not found", "executionTime": duration})

		return
	}

	newPath := make([]string, 0, len(paths))
	for _, path := range paths {
		newPath = append(newPath, services.DecodePercentEncodedString(path))
	}

	ctx.JSON(http.StatusOK, gin.H{"paths": [][]string{newPath}, "found": true, "message": "path found", "executionTime": duration, "countChecked": countChecked})
	// return

}

func SearchBFSMulti(ctx *gin.Context) {
	var input models.SearchBodyRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error Bind JSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error Bind JSON"})
		return
	}
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error Validator"})
		return
	}

	// fullLinkStart := "https://en.wikipedia.org/wiki/" + input.Start
	// fullLinkGoal := "https://en.wikipedia.org/wiki/" + input.Goal

	// paths, err := services.MainBFS(fullPathStart, fullPathGoal)
	// paths, _ := services.BFS2(fullLinkStart, fullLinkGoal)
	start := time.Now()
	paths, countChecked := services.AsyncBFSMulti(input.Start, input.Goal)
	duration := time.Since(start).Milliseconds()

	if paths == nil {
		// fmt.Println(err.Error())

		ctx.JSON(http.StatusNotFound, gin.H{"found": false, "message": "path not found", "executionTime": duration})

		return
	}

	// newPaths := make([][]string, 0, len(paths))

	// for _, nPath := range paths {

	// 	newPath := make([]string, 0, len(nPath))
	// 	for _, path := range nPath {
	// 		newPath = append(newPath, services.DecodePercentEncodedString(path))
	// 	}

	// 	same := false
	// 	for _, p := range newPaths {
	// 		if services.CompareArrays(p, newPath) {
	// 			same = true
	// 			break
	// 		}
	// 	}

	// 	if !same {

	// 		newPaths = append(newPaths, newPath)
	// 	}
	// }

	uniquePaths := make(map[string]struct{}) // Using a map as a set
	newPaths := make([][]string, 0)

	for _, nPath := range paths {
		newPath := make([]string, 0, len(nPath))
		for _, path := range nPath {
			decodedPath := services.DecodePercentEncodedString(path)
			newPath = append(newPath, decodedPath)
		}

		// Create a string key to use in the map for uniqueness check
		pathKey := strings.Join(newPath, ",") // Join paths as a single string with commas
		if _, exists := uniquePaths[pathKey]; !exists {
			uniquePaths[pathKey] = struct{}{}
			newPaths = append(newPaths, newPath)
		}
	}

	// fmt.Println(newPaths)

	ctx.JSON(http.StatusOK, gin.H{"paths": newPaths, "found": true, "message": "path found", "executionTime": duration, "countChecked": countChecked})

}

func SearchBFS2(ctx *gin.Context) {
	var input models.SearchBodyRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error Bind JSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error Bind JSON"})
		return
	}
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error Validator"})
		return
	}

	// fullLinkStart := "https://en.wikipedia.org/wiki/" + input.Start
	// fullLinkGoal := "https://en.wikipedia.org/wiki/" + input.Goal

	// paths, err := services.MainBFS(fullPathStart, fullPathGoal)
	// paths, _ := services.BFS2(fullLinkStart, fullLinkGoal)
	start := time.Now()
	paths, countChecked := services.AsyncBFS5(input.Start, input.Goal)
	duration := time.Since(start).Milliseconds()

	if paths == nil {
		// fmt.Println(err.Error())

		ctx.JSON(http.StatusNotFound, gin.H{"found": false, "message": "path not found", "executionTime": duration})

		return
	}

	newPath := make([]string, 0, len(paths))
	for _, path := range paths {
		newPath = append(newPath, services.DecodePercentEncodedString(path))
	}

	ctx.JSON(http.StatusOK, gin.H{"paths": [][]string{newPath}, "found": true, "message": "path found", "executionTime": duration, "countChecked": countChecked})

}

func SearchBFS3(ctx *gin.Context) {
	var input models.SearchBodyRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error Bind JSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error Bind JSON"})
		return
	}
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error Validator"})
		return
	}

	// fullLinkStart := "https://en.wikipedia.org/wiki/" + input.Start
	// fullLinkGoal := "https://en.wikipedia.org/wiki/" + input.Goal

	// paths, err := services.MainBFS(fullPathStart, fullPathGoal)
	// paths, _ := services.BFS2(fullLinkStart, fullLinkGoal)
	start := time.Now()
	paths, countChecked := services.AsyncBFS7(input.Start, input.Goal)
	duration := time.Since(start).Milliseconds()

	if paths == nil {
		// fmt.Println(err.Error())

		ctx.JSON(http.StatusNotFound, gin.H{"found": false, "message": "path not found", "executionTime": duration})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"paths": [][]string{paths}, "found": true, "message": "path found", "executionTime": duration, "countChecked": countChecked})
	// return

}

func SearchDoubleBFS(ctx *gin.Context) {
	var input models.SearchBodyRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error Bind JSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error Bind JSON"})
		return
	}
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error Validator"})
		return
	}

	fullLinkStart := "https://en.wikipedia.org/wiki/" + input.Start
	fullLinkGoal := "https://en.wikipedia.org/wiki/" + input.Goal

	// paths, err := services.MainBFS(fullPathStart, fullPathGoal)
	paths, _ := services.BFS2(fullLinkStart, fullLinkGoal)
	// paths := services.AsyncBFS5(fullLinkStart, fullLinkGoal)

	if paths == nil {
		// fmt.Println(err.Error())

		ctx.JSON(http.StatusNotFound, gin.H{"found": false, "message": "path not found"})

		return
	}

	newPath := make([]string, 0, len(paths))
	for _, path := range paths {
		newPath = append(newPath, services.DecodePercentEncodedString(path))
	}

	ctx.JSON(http.StatusOK, gin.H{"paths": [][]string{newPath}, "found": true, "message": "path found"})
	// return

}
