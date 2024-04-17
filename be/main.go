package main

import (
	"be/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)


func main(){
	app := gin.Default()

	routes.Init(app)

	var port string = "3000"

	fmt.Println("Hello World!")
	fmt.Println("Server running on port " + port)
	app.Run(":" + port)
}