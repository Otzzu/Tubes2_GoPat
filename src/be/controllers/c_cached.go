package controllers

import (
	"be/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ClearCached(ctx *gin.Context) {
	services.LinkCache = make(map[string][]string)

	ctx.JSON(http.StatusOK, gin.H{"message": "clear cached success"})

}