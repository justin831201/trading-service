package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "ok")
	})
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
