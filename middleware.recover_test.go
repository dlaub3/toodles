package main

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	gofight "gopkg.in/appleboy/gofight.v2"
)

func getRecoverRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r = gin.New()
	r.Use(middlewareRecover())
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	r.GET("/ping", func(c *gin.Context) {
		panic("Something went wrong.")
	})
	return r
}

func TestMiddlewareRecover(t *testing.T) {
	f := gofight.New()
	r := getRecoverRouter()

	f.GET("/ping").
		SetHeader(gofight.H{
			"Accept": "text/html",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Contains(t, f.Body.String(), "😑 oh snap! Please try again.")
			assert.Equal(t, f.Code, http.StatusInternalServerError)
		})
}
