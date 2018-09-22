package main

import (
	"github.com/globalsign/mgo/bson"
	mgo "gopkg.in/mgo.v2"
)

// Toodle model
type Toodle struct {
	ID      bson.ObjectId `bson:"_id" form:"id" json:"id" `
	Title   string        `bson:"title" form:"title" json:"title" binding:"required,max=150"`
	Content string        `bson:"content" form:"content" json:"content" binding:"max=2000"`
	Status  string        `bson:"status" form:"status" json:"status"`
}

// Toodles model
type Toodles struct {
	ID      bson.ObjectId `bson:"_id" form:"id" json:"id" `
	Toodles []Toodle      `bson:"toodles"`
}

func getToodles(UID string, toodles *Toodles) error {
	return mongo.C(collectionToodles).FindId(bson.ObjectIdHex(UID)).One(&toodles)
}

func createToodle(UID string, toodle *Toodle) (*mgo.ChangeInfo, error) {
	query := bson.M{"_id": bson.ObjectIdHex(UID)}
	update := bson.M{"$push": bson.M{"toodles": &toodle}}
	return mongo.C(collectionToodles).Upsert(query, update)
}

func completeToodle(UID string, id string) error {
	query := bson.M{"_id": bson.ObjectIdHex(UID), "toodles._id": bson.ObjectIdHex(id)}
	update := bson.M{"$set": bson.M{"toodles.$.status": "complete"}}
	return mongo.C(collectionToodles).Update(query, update)
}

func updateToodle(UID string, id string, toodle *Toodle) error {
	query := bson.M{"_id": bson.ObjectIdHex(UID), "toodles._id": bson.ObjectIdHex(id)}
	update := bson.M{"$set": bson.M{"toodles.$": &toodle}}
	return mongo.C(collectionToodles).Update(query, update)
}

func deleteToodle(UID string, id string) (*mgo.ChangeInfo, error) {
	query := bson.M{"_id": bson.ObjectIdHex(UID)}
	update := bson.M{"$pull": bson.M{"toodles": bson.M{"_id": bson.ObjectIdHex(id)}}}
	return mongo.C(collectionToodles).Upsert(query, update)
}
