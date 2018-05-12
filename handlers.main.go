package main

import (
	"github.com/gin-gonic/gin"
)

func showHomePage(c *gin.Context) {
	render(c, gin.H{
		"title": "Golang Todo Applicaiton"}, "index.html")
}

func showLoginPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Golang Todo Applicaiton"}, "login.html")
}
func showSignupPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Golang Todo Applicaiton"}, "signup.html")
}
