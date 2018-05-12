package main

import (
	. "github.com/dlaub3/toodles/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func getAllTodos(c *gin.Context) {
	var todos []Todo
	Mongo.C("todo").Find(bson.M{}).All(&todos)
	render(c, gin.H{
		"title":   "Todo",
		"payload": todos}, "todos.html")
}

func getATodo(c *gin.Context) {
	id := c.Param("todo_id")
	todo := Todo{}
	Mongo.C("todo").FindId(bson.ObjectIdHex(id)).One(&todo)
	render(c, gin.H{
		"title":   "Todo",
		"payload": todo}, "todo.html")
}

func createATodo(c *gin.Context) {
	todo := Todo{}
	c.BindJSON(&todo)
	todo.ID = bson.NewObjectId()
	Mongo.C("todo").Insert(&todo)
	render(c, gin.H{
		"title":   "Todo",
		"payload": todo}, "todo.html")
}

func updateATodo(c *gin.Context) {
	todo := Todo{}
	c.BindJSON(&todo)
	Mongo.C("todo").UpdateId(todo.ID, &todo)
	render(c, gin.H{
		"title":   "Todo",
		"payload": todo}, "todo.html")
}

func deleteATodo(c *gin.Context) {
	todo := Todo{}
	c.BindJSON(&todo)
	Mongo.C("todo").Remove(&todo)
	render(c, gin.H{
		"title":   "Todo",
		"payload": todo}, "todo.html")
}
