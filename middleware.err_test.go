package main

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"

	gofight "gopkg.in/appleboy/gofight.v2"
)

func getErrRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	g := gin.New()
	g.Use(middlewareErrors())
	g.LoadHTMLGlob("templates/*")
	g.Static("/assets", "./assets")
	g.GET("/ErrorTypePublic", func(c *gin.Context) {
		c.AbortWithError(http.StatusBadRequest, errorInternalError).SetType(gin.ErrorTypePublic)
	})
	g.POST("/ErrorTypeBind", func(c *gin.Context) {
		binding.Validator = new(defaultValidator)
		type User struct {
			Username string `bson:"username" form:"username" json:"username" binding:"required,min=10,max=25"`
			Password string `bson:"password" form:"password" json:"password" binding:"required,alphanum,min=10,max=30"`
		}

		user := User{}
		c.Bind(&user)
		c.JSON(http.StatusOK, gin.H{
			"username": user.Username,
			"password": user.Password,
		})
	})
	g.POST("/ErrorDefault", func(c *gin.Context) {
		c.AbortWithError(http.StatusInternalServerError, errorInternalError).SetType(gin.ErrorTypePrivate)
	})

	return g
}

func TestMiddlewareErrTypePublic(t *testing.T) {
	g := getErrRouter()
	f := gofight.New()

	f.GET("/ErrorTypePublic").
		SetHeader(gofight.H{
			"Accept": "text/html",
		}).
		Run(g, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Contains(t, f.Body.String(), "Our servers are busy please stand bye. 😞")
			assert.Equal(t, f.Code, http.StatusBadRequest)
		})
}

func TestMiddlewareErrTypeBind(t *testing.T) {
	g := getErrRouter()
	f := gofight.New()

	f.POST("/ErrorTypeBind").
		SetForm(gofight.H{
			"username": "username",
			"password": "password",
		}).
		SetHeader(gofight.H{
			"Accept": "application/json",
		}).
		Run(g, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Contains(t, f.Body.String(), "Password must be longer than 10")
			assert.Contains(t, f.Body.String(), "Username must be longer than 10")
		})
}

func TestMiddlewareErrDefault(t *testing.T) {
	g := getErrRouter()
	f := gofight.New()

	f.POST("/ErrorDefault").
		SetHeader(gofight.H{
			"Accept": "application/json",
		}).
		Run(g, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Contains(t, f.Body.String(), "Our servers are busy please stand bye. 😞")
		})
}
