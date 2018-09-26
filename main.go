package main

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var r *gin.Engine
var config Config

func main() {
	initConfig()

	r = gin.New()
	r.Use(middlewareRecover())
	r.Use(gin.Logger())
	r.Use(middlewareCSRF())
	r.Use(middlewareErrors())

	binding.Validator = new(defaultValidator)

	// Process the templates at the start so that they don't have to be loaded
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	initializeRoutes()
	dbConnect()
	r.Run()
}

func initConfig() {
	config = Config{}
	config.Read()
	config.LogPath = config.LogPath + "/gin/"

	// app log to file
	fh, err := os.OpenFile(config.LogPath+"app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	log.SetOutput(fh)
	log.Println("Logging setup successfully :)")

	// Log HTTP requests to file
	gin.DisableConsoleColor()

	f, err := os.Create(config.LogPath + "http.log")
	if err != nil {
		log.Fatal(err)
	}
	gin.DefaultWriter = io.MultiWriter(f)
}
