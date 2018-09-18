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

	// Set the router as the default one provided by Gin
	r = gin.Default()
	binding.Validator = new(defaultValidator)
	r.Use(middlewareErrors())

	// Process the templates at the start so that they don't have to be loaded
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	initializeRoutes()
	dbConnect()
	r.Run()
}
