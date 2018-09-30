package main

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	gofight "gopkg.in/appleboy/gofight.v2"
)

func getCSRFRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r = gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.Use(middlewareCSRF())
	return r
}

func TestMiddlewareCSRF(t *testing.T) {
	r := getCSRFRouter()
	f := gofight.New()

	f.POST("/ping").
		SetHeader(gofight.H{
			"Accept": "text/html",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			cookies := f.HeaderMap["Set-Cookie"]
			assert.Contains(t, cookies[0], "csrf=")
		})
}
