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

	if c.Keys["uid"] != nil {
		data["loggedin"] = true
	}

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data)
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data)
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}

}
