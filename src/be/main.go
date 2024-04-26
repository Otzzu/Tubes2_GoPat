package main

import (
	"be/repository"
	"be/routes"
	"fmt"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)


func main(){
	host := os.Getenv("DATABASE_HOST")
	portDb := os.Getenv("DATABASE_PORT")
	stringConnect := fmt.Sprintf("user=postgres dbname=mydb sslmode=disable password=postgres host=%s port=%s", host, portDb)
	db, err := sqlx.Connect("postgres", stringConnect)
    if err != nil {
        panic(err)
    }

	defer db.Close()

	repository.Db = db
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := gin.Default()

	routes.Init(app)

	var port string = "8080"

	fmt.Println("Hello World!")
	fmt.Println("Server running on port " + port)
	app.Run(":" + port)
}