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

	start := time.Now()
	paths, countChecked := services.AsyncBFSMulti(input.Start, input.Goal)
	duration := time.Since(start).Milliseconds()

	if paths == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"found": false, "message": "path not found", "executionTime": duration})
		return
	}

	uniquePaths := make(map[string]struct{}) 
	newPaths := make([][]string, 0)

	for _, nPath := range paths {
		newPath := make([]string, 0, len(nPath))
		for _, path := range nPath {
			decodedPath := services.DecodePercentEncodedString(path)
			newPath = append(newPath, decodedPath)
		}

		// ensure uniqueness
		pathKey := strings.Join(newPath, ",") 
		if _, exists := uniquePaths[pathKey]; !exists {
			uniquePaths[pathKey] = struct{}{}
			newPaths = append(newPaths, newPath)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"paths": newPaths, "found": true, "message": "path found", "executionTime": duration, "countChecked": countChecked})
}

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

	start := time.Now()
	paths, countChecked := services.AsyncBFS(input.Start, input.Goal)
	fmt.Println("TES")
	duration := time.Since(start).Milliseconds()

	if paths == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"found": false, "message": "path not found", "executionTime": duration})

		return
	}

	newPath := make([]string, 0, len(paths))
	for _, path := range paths {
		newPath = append(newPath, services.DecodePercentEncodedString(path))
	}

	ctx.JSON(http.StatusOK, gin.H{"paths": [][]string{newPath}, "found": true, "message": "path found", "executionTime": duration, "countChecked": countChecked})
}