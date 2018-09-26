package main

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	gofight "gopkg.in/appleboy/gofight.v2"
)

var g *gin.Engine

func init() {
	g = gin.New()
	g.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	g.Use(middlewareCSRF())
}

func TestCSRFTokenIsSetInCookie(t *testing.T) {
	f := gofight.New()

	f.POST("/ping").
		SetHeader(gofight.H{
			"Accept": "text/html",
		}).
		Run(g, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			cookies := f.HeaderMap["Set-Cookie"]
			assert.Contains(t, cookies[0], "csrf=")
		})
}
