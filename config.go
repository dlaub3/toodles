package main

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

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

	binding.Validator = new(defaultValidator)
}
