package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func showHomePage(c *gin.Context) {
	render(c, gin.H{
		"title":    "Toodles, stay organized. Get stuff done!",
		"subtitle": "Login or create an account to get started.",
	}, "index.html")
}

func showLoginPage(c *gin.Context) {
	render(c, gin.H{
		"title":    "Login to access your toodles.",
		"subtitle": "Signup to create an account.",
	}, "login.html")
}

func showSignupPage(c *gin.Context) {
	render(c, gin.H{
		"title":    "Signup to start toodling today.",
		"subtitle": "Complete the form below to create your account.",
	}, "signup.html")
}

func logout(c *gin.Context) {
	invalidateCookies(c)
	c.Redirect(http.StatusFound, "/")
}
