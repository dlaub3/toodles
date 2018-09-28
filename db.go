package main

import (
	"log"

	"github.com/globalsign/mgo"
)

var mongo *mgo.Database

const (
	// todos
	collectionToodles = "toodles"
	// users
	collectionToodlers = "toodlers"
)

func initDB() {

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
	mongo = session.DB(connect.Database)
}
