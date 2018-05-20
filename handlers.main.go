package main

import (
	. "github.com/dlaub3/toodles/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
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

func registerNewUser(c *gin.Context) {

	user := User{}
	user.ID = bson.NewObjectId()
	user.UID = user.ID.Hex()
	c.Bind(&user)
	Mongo.C(CollectionToodlers).Insert(&user)

	Toodles := Toodles{}
	Toodles.ID = user.ID
	Mongo.C(CollectionToodles).Insert(&Toodles)

	render(c, gin.H{
		"title": "Golang Todo Applicaiton"}, "login.html")
}
