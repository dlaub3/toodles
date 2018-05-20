package main

import (
	jwt "github.com/dlaub3/gin-jwt"
	. "github.com/dlaub3/toodles/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func getAllTodos(c *gin.Context) {
	showAllToodles(c)
}

func getATodo(c *gin.Context) {
	id := c.Param("todo_id")
	userid := getUserIDFromJWT(c)
	toodleuser := getToodleUserID(id)
	if userid == toodleuser {
		todo := Todo{}
		todo.ID = bson.ObjectIdHex(id)
		Mongo.C(CollectionToodles).FindId(bson.ObjectIdHex(id)).One(&todo)
		render(c, gin.H{
			"title":   "Todo",
			"payload": todo}, "todo.html")
	} else {
		render(c, gin.H{
			"title":   "403 Can't touch this.",
			"payload": "403 Can't touch this."}, "error.html")
	}
}

func createATodo(c *gin.Context) {
	todo := Todo{}
	c.Bind(&todo)
	todo.ID = bson.NewObjectId()
	todo.UserID = getUserIDFromJWT(c)
	Mongo.C(CollectionToodles).Insert(&todo)
	showAToodle(c, todo)
}

func updateATodo(c *gin.Context) {
	id := c.Param("todo_id")
	todo := Todo{}
	c.Bind(&todo)
	todo.ID = bson.ObjectIdHex(id)
	userid := getUserIDFromJWT(c)
	toodleuser := getToodleUserID(id)
	if userid == toodleuser {
		Mongo.C(CollectionToodles).UpdateId(todo.ID, &todo)
		showAToodle(c, todo)
	} else {
		render(c, gin.H{
			"title":   "403 Can't touch this.",
			"payload": "403 Can't touch this."}, "error.html")
	}
}

func deleteATodo(c *gin.Context) {
	todo := Todo{}
	id := c.Param("todo_id")
	userid := getUserIDFromJWT(c)
	toodleuser := getToodleUserID(id)
	if userid == toodleuser {
		Mongo.C(CollectionToodles).RemoveId(bson.ObjectIdHex(id))
		showAToodle(c, todo)
	} else {
		render(c, gin.H{
			"title":   "403 Can't touch this.",
			"payload": "403 Can't touch this."}, "error.html")
	}
}

/*
	Support update/delete when JS is disabled by
	using a hidden form field.
*/

func updateOrDeleteTodo(c *gin.Context) {
	todo := Todo{}

	method := c.PostForm("method")
	id := c.Param("todo_id")

	userid := getUserIDFromJWT(c)
	toodleuser := getToodleUserID(id)
	if userid == toodleuser {
		if method == "put" {
			c.Bind(&todo)
			todo.ID = bson.ObjectIdHex(id)
			Mongo.C(CollectionToodles).UpdateId(bson.ObjectIdHex(id), &todo)
			render(c, gin.H{
				"title":   "Todo",
				"payload": todo}, "todo.html")
		} else if method == "delete" {
			Mongo.C(CollectionToodles).RemoveId(bson.ObjectIdHex(id))
			showAllToodles(c)
		}
	} else {
		render(c, gin.H{
			"title":   "403 Can't touch this.",
			"payload": "403 Can't touch this."}, "error.html")
	}
}

func showAToodle(c *gin.Context, todo Todo) {
	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "application/json" {
		render(c, gin.H{
			"title":   "Todo",
			"payload": todo}, "todo.html")
	} else {
		showAllToodles(c)
	}
}

func showAllToodles(c *gin.Context) {
	var todos []Todo
	claims := jwt.ExtractClaims(c)
	Mongo.C(CollectionToodles).Find(bson.M{"userid": claims["id"]}).All(&todos)
	render(c, gin.H{
		"title":   "Todo",
		"payload": todos}, "todos.html")
}

func getUserIDFromJWT(c *gin.Context) string {
	claims := jwt.ExtractClaims(c)
	userid := claims["id"].(string)
	return userid
}

func getToodleUserID(id string) string {
	todo := Todo{}
	Mongo.C(CollectionToodles).FindId(bson.ObjectIdHex(id)).One(&todo)
	return todo.UserID
}
