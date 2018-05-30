package main

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

// CsrfToken binds with form submit csrf
type CsrfToken struct {
	CsrfToken string `form:"csrf" json:"csrf"`
}

func getAllToodles(c *gin.Context) {
	showAllToodles(c)
}

func getAToodle(c *gin.Context) {
	c.Keys["showsingle"] = true
	id := bson.ObjectIdHex(c.Param("toodle_id"))
	toodles := Toodles{}
	toodle := Toodle{}
	UID := c.Keys["uid"].(string)
	query := bson.M{"_id": bson.ObjectIdHex(UID)}
	Mongo.C(CollectionToodles).Find(query).One(&toodles)

	for _, item := range toodles.Toodles {
		if item.ID == id {
			toodle.ID = item.ID
			toodle.Title = item.Title
			toodle.Content = item.Content
		}
	}

	showAToodle(c, toodle)
}

func createAToodle(c *gin.Context) {

	validRequest := IsCSRFTokenValid(c)
	if !validRequest {
		returnError(c)
		return
	}

	toodle := Toodle{}
	toodle.ID = bson.NewObjectId()
	UID := c.Keys["uid"].(string)
	c.Bind(&toodle)

	query := bson.M{"_id": bson.ObjectIdHex(UID)}
	update := bson.M{"$push": bson.M{"toodles": &toodle}}
	Mongo.C(CollectionToodles).Upsert(query, update)

	showAToodle(c, toodle)
}

func updateAToodle(c *gin.Context) {
	validRequest := IsCSRFTokenValid(c)
	if !validRequest {
		returnError(c)
		return
	}

	id := c.Param("toodle_id")
	toodle := Toodle{}
	c.Bind(&toodle)
	UID := c.Keys["uid"].(string)
	toodle.ID = bson.ObjectIdHex(id)
	query := bson.M{"_id": bson.ObjectIdHex(UID), "toodles._id": bson.ObjectIdHex(id)}
	update := bson.M{"$set": bson.M{"toodles.$": &toodle}}
	Mongo.C(CollectionToodles).Update(query, update)
	showAToodle(c, toodle)
}

func deleteAToodle(c *gin.Context) {

	validRequest := IsCSRFTokenValid(c)
	if !validRequest {
		returnError(c)
		return
	}

	id := c.Param("toodle_id")
	UID := c.Keys["uid"].(string)
	query := bson.M{"_id": bson.ObjectIdHex(UID)}
	update := bson.M{"$pull": bson.M{"toodles": bson.M{"_id": bson.ObjectIdHex(id)}}}
	Mongo.C(CollectionToodles).Upsert(query, update)
	showAllToodles(c)
}

/*
	Support update/delete when JS is disabled by
	using a hidden form field.
*/

func updateOrDeleteToodle(c *gin.Context) {

	method := c.PostForm("method")

	if method == "put" {
		c.Keys["showsingle"] = true
		updateAToodle(c)
	} else if method == "delete" {
		deleteAToodle(c)
	}
}
func showAToodle(c *gin.Context, toodle Toodle) {
	contentType := c.Request.Header.Get("Content-Type")
	showSingle := c.Keys["showsingle"]
	if contentType == "application/json" || showSingle == true {
		render(c, gin.H{
			"title":     "Toodle",
			"csrfToken": c.Keys["csrftoken"],
			"payload":   toodle}, "toodle.html")
	} else {
		showAllToodles(c)
	}
}

func showAllToodles(c *gin.Context) {
	toodles := Toodles{}
	UID := c.Keys["uid"].(string)
	Mongo.C(CollectionToodles).FindId(bson.ObjectIdHex(UID)).One(&toodles)
	render(c, gin.H{
		"title":     "All your Toodles",
		"csrfToken": c.Keys["csrftoken"],
		"payload":   toodles.Toodles}, "toodles.html")
}

func returnError(c *gin.Context) {
	render(c, gin.H{
		"payload": "Our servers are busy please stand bye.",
	}, "error.html")
}

func IsCSRFTokenValid(c *gin.Context) bool {
	csrfToken := CsrfToken{}
	// save the request body
	body, _ := ioutil.ReadAll(c.Request.Body)
	// restore the request body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	// bind the token
	c.Bind(&csrfToken)
	// restore the request body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	// return token validation
	return csrfToken.CsrfToken == c.Keys["csrftoken"].(string)
}
