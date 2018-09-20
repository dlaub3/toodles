package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

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
	if err := mongo.C(collectionToodles).Find(query).One(&toodles); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errorInternalError).SetType(gin.ErrorTypePublic)
		return
	}

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
	UID := c.Keys["uid"].(string)
	if err := c.Bind(&toodle); err != nil {
		return
	}

	query := bson.M{"_id": bson.ObjectIdHex(UID)}
	update := bson.M{"$push": bson.M{"toodles": &toodle}}
	if _, err := mongo.C(collectionToodles).Upsert(query, update); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errorInternalError).SetType(gin.ErrorTypePublic)
		return
	}

	showAToodle(c, toodle)
}

func updateAToodle(c *gin.Context) {

	id := c.Param("toodle_id")
	toodle := Toodle{}
	if err := c.Bind(&toodle); err != nil {
		return
	}
	UID := c.Keys["uid"].(string)
	toodle.ID = bson.ObjectIdHex(id)
	query := bson.M{"_id": bson.ObjectIdHex(UID), "toodles._id": bson.ObjectIdHex(id)}
	update := bson.M{"$set": bson.M{"toodles.$": &toodle}}
	if err := mongo.C(collectionToodles).Update(query, update); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errorInternalError).SetType(gin.ErrorTypePublic)
		return
	}
	showAToodle(c, toodle)
}

func completeAToodle(c *gin.Context) {

	id := c.Param("toodle_id")
	UID := c.Keys["uid"].(string)

	query := bson.M{"_id": bson.ObjectIdHex(UID), "toodles._id": bson.ObjectIdHex(id)}
	update := bson.M{"$set": bson.M{"toodles.$.status": "complete"}}
	if err := mongo.C(collectionToodles).Update(query, update); err == nil {
		c.AbortWithError(http.StatusInternalServerError, errorInternalError).SetType(gin.ErrorTypePublic)
		return
	}

	toodles := Toodles{}
	toodle := Toodle{}
	query = bson.M{"_id": bson.ObjectIdHex(UID)}
	if err := mongo.C(collectionToodles).Find(query).One(&toodles); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errorInternalError).SetType(gin.ErrorTypePublic)
		return
	}

	for _, item := range toodles.Toodles {
		if item.ID == bson.ObjectIdHex(id) {
			toodle.ID = item.ID
			toodle.Status = item.Status
			toodle.Title = item.Title
			toodle.Content = item.Content
		}
	}

	showAToodle(c, toodle)
}

func deleteAToodle(c *gin.Context) {

	id := c.Param("toodle_id")
	UID := c.Keys["uid"].(string)
	query := bson.M{"_id": bson.ObjectIdHex(UID)}
	update := bson.M{"$pull": bson.M{"toodles": bson.M{"_id": bson.ObjectIdHex(id)}}}
	if _, err := mongo.C(collectionToodles).Upsert(query, update); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errorInternalError).SetType(gin.ErrorTypePublic)
		return
	}
	showAllToodles(c)
}

/*
	Support update/delete when JS is disabled by
	using a hidden form field.
*/

func updateOrDeleteAToodle(c *gin.Context) {

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
			"title":   "Toodle",
			"payload": toodle}, "toodle.html")
	} else {
		showAllToodles(c)
	}
}

func showAllToodles(c *gin.Context) {
	toodles := Toodles{}
	UID := c.Keys["uid"].(string)
	if err := mongo.C(collectionToodles).FindId(bson.ObjectIdHex(UID)).One(&toodles); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errorInternalError).SetType(gin.ErrorTypePublic)
		return
	}

	activeToodles := Toodles{}
	active := 0
	completed := 0
	for _, toodle := range toodles.Toodles {
		if toodle.Status != "complete" {
			active++
			activeToodles.Toodles = append(activeToodles.Toodles, toodle)
		} else {
			completed++
		}
	}

	render(c, gin.H{
		"active":    active,
		"completed": completed,
		"title":     "All your Toodles",
		"payload":   activeToodles.Toodles,
	}, "toodles.html")
}
