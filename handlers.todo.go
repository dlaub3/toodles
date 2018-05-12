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
	render(c, gin.H{
		"title":   "Todo",
		"payload": todos}, "todos.html")
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
	render(c, gin.H{
		"title":   "Todo",
		"payload": todo}, "todo.html")
}

func updateATodo(c *gin.Context) {
	todo := Todo{}
	c.Bind(&todo)
	Mongo.C(collectionTodo).UpdateId(todo.ID, &todo)
	render(c, gin.H{
		"title":   "Todo",
		"payload": todo}, "todo.html")
}

func deleteATodo(c *gin.Context) {
	todo := Todo{}
	id := c.Param("todo_id")
	Mongo.C(collectionTodo).RemoveId(bson.ObjectIdHex(id))
	render(c, gin.H{
		"title":   "Todo",
		"payload": todo}, "todo.html")
}
