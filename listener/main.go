package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rmeharg/gomessage/sub"
)

func server() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func main() {
	fmt.Println("Hello, main world!")

	sub.HelloExportedFunction()

	server()
}
