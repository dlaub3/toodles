package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CsrfToken binds with form submit csrf
type csrfToken struct {
	CsrfToken string `form:"csrf" json:"csrf"`
}

// IsCSRFTokenValid checks the request for a valid CSRF token
func isCSRFTokenValid(c *gin.Context) bool {
	csrfToken := csrfToken{}
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
func invalidateCookies(c *gin.Context) {
	invalidateCSRF(c)
	invalidateJWT(c)
}

// InvalidateCSRF Cookie
func invalidateCSRF(c *gin.Context) {
	csrfcookie := http.Cookie{
		Name:    "csrf",
		Path:    "/",
		Expires: time.Now().UTC(),
	}
	http.SetCookie(c.Writer, &csrfcookie)
}

// InvalidateJWT Cooke
func invalidateJWT(c *gin.Context) {
	jwtcookie := http.Cookie{
		Name:    "token",
		Path:    "/",
		Expires: time.Now().UTC(),
	}
	http.SetCookie(c.Writer, &jwtcookie)
}

// ShowErrorPage for bad requests
func showErrorPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Our servers are busy please stand bye. 😞",
	}, "error.html")
}

// HandleUnauthorized request repsonses
func handleUnauthorized(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "application/json" {
		c.Set("error", "unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{})
	} else {
		c.Redirect(302, "/login")
	}
}
