package main

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	gofight "gopkg.in/appleboy/gofight.v2"
)

func getRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	g := gin.New()
	g.POST("/isCSRFTokenValid", isCSRFTokenValidTest)
	g.GET("/csrf", csrfTest)
	g.GET("/invalidateCSRF", invalidateCSRFTest)
	g.GET("/invalidateJWT", invalidateJWTTest)
	g.GET("/invalidateCookies", invalidateCookiesTest)
	g.Run()
	return g
}

func invalidateCSRFTest(c *gin.Context) {
	invalidateCSRF(c)
	c.JSON(http.StatusOK, gin.H{})
}

func invalidateJWTTest(c *gin.Context) {
	invalidateJWT(c)
	c.JSON(http.StatusOK, gin.H{})
}

func invalidateCookiesTest(c *gin.Context) {
	invalidateCookies(c)
	c.JSON(http.StatusOK, gin.H{})
}

func csrfTest(c *gin.Context) {
	cookie, _ := csrf(c)
	c.JSON(http.StatusOK, gin.H{
		"cookie": cookie,
	})
}

func isCSRFTokenValidTest(c *gin.Context) {
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
	g := getRouter()
	f := gofight.New()

	f.POST("/isCSRFTokenValid").
		SetForm(gofight.H{
			"csrf": "aalkj3035555hwwe002jl21",
		}).
		SetHeader(gofight.H{
			"Cookie": "csrf=aalkj3035555hwwe002jl21;",
		}).
		Run(g, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Contains(t, f.Body.String(), "success")
		})
}

func TestIsCSRFTokenValidFailure(t *testing.T) {
	g := getRouter()
	f := gofight.New()

	f.POST("/isCSRFTokenValid").
		SetForm(gofight.H{
			"csrf": "10aalkj3035555hwwe002jl21",
		}).
		SetHeader(gofight.H{
			"Cookie": "csrf=aalkj3035555hwwe002jl21;",
		}).
		Run(g, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Contains(t, f.Body.String(), "failure")
		})
}

func TestCsrfSetInCookie(t *testing.T) {
	g := getRouter()
	f := gofight.New()

	f.GET("/csrf").
		Run(g, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			cookies := f.HeaderMap["Set-Cookie"]
			assert.Contains(t, cookies[0], "csrf=")
		})
}

func TestInvalidateCSRF(t *testing.T) {
	g := getRouter()
	f := gofight.New()

	f.GET("/invalidateCSRF").
		Run(g, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			cookies := f.HeaderMap["Set-Cookie"]
			assert.Contains(t, cookies[0], "csrf=;")
		})
}

func TestInvalidateJWT(t *testing.T) {
	g := getRouter()
	f := gofight.New()

	f.GET("/invalidateJWT").
		Run(g, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			cookies := f.HeaderMap["Set-Cookie"]
			assert.Contains(t, cookies[0], "token=;")
		})
}

func TestInvalidateCookies(t *testing.T) {
	g := getRouter()
	f := gofight.New()

	f.GET("/invalidateCookies").
		Run(g, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			cookies := f.HeaderMap["Set-Cookie"]
			assert.Contains(t, cookies[0], "csrf=;")
			assert.Contains(t, cookies[1], "token=;")
		})
}
