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

func SearchIDS(ctx *gin.Context) {
	var input models.SearchBodyRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error Bind JSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error Bind JSON"})
		return
	}
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		fmt.Println("Error Validator")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error Validator"})
		return
	}


	// paths, err := services.IDS(fullPathStart, fullPathGoal, 5)
	start := time.Now()

	paths, countChecked := services.IDS(input.Start, input.Goal, 8)
	duration := time.Since(start).Milliseconds()


	if (paths == nil) {
		// fmt.Println(err.Error())

		ctx.JSON(http.StatusNotFound, gin.H{"found": false, "message": "path not found", "executionTime" : duration})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"paths": paths ,"found": true, "countChecked" : countChecked, "message": "path found", "executionTime" : duration})
	return

}