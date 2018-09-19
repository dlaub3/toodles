package main

import "github.com/globalsign/mgo/bson"

// User model
type User struct {
	ID       bson.ObjectId `bson:"_id" `
	Email    string        `bson:"email" form:"email" json:"email" binding:"required,min=4,max=25"`
	Password string        `bson:"password" form:"password" json:"password" binding:"required,alphanum,min=4,max=30"`
	UID      string        `bson:"uid"`
}
