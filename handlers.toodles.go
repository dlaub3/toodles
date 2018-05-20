package main

import (
	. "github.com/dlaub3/toodles/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

// Used as a flag to indicate the desited response is a single toodle
var showSingle bool

func getAllToodles(c *gin.Context) {
	showAllToodles(c)
}

func getAToodle(c *gin.Context) {
	showSingle = true
	id := bson.ObjectIdHex(c.Param("toodle_id"))
	toodles := Toodles{}
	toodle := Toodle{}
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
	toodle := Toodle{}
	toodle.ID = bson.NewObjectId()
	c.Bind(&toodle)
	query := bson.M{"_id": bson.ObjectIdHex(UID)}
	update := bson.M{"$push": bson.M{"toodles": &toodle}}
	Mongo.C(CollectionToodles).Upsert(query, update)
	showAToodle(c, toodle)
}

func updateAToodle(c *gin.Context) {
	id := c.Param("toodle_id")
	toodle := Toodle{}
	c.Bind(&toodle)
	toodle.ID = bson.ObjectIdHex(id)
	query := bson.M{"_id": bson.ObjectIdHex(UID), "toodles._id": bson.ObjectIdHex(id)}
	update := bson.M{"$set": bson.M{"toodles.$": &toodle}}
	Mongo.C(CollectionToodles).Update(query, update)
	showAToodle(c, toodle)
}

func deleteAToodle(c *gin.Context) {
	id := c.Param("toodle_id")
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
		showSingle = true
		updateAToodle(c)
	} else if method == "delete" {
		deleteAToodle(c)
	}
}

func showAToodle(c *gin.Context, toodle Toodle) {
	contentType := c.Request.Header.Get("Content-Type")

	if contentType == "application/json" || showSingle == true {
		render(c, gin.H{
			"title":   "Toodle",
			"payload": toodle}, "toodle.html")
	} else {
		showAllToodles(c)
	}
}

func showAllToodles(c *gin.Context) {
	toodles := Toodles{}

	Mongo.C(CollectionToodles).FindId(bson.ObjectIdHex(UID)).One(&toodles)
	render(c, gin.H{
		"title":   "All your Toodles",
		"payload": toodles.Toodles}, "toodles.html")
}
