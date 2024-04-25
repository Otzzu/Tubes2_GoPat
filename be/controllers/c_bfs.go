package controllers

import (
	"be/models"
	"be/services"
	"fmt"
	"net/http"
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

	newPaths := make([][]string, 0, len(paths))

	for _, nPath := range paths {

		newPath := make([]string, 0, len(nPath))
		for _, path := range nPath {
			newPath = append(newPath, services.DecodePercentEncodedString(path))
		}

		newPaths = append(newPaths, newPath)
	}

	ctx.JSON(http.StatusOK, gin.H{"paths": newPaths, "found": true, "message": "path found", "executionTime": duration, "countChecked": countChecked})
	// return

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
	// return

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
