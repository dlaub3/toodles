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

	getUser(UID.(string), &user)

	render(c, gin.H{
		"user":  user,
		"title": "Accout Info",
	}, "user.html")
}

func deleteAUser(c *gin.Context) {
	UID, _ := c.Get("uid")
	deleteUser(UID.(string))
	deleteAllToodles(UID.(string))
	logout(c)
}

func updateAUser(c *gin.Context) {

}

func updateOrDeleteAUser(c *gin.Context) {
	method := c.PostForm("method")

	if method == "put" {

	} else if method == "delete" {
		deleteAUser(c)
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
