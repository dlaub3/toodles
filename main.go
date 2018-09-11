package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func main() {

	// Set the router as the default one provided by Gin
	r = gin.Default()
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

// Render one of HTML, JSON or XML based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func render(c *gin.Context, data gin.H, templateName string) {
	error, _ := c.Get("error")
	data["error"] = error

	var httpStatus int
	status, _ := c.Get("httpStatus")

	if status != nil {
		httpStatus = status.(int)
	}

	cookie, _ := c.Request.Cookie("token")

	if cookie != nil {
		data["loggedin"] = true
	}

	if httpStatus < 1 {
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
	}

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(httpStatus, data)
	case "application/xml":
		c.XML(httpStatus, data)
	default:
		// Respond with HTML
		c.HTML(httpStatus, templateName, data)
	}

}
