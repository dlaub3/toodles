package model

import (
	"log"

	. "github.com/dlaub3/toodles/config"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// mongo databse ORM
var mongo *mgo.Database

type connection struct {
	Server   string
	Database string
}

// DB handle
var connect = connection{}
var config = Config{}

func init() {
	config.Read()
	connect.Server = config.Server
	connect.Database = config.Database
	session, err := mgo.Dial(connect.Server)
	if err != nil {
		log.Fatal(err)
	}
	mongo = session.DB(connect.Database)
}

const (
	// collectionToodles contains toodles
	collectionToodles = "toodles"
	// collectionToodlers contains toodlers info
	collectionToodlers = "toodlers"
)

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

// User model
type User struct {
	ID       bson.ObjectId `bson:"_id" `
	Email    string        `bson:"email" form:"email" json:"email" binding:"required"`
	Password string        `bson:"password" form:"password" json:"password"`
	UID      string        `bson:"uid"`
}
