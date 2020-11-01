package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"go-trivia/trivia"

	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
)

var slide_data *trivia.SlideList
var game_data *trivia.Game

const slides_path = "./slides.json"
const game_path = "./game.json"
const Port = "8080"

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		if strings.HasPrefix(addr.String(), "192") {
			address := strings.TrimSuffix(addr.String(), "/24")
			if len(Port) > 0 {
				address = fmt.Sprintf("%v:%v", address, Port)
			}
			return address
		}
	}
	return "unknown"
}

func myGame(c *gin.Context, playerName string) {
	data, err := game_data.ForPlayer(playerName)
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

	// add player if they don't exist
	_, found := game_data.FindPlayer(player)
	if !found {
		game_data.AddPlayer(trivia.Player{Name: player})
	}

	if err := game_data.AddAnswer(player, slide, answer); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Could not set answer: %v: %v", message, err))
		return
	}
	myGame(c, player)
}

func main() {
	var err error

	fmt.Println("\n******************************************")
	fmt.Printf("* Local address: %v\n", localIP())
	fmt.Println("******************************************\n")

	slide_data, err = trivia.GetSlideJSON(slides_path)
	if err != nil {
		fmt.Println("Could not open slide data: ", slides_path)
		fmt.Println(err)
		os.Exit(1)
	}

	game_data, err = trivia.GetGameJSON(game_path)
	if err != nil {
		fmt.Println("Could not open slide data: ", game_path)
		fmt.Println(err)
		os.Exit(1)
	}
	trivia.NewGame(game_data, slide_data)

	router := gin.Default()
	router.Use(favicon.New("./public/favicon.ico"))

	// for debugging, limit caching
	// router.Use(func(c *gin.Context) {
	// 	c.Header("Cache-Control", "max-age=5")
	// })

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/slides", func(c *gin.Context) {
		c.JSON(http.StatusOK, game_data.Slides)
	})

	router.GET("/game/:playerName", func(c *gin.Context) {
		myGame(c, c.Param("playerName"))
	})

	router.POST("/game/:playerName", func(c *gin.Context) {
		p := trivia.Player{Name: c.Param("playerName")}
		game_data.AddPlayer(p)
		myGame(c, p.Name)
	})

	router.POST("/answer/:slide", postAnswer)

	router.Static("/public", "./public")
	router.StaticFile("/", "index.html")

	//router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	// TODO: add port setting
	srv := &http.Server{
		Addr:    ":" + Port,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	if err := game_data.Save(game_path); err != nil {
		fmt.Println("Error saving game:", err)
	}

	// The context is used to inform the server it has 2 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
