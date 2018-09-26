package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
	gofight "gopkg.in/appleboy/gofight.v2"
)

func init() {
	r = gin.New()
	initConfig()
	initializeRoutes()
	dbConnect()
	binding.Validator = new(defaultValidator)
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestPing(t *testing.T) {
	w := performRequest(r, "GET", "/")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIndexPage(t *testing.T) {
	w := performRequest(r, "GET", "/")
	p, _ := ioutil.ReadAll(w.Body)
	assert.Contains(t, string(p), "Toodles, stay organized. Get stuff done!")
}

func TestSignupRequresEmail(t *testing.T) {
	f := gofight.New()

	f.POST("/signup").
		SetForm(gofight.H{
			"email":    "",
			"password": "testing123",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(f.Body.String())

			assert.Contains(t, string(data), "Email is required")
			assert.Equal(t, http.StatusBadRequest, f.Code)
		})
}

func TestSignupRequresPassword(t *testing.T) {
	f := gofight.New()

	f.POST("/signup").
		SetForm(gofight.H{
			"email":    "test@example.com",
			"password": "",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(f.Body.String())

			assert.Contains(t, string(data), "Password is required")
			assert.Equal(t, http.StatusBadRequest, f.Code)
		})
}

func TestLoginFailure(t *testing.T) {
	f := gofight.New()

	f.POST("/login").
		SetForm(gofight.H{
			"username": "admin@example.com",
			"password": "testing123",
		}).
		SetHeader(gofight.H{
			"Accept": "text/html",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(f.Body.String())

			assert.Contains(t, string(data), "Your username and password do not match.")
			assert.Equal(t, http.StatusUnauthorized, f.Code)
		})
}

func TestLoginRedirectOnSuccess(t *testing.T) {
	f := gofight.New()

	f.POST("/login").
		SetForm(gofight.H{
			"username": "admin@example.com",
			"password": "password",
		}).
		SetHeader(gofight.H{
			"Accept": "text/html",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusFound, f.Code)
		})
}
