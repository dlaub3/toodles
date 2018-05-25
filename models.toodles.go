package main

import "github.com/globalsign/mgo/bson"

// Toodle model
type Toodle struct {
	ID      bson.ObjectId `bson:"_id" form:"id" json:"id" `
	Title   string        `bson:"title" form:"title" json:"title" binding:"required"`
	Content string        `bson:"content" form:"content" json:"content"`
}

// Toodles model
type Toodles struct {
	ID      bson.ObjectId `bson:"_id" form:"id" json:"id" `
	Toodles []Toodle      `bson:"toodles"`
}
