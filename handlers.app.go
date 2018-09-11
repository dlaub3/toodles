package main

import (
	"net/http"

	"github.com/dlaub3/toodles/crypt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
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

func registerNewUser(c *gin.Context) {

	user := User{}
	user.ID = bson.NewObjectId()
	user.UID = user.ID.Hex()
	c.Bind(&user)

	query := bson.M{"email": user.Email}
	existingUser := User{}
	Mongo.C(CollectionToodlers).Find(query).One(&existingUser)

	if existingUser.Email == user.Email {
		c.Set("httpStatus", 400)
		c.Set("error", "Please choose a different username.")
		showSignupPage(c)
		return
	}

	user.Password, _ = crypt.HashPassword(user.Password, 32)
	Mongo.C(CollectionToodlers).Insert(&user)

	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "application/json" {
		c.JSON(http.StatusCreated, gin.H{})
	} else {
		c.Redirect(302, "/login")
	}

}

func logout(c *gin.Context) {
	InvalidateCookies(c)
	c.Redirect(302, "/")
}
