package main

import (
	"net/http"
	"time"

	"github.com/dlaub3/toodles/crypt"
	"github.com/gin-gonic/gin"
)

func initializeRoutes() {

	auth := r.Group("/")
	authMiddleware := jwtMiddleware()
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/toodles", getAllToodles)
		auth.POST("/toodles", createAToodle)
		auth.GET("/toodles/:toodle_id", getAToodle)
		auth.PUT("/toodles/:toodle_id", updateAToodle)
		auth.DELETE("/toodles/:toodle_id", deleteAToodle)
		auth.PUT("/toodles/:toodle_id/complete", completeToodle)

		//Routes specifically for form submitalls and not AJAX
		auth.POST("/toodles/:toodle_id", updateOrDeleteToodle)
		auth.POST("/toodles/:toodle_id/complete", completeToodle)

		auth.GET("refresh_token", authMiddleware.RefreshHandler)
	}

	r.GET("/", showHomePage)
	r.GET("/signup", showSignupPage)
	r.POST("/signup", registerNewUser)
	r.GET("/login", showLoginPage)
	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/logout", logout)

}

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
