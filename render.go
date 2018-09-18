package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// send the http response
func render(c *gin.Context, data gin.H, templateName string) {

	if error, _ := c.Get("error"); error != nil {
		data["error"] = error
	}

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

	// over ride HTTP status with alternate httpStatus
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
