package main

import "github.com/globalsign/mgo/bson"

// Toodle model
type Toodle struct {
	ID      bson.ObjectId `bson:"_id" form:"id" json:"id" `
	Title   string        `bson:"title" form:"title" json:"title"`
	Content string        `bson:"content" form:"content" json:"content"`
	Status  string        `bson:"status" form:"status" json:"status"`
}

// Toodles model
type Toodles struct {
	ID      bson.ObjectId `bson:"_id" form:"id" json:"id" `
	Toodles []Toodle      `bson:"toodles"`
}
