package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	gofight "gopkg.in/appleboy/gofight.v2"
)

func init() {
	initConfig()
	initDB()
	initRoutes()
}

func TestShowHomePage(t *testing.T) {
	f := gofight.New()

	f.GET("/").
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Contains(t, f.Body.String(), "Toodles, stay organized. Get stuff done!")
			assert.Equal(t, http.StatusOK, f.Code)
		})
}

func TestShowSignupPage(t *testing.T) {
	f := gofight.New()

	f.GET("/signup").
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Contains(t, f.Body.String(), "Signup to start toodling today.")
			assert.Equal(t, http.StatusOK, f.Code)
		})
}

func TestShowLoginPage(t *testing.T) {
	f := gofight.New()

	f.GET("/login").
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Contains(t, f.Body.String(), "Login to access your toodles.")
			assert.Equal(t, http.StatusOK, f.Code)
		})
}

func TestSignupRequresEmail(t *testing.T) {
	f := gofight.New()

	f.POST("/signup").
		SetForm(gofight.H{
			"email":    "",
			"password": "testing123",
			"csrf":     "aalkj3035555hwwe002jl21",
		}).
		SetHeader(gofight.H{
			"Accept": "text/html",
			"Cookie": "csrf=aalkj3035555hwwe002jl21;",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Contains(t, f.Body.String(), "Email is required")
			assert.Equal(t, http.StatusBadRequest, f.Code)
		})
}

func TestSignupRequresPassword(t *testing.T) {
	f := gofight.New()

	f.POST("/signup").
		SetForm(gofight.H{
			"email":    "test@example.com",
			"password": "",
			"csrf":     "aalkj3035555hwwe002jl21",
		}).
		SetHeader(gofight.H{
			"Accept": "text/html",
			"Cookie": "csrf=aalkj3035555hwwe002jl21;",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Contains(t, f.Body.String(), "Password is required")
			assert.Equal(t, http.StatusBadRequest, f.Code)
		})
}

func TestLoginFailure(t *testing.T) {
	f := gofight.New()

	f.POST("/login").
		SetForm(gofight.H{
			"username": "admin@example.com",
			"password": "testing123",
			"csrf":     "aalkj3035555hwwe002jl21",
		}).
		SetHeader(gofight.H{
			"Accept": "text/html",
			"Cookie": "csrf=aalkj3035555hwwe002jl21;",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Contains(t, f.Body.String(), "Your username and password do not match.")
			assert.Equal(t, http.StatusUnauthorized, f.Code)
		})
}

func TestLoginRedirectOnSuccess(t *testing.T) {
	f := gofight.New()

	f.POST("/login").
		SetForm(gofight.H{
			"username": "admin@example.com",
			"password": "password",
			"csrf":     "aalkj3035555hwwe002jl21",
		}).
		SetHeader(gofight.H{
			"Accept": "text/html",
			"Cookie": "csrf=aalkj3035555hwwe002jl21;",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusFound, f.Code)
		})
}

func TestLogout(t *testing.T) {
	f := gofight.New()

	f.GET("/logout").
		SetHeader(gofight.H{
			"Accept": "text/html",
			"Cookie": "csrf=aalkj3035555hwwe002jl21;",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			cookies := f.HeaderMap["Set-Cookie"]
			assert.Contains(t, cookies[1], "token=;")
			assert.Contains(t, cookies[0], "csrf=;")
			assert.Equal(t, http.StatusFound, f.Code)
		})
}
