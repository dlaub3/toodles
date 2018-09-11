package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CsrfToken binds with form submit csrf
type CsrfToken struct {
	CsrfToken string `form:"csrf" json:"csrf"`
}

// IsCSRFTokenValid checks the request for a valid CSRF token
func IsCSRFTokenValid(c *gin.Context) bool {
	csrfToken := CsrfToken{}
	// save the request body
	body, _ := ioutil.ReadAll(c.Request.Body)
	// restore the request body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	c.Bind(&csrfToken)
	// restore the request body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return csrfToken.CsrfToken == c.Keys["csrftoken"].(string)
}

// InvalidateCookies for JWT and CSRF
func InvalidateCookies(c *gin.Context) {
	InvalidateCSRF(c)
	InvalidateJWT(c)
}

// InvalidateCSRF Cookie
func InvalidateCSRF(c *gin.Context) {
	csrfcookie := http.Cookie{
		Name:    "csrf",
		Path:    "/",
		Expires: time.Now().UTC(),
	}
	http.SetCookie(c.Writer, &csrfcookie)
}

// InvalidateJWT Cooke
func InvalidateJWT(c *gin.Context) {
	jwtcookie := http.Cookie{
		Name:    "token",
		Path:    "/",
		Expires: time.Now().UTC(),
	}
	http.SetCookie(c.Writer, &jwtcookie)
}

// ShowErrorPage for bad requests
func ShowErrorPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Our servers are busy please stand bye. 😞",
	}, "error.html")
}

// HandleUnauthorized request repsonses
func HandleUnauthorized(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "application/json" {
		c.Set("error", "unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{})
	} else {
		c.Redirect(302, "/login")
	}
}
