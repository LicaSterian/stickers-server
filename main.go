package main

import (
	"github.com/gin-gonic/gin"
	"./routers"
	"./helpers"
)

func main() {
	helpers.ConnectToMongo()

	r := gin.Default()
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("./views/*")

	routers.InitStickersRouter(r)

	r.Run("0.0.0.0:8082")
}