package main

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

// Mongo databse ORM
var Mongo *mgo.Database

const (
	// CollectionToodles is todos
	CollectionToodles = "toodles"
	// CollectionToodlers is users
	CollectionToodlers = "toodlers"
)

func dbConnect() {

	type connection struct {
		Server   string
		Database string
		Mongo    *mgo.Database
	}

	var connect = connection{}
	connect.Server = config.Server
	connect.Database = config.Database
	session, err := mgo.Dial(connect.Server)
	if err != nil {
		log.Fatal(err)
	}
	Mongo = session.DB(connect.Database)
}
