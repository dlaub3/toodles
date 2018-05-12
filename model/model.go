package model

import (
	"log"

	. "github.com/dlaub3/todo/config"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Mongo databse ORM
var Mongo *mgo.Database

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
	Mongo = session.DB(connect.Database)
}

// Todo model
type Todo struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Content string        `bson:"content" json:"content"`
}

// Note model
type Note struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Content string        `bson:"content" json:"content"`
}
