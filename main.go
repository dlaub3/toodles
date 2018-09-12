package main

import (
	"github.com/gin-gonic/gin"
)

var r *gin.Engine
var config Config

func main() {

	config = Config{}
	config.Read()

	// Set the router as the default one provided by Gin
	r = gin.Default()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Process the templates at the start so that they don't have to be loaded
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	initializeRoutes()
	dbConnect()
	r.Run()

}
