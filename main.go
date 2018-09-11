package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func main() {

	// Set the router as the default one provided by Gin
	r = gin.Default()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	r.LoadHTMLGlob("templates/*")

	// Loads assets
	r.Static("/assets", "./assets")

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	r.Run()

}

// send the http response
func render(c *gin.Context, data gin.H, templateName string) {

	error, _ := c.Get("error")
	data["error"] = error

	cookie, _ := c.Request.Cookie("token")
	if cookie != nil {
		data["loggedin"] = true
	}

	var httpStatus int
	switch c.Request.Method {
	case "GET":
		httpStatus = http.StatusOK
	case "POST":
		httpStatus = http.StatusCreated
	case "PUT":
		httpStatus = http.StatusCreated
	case "DELETE":
		httpStatus = http.StatusOK
	}

	// get alternate httpStatus
	status, _ := c.Get("httpStatus")
	if status != nil {
		httpStatus = status.(int)
	}

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(httpStatus, data)
	case "application/xml":
		c.XML(httpStatus, data)
	default:
		c.HTML(httpStatus, templateName, data)
	}

}
