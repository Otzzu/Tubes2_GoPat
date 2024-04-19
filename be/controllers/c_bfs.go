package controllers

import (
	"be/models"
	"be/services"
	"fmt"
	"net/http"

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

	fullLinkStart := "https://en.wikipedia.org/wiki/" + input.Start
	fullLinkGoal := "https://en.wikipedia.org/wiki/" + input.Goal

	// paths, err := services.MainBFS(fullPathStart, fullPathGoal)
	// paths, err := services.BFS3(fullPathStart, fullPathGoal)
	paths, err := services.BreadthFirstSearch(fullLinkStart, fullLinkGoal)

	if (err != nil) {
		fmt.Println(err.Error())

		ctx.JSON(http.StatusNotFound, gin.H{"found": false, "message": "path not found"})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"paths": paths ,"found": true, "message": "path found"})
	return

}