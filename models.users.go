package main

import (
	"github.com/globalsign/mgo/bson"
)

// User model
type User struct {
	ID       bson.ObjectId `bson:"_id" `
	Email    string        `bson:"email" form:"email" json:"email" binding:"required,min=4,max=25"`
	Password string        `bson:"password" form:"password" json:"password" binding:"required,alphanum,min=4,max=30"`
	UID      string        `bson:"uid"`
}

func createUser(user User) error {
	return mongo.C(collectionToodlers).Insert(&user)
}

func updateUser(UID string, user User) error {
	query := bson.M{"uid": UID}
	update := bson.M{"$set": &user}
	return mongo.C(collectionToodlers).Update(query, update)
}

func deleteUser(UID string) error {
	query := bson.M{"uid": UID}
	return mongo.C(collectionToodlers).Remove(query)
}

func getUser(UID string, user *User) error {
	query := bson.M{"uid": UID}
	return mongo.C(collectionToodlers).Find(query).One(&user)
}
