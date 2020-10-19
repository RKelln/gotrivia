package main

import "github.com/gin-gonic/gin"

import (
	"fmt"
	"net/http"
	"os"
)

func postAnswer(c *gin.Context) {
	player := c.PostForm("player")
	slide := c.Param("slide")
	answer := c.PostForm("answer")
	message := player + " anwered " + slide + " with " + answer
	c.String(http.StatusOK, message)
}

func main() {
	slide_data, err := GetSlideJSON(slides_path)
	if err != nil {
		fmt.Println("Could not open slide data: ", slides_path)
		fmt.Println(err)
		os.Exit(1)
	}

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/slides", func(c *gin.Context) {
		c.JSON(http.StatusOK, slide_data)
	})

	router.POST("/answer/:slide", postAnswer)

	router.Static("/public", "./public")
	router.StaticFile("/", "index.html")
	router.StaticFile("/favicon.ico", "./public/favicon.ico")

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
