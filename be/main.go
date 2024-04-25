package main

import (
	"be/routes"
	"fmt"
	"runtime"
	"github.com/gin-gonic/gin"
)


func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := gin.Default()

	routes.Init(app)

	var port string = "8080"

	fmt.Println("Hello World!")
	fmt.Println("Server running on port " + port)
	app.Run(":" + port)
}