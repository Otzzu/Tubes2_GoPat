package routes

import (
	"be/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	
	config := cors.DefaultConfig();

	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Authorization, content-type"}
	router.Use(cors.New(config))
	router.Use(gin.Recovery())

	search := router.Group("/search")
	search.POST("/BFS/multi", controllers.SearchBFSMulti)
	search.POST("/BFS", controllers.SearchBFS)
	search.POST("/IDS", controllers.SearchIDS)
}