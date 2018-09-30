package main

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	gofight "gopkg.in/appleboy/gofight.v2"
)

func getRouter() *gin.Engine {
	// gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/", csrfTokenTest)
	r.Run()
	return r
}

func csrfTokenTest(c *gin.Context) {
	csrfToken, _ := c.Request.Cookie("csrf")
	c.Set("csrftoken", csrfToken.Value)
	validCSRF := isCSRFTokenValid(c)

	if validCSRF {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "failure",
		})
	}
}

func TestIsCSRFTokenValidSuccess(t *testing.T) {
	r := getRouter()
	f := gofight.New()

	f.POST("/").
		SetForm(gofight.H{
			"csrf": "aalkj3035555hwwe002jl21",
		}).
		SetHeader(gofight.H{
			"Cookie": "csrf=aalkj3035555hwwe002jl21;",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Contains(t, f.Body.String(), "success")
		})
}

func TestIsCSRFTokenValidFailure(t *testing.T) {
	r := getRouter()
	f := gofight.New()

	f.POST("/").
		SetForm(gofight.H{
			"csrf": "10aalkj3035555hwwe002jl21",
		}).
		SetHeader(gofight.H{
			"Cookie": "csrf=aalkj3035555hwwe002jl21;",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Contains(t, f.Body.String(), "failure")
		})
}
