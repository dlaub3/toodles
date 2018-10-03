package main

import (
	"log"
	"net/http"

	"github.com/dlaub3/toodles/crypt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func getAUser(c *gin.Context) {
	user := User{}
	UID, _ := c.Get("uid")

	if err := getUserByUID(UID.(string), &user); err != nil {
		c.Set("genError", "😨 error getting account info")
		c.Set("httpStatus", http.StatusInternalServerError)
		log.Println("getUserByUID: " + UID.(string))
	}

	render(c, gin.H{
		"user":  user,
		"title": "Accout Info",
	}, "user.html")
}

func deleteAccount(c *gin.Context) {
	UID, _ := c.Get("uid")
	if err := deleteUser(UID.(string)); err != nil {
		c.Set("genError", "😨 error deleting account: "+UID.(string))
		c.Set("httpStatus", http.StatusInternalServerError)
		log.Println("deleteUser: " + UID.(string))
	}

	if err := deleteAllToodles(UID.(string)); err != nil {
		c.Set("genError", "😨 error deleting toodles: "+UID.(string))
		c.Set("httpStatus", http.StatusInternalServerError)
		log.Println("deleteAllToodles: " + UID.(string))
	}

	logout(c)
}

func updateOrDeleteAUser(c *gin.Context) {
	method := c.PostForm("method")

	if method == "put" {

	} else if method == "delete" {
		deleteAccount(c)
	}
}

func registerNewUser(c *gin.Context) {

	user := User{}
	user.ID = bson.NewObjectId()
	user.UID = user.ID.Hex()
	if err := c.ShouldBind(&user); err != nil {
		c.Set("error", getValidationErrorMsg(err))
		c.Set("httpStatus", http.StatusBadRequest)
		showSignupPage(c)
		return
	}

	query := bson.M{"email": user.Email}
	existingUser := User{}
	if err := mongo.C(collectionToodlers).Find(query).One(&existingUser); err == nil {
		c.Set("httpStatus", http.StatusBadRequest)
		errors := make(map[string]string)
		errors["Email"] = "Please choose a different email."
		c.Set("error", errors)
		showSignupPage(c)
		return
	}

	user.Password, _ = crypt.HashPassword(user.Password, 32)
	if err := mongo.C(collectionToodlers).Insert(&user); err != nil {
		c.Set("genError", "😨 failed to register account. Please try again.")
		c.Set("httpStatus", http.StatusInternalServerError)
		log.Println("registerNewUser: " + err.Error())
		log.Println("Params: UID=" + user.UID + " password=" + user.Password)
		showSignupPage(c)
		return
	}

	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "application/json" {
		c.JSON(http.StatusCreated, gin.H{})
	} else {
		c.Redirect(http.StatusFound, "/login")
	}
}
