package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRouters() {

	r := gin.Default()
	r.GET("/hello", helloHandler)
	r.Run()

}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, Valhalla",
	})

}
