package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func getAllToodles(c *gin.Context) {
	showAllToodles(c)
}

func getAToodle(c *gin.Context) {
	c.Set("showsingle", true)

	UID, _ := c.Get("uid")
	id := c.Param("toodle_id")

	if id == "" {
		c.Set("genError", "😨 cannot find toodle. Please try again.")
		c.Set("httpStatus", http.StatusBadRequest)
		log.Println("getAToodle: missing id")
		showAllToodles(c)
		return
	}

	if success, toodle := findToodle(UID.(string), id); success == true {
		showAToodle(c, toodle)
		return
	}

	c.Set("genError", "😨 cannot find toodle: "+id)
	c.Set("httpStatus", http.StatusInternalServerError)
	log.Println("getAToodle: invalid id: " + id)
	showAllToodles(c)
}

func createAToodle(c *gin.Context) {

	toodle := Toodle{}
	toodle.ID = bson.NewObjectId()
	UID, _ := c.Get("uid")

	if err := c.ShouldBind(&toodle); err != nil {
		c.Set("error", getValidationErrorMsg(err))
		c.Set("httpStatus", http.StatusBadRequest)
		log.Println("validation error :" + err.Error())
		showAToodle(c, toodle)
		return
	}

	if _, err := createToodle(UID.(string), &toodle); err != nil {
		c.Set("genError", "😨 failed to create toodle. Please try again.")
		c.Set("httpStatus", http.StatusInternalServerError)
		log.Println("createAToodle: " + err.Error())
		log.Println("Params: UID=" + UID.(string))
	}

	showAToodle(c, toodle)
}

func updateAToodle(c *gin.Context) {

	toodle := Toodle{}
	id := c.Param("toodle_id")
	UID, _ := c.Get("uid")
	toodle.ID = bson.ObjectIdHex(id)

	if id == "" {
		c.Set("genError", "😨 cannot find toodle. Please try again.")
		c.Set("httpStatus", http.StatusBadRequest)
		log.Println("updateAToodle: missing id")
		showAToodle(c, toodle)
		return
	}

	if err := c.ShouldBind(&toodle); err != nil {
		c.Set("error", getValidationErrorMsg(err))
		c.Set("httpStatus", http.StatusBadRequest)
		showAToodle(c, toodle)
		return
	}

	if err := updateToodle(UID.(string), id, &toodle); err != nil {
		c.Set("genError", "😨 failed to update toodle. Please try again.")
		c.Set("httpStatus", http.StatusInternalServerError)
		log.Println("updateAToodle: " + err.Error())
		log.Println("Params: id=" + id + " UID=" + UID.(string))
	}

	showAToodle(c, toodle)
}

func findToodle(UID string, id string) (bool, Toodle) {
	success := false

	toodles := Toodles{}
	toodle := Toodle{}
	toodle.ID = bson.ObjectIdHex(id)

	if err := getToodles(UID, &toodles); err != nil {
		log.Println("findToodle: " + err.Error())
		log.Println("Params: id=" + id + " UID=" + UID)
	}

	for _, item := range toodles.Toodles {
		if item.ID == toodle.ID {
			toodle.Title = item.Title
			toodle.Content = item.Content
			success = true
		}
	}
	return success, toodle
}

func completeAToodle(c *gin.Context) {

	id := c.Param("toodle_id")
	UID, _ := c.Get("uid")

	if id == "" {
		c.Set("genError", "😨 cannot find toodle. Please try again.")
		c.Set("httpStatus", http.StatusBadRequest)
		log.Println("completeAToodle: missing id")
		showAllToodles(c)
		return
	}

	if err := completeToodle(UID.(string), id); err != nil {
		c.Set("genError", "😨 failed to update toodle. Please try again.")
		c.Set("httpStatus", http.StatusInternalServerError)
		log.Println("completeAToodle: " + err.Error())
		log.Println("Params: id=" + id + " UID=" + UID.(string))
	}

	showAllToodles(c)
}

func deleteAToodle(c *gin.Context) {

	id := c.Param("toodle_id")
	UID, _ := c.Get("uid")

	if id == "" {
		c.Set("genError", "😨 cannot find toodle. Please try again.")
		c.Set("httpStatus", http.StatusBadRequest)
		log.Println("updateAToodle: missing id")
		showAllToodles(c)
		return
	}

	if _, err := deleteToodle(UID.(string), id); err != nil {
		c.Set("genError", "😨 failed to delete toodle. Please try again.")
		c.Set("httpStatus", http.StatusInternalServerError)
		log.Println("deleteAToodle: " + err.Error())
		log.Println("Params: id=" + id + " UID=" + UID.(string))
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
		c.Set("showsingle", true)
		updateAToodle(c)
	} else if method == "delete" {
		deleteAToodle(c)
	} else {
		c.AbortWithError(http.StatusInternalServerError, errorInternalError).SetType(gin.ErrorTypePublic)
		log.Println("coupdateOrDeleteAToodle: invalid HTTP method:" + method)
		return
	}
}

func showAToodle(c *gin.Context, toodle Toodle) {
	contentType := c.Request.Header.Get("Content-Type")
	showSingle, exists := c.Get("showsingle")
	if exists == false {
		showSingle = false
	}

	if contentType == "application/json" || showSingle.(bool) == true {
		render(c, gin.H{
			"title":  "Toodle",
			"toodle": toodle}, "toodle.html")
	} else {
		showAllToodles(c)
	}
}

func showAllToodles(c *gin.Context) {
	toodles := Toodles{}
	UID, _ := c.Get("uid")
	getToodles(UID.(string), &toodles)

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
		"toodles":   activeToodles.Toodles,
	}, "toodles.html")
}
