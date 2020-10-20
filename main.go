package main

import "github.com/gin-gonic/gin"

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var slide_data *Slides
var game_data *Game

func myGame(c *gin.Context, playerName string) {
	data, err := game_data.forPlayer(playerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, data)
}

func postAnswer(c *gin.Context) {
	player := c.PostForm("player")
	slide, err := strconv.Atoi(c.Param("slide"))
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid slide param: %v: %v", c.Param("slide"), err))
		return
	}
	answer, err := strconv.Atoi(c.PostForm("answer"))
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid answer param: %v: %v", c.PostForm("answer"), err))
		return
	}

	message := fmt.Sprintf("%v answered slide %v with %v", player, slide, answer)
	fmt.Println(message)

	if err := game_data.addAnswer(player, slide, answer); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Could not set answer: %v: %v", message, err))
		return
	}
	fmt.Print(game_data)

	myGame(c, player)
}

func main() {
	var err error

	slide_data, err = GetSlideJSON(slides_path)
	if err != nil {
		fmt.Println("Could not open slide data: ", slides_path)
		fmt.Println(err)
		os.Exit(1)
	}

	game_data, err = GetGameJSON(game_path)
	if err != nil {
		fmt.Println("Could not open slide data: ", game_path)
		fmt.Println(err)
		os.Exit(1)
	}
	NewGame(game_data, slide_data)
	game_data.addPlayer(Player{Name: "player1"})

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/slides", func(c *gin.Context) {
		c.JSON(http.StatusOK, game_data.Slides)
	})

	router.GET("/game/:g/:playerName", func(c *gin.Context) {
		myGame(c, c.Param("playerName"))
	})

	router.POST("/answer/:slide", postAnswer)

	router.Static("/public", "./public")
	router.StaticFile("/", "index.html")
	router.StaticFile("/favicon.ico", "./public/favicon.ico")

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
