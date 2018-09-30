package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dlaub3/toodles/crypt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// csrfToken binds with form submit csrf
type csrfToken struct {
	CsrfToken string `form:"csrf" json:"csrf"`
}

// isCSRFTokenValid checks the request for a valid CSRF token
func isCSRFTokenValid(c *gin.Context) bool {
	var err error
	formValue := csrfToken{}
	// save the request body
	body, err := ioutil.ReadAll(c.Request.Body)
	// restore the request body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	err = c.Bind(&formValue)
	// restore the request body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		log.Panic(err)
	}

	csrfCookie, _ := c.Get("csrftoken")

	return formValue.CsrfToken == csrfCookie.(string)
}

// csrf sets a csrf token in a cookie
func csrf(c *gin.Context) (*http.Cookie, error) {

	cookie, err := c.Request.Cookie("csrf")

	if err != nil {
		expire := time.Now().UTC().Add(time.Hour)
		maxage := int(expire.Unix() - time.Now().Unix())

		csrf, err := crypt.GenerateRandomString(32)
		cookie := http.Cookie{
			Name:     "csrf",
			Value:    csrf,
			Path:     "/",
			Expires:  expire,
			MaxAge:   maxage,
			HttpOnly: false, // only access with the secure option
			Secure:   false, //@dev change when in prod mode
			// No support for SameSite yet https://golang.org/src/net/http/cookie.go
		}
		http.SetCookie(c.Writer, &cookie)
		return &cookie, err
	}
	return cookie, err
}

// invalidateCookies for JWT and CSRF
func invalidateCookies(c *gin.Context) {
	invalidateCSRF(c)
	invalidateJWT(c)
}

// invalidateCSRF Cookie
func invalidateCSRF(c *gin.Context) {
	csrfcookie := http.Cookie{
		Name:    "csrf",
		Path:    "/",
		Expires: time.Now().UTC(),
	}
	http.SetCookie(c.Writer, &csrfcookie)
}

// invalidateJWT Cooke
func invalidateJWT(c *gin.Context) {
	jwtcookie := http.Cookie{
		Name:    "token",
		Path:    "/",
		Expires: time.Now().UTC(),
	}
	http.SetCookie(c.Writer, &jwtcookie)
}

// showErrorPage for bad requests
func showErrorPage(c *gin.Context) {
	render(c, gin.H{
		"error": "Our servers are busy please stand bye. 😞",
	}, "error.html")
}

func getValidationErrorMsg(e error) map[string]string {
	list := make(map[string]string)
	for _, field := range e.(validator.ValidationErrors) {
		list[field.Field()] = validationErrorToText(field)
	}
	return list
}

func getGeneralErrorMsg(key string) map[string]string {
	list := make(map[string]string)
	errMsgs := map[string]string{
		"database": "There was an error processing your request, please try again.",
	}
	list["general"] = errMsgs[key]
	return list
}

// handleUnauthorized request repsonses
func handleUnauthorized(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "application/json" {
		c.Set("error", "unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{})
	} else {
		c.Redirect(302, "/login")
	}
}
