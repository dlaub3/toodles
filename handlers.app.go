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
	if err := c.Bind(&user); err != nil {
		return
	}

	query := bson.M{"email": user.Email}
	existingUser := User{}
	if err := mongo.C(collectionToodlers).Find(query).One(&existingUser); err == nil {
		c.Set("httpStatus", 400)
		errors := make(map[string]string)
		errors["Email"] = "Please choose a different email."
		c.Set("error", errors)
		showSignupPage(c)
		return
	}

	user.Password, _ = crypt.HashPassword(user.Password, 32)
	if err := mongo.C(collectionToodlers).Insert(&user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errorInternalError).SetType(gin.ErrorTypePublic)
	}

	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "application/json" {
		c.JSON(http.StatusCreated, gin.H{})
	} else {
		c.Redirect(302, "/login")
	}
}

func logout(c *gin.Context) {
	invalidateCookies(c)
	c.Redirect(302, "/")
}
