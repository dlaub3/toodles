package main

import "github.com/globalsign/mgo/bson"

// User model
type User struct {
	ID       bson.ObjectId `bson:"_id" `
	Email    string        `bson:"email" form:"email" json:"email" binding:"required"`
	Password string        `bson:"password" form:"password" json:"password"`
	UID      string        `bson:"uid"`
}
