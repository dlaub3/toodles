package main

import (
	. "github.com/dlaub3/toodles/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

const (
	collectionTodo = "toodles"
)

func getAllTodos(c *gin.Context) {
	var todos []Todo
	Mongo.C(collectionTodo).Find(bson.M{}).All(&todos)
	showAllToodles(c)
}

func getATodo(c *gin.Context) {
	id := c.Param("todo_id")
	todo := Todo{}
	Mongo.C(collectionTodo).FindId(bson.ObjectIdHex(id)).One(&todo)
	render(c, gin.H{
		"title":   "Todo",
		"payload": todo}, "todo.html")
}

func createATodo(c *gin.Context) {
	todo := Todo{}
	c.Bind(&todo)
	todo.ID = bson.NewObjectId()
	Mongo.C(collectionTodo).Insert(&todo)
	showAToodle(c, todo)
}

func updateATodo(c *gin.Context) {
	id := c.Param("todo_id")
	todo := Todo{}
	c.Bind(&todo)
	todo.ID = bson.ObjectIdHex(id)
	Mongo.C(collectionTodo).UpdateId(todo.ID, &todo)
	showAToodle(c, todo)
}

func deleteATodo(c *gin.Context) {
	todo := Todo{}
	id := c.Param("todo_id")
	Mongo.C(collectionTodo).RemoveId(bson.ObjectIdHex(id))
	showAToodle(c, todo)
}

/*
	Support update/delete when JS is disabled by
	using a hidden form field.
*/

func updateOrDeleteTodo(c *gin.Context) {
	todo := Todo{}

	method := c.PostForm("method")
	id := c.Param("todo_id")

	if method == "put" {
		c.Bind(&todo)
		todo.ID = bson.ObjectIdHex(id)
		Mongo.C(collectionTodo).UpdateId(bson.ObjectIdHex(id), &todo)
		render(c, gin.H{
			"title":   "Todo",
			"payload": todo}, "todo.html")
	} else if method == "delete" {
		Mongo.C(collectionTodo).RemoveId(bson.ObjectIdHex(id))
		showAllToodles(c)
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
	Mongo.C(collectionTodo).Find(bson.M{}).All(&todos)
	render(c, gin.H{
		"title":   "Todo",
		"payload": todos}, "todos.html")
}
