package main

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

var db *mgo.Database

type DBServer struct {
	Server   string
	Database string
}

func (m *DBServer) Connect() {
	session, err := mgo.Dial(m.Server)

	if err != nil {
		log.Fatal(err)
	}
	if err := session.Ping(); err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

func (m *DBServer) FindAll() ([]People, error) {
	var p []People
	err := db.C("peoples").Find(bson.M{}).All(&p)
	return p, err
}

func (m *DBServer) FindById(id string) (People, error) {
	var p People
	err := db.C("peoples").FindId(bson.ObjectIdHex(id)).One(&p)
	return p, err
}

func (m *DBServer) Insert(p People) error {
	return db.C("peoples").Insert(&p)
}

func (m *DBServer) Delete(p People) error {
	return db.C("peoples").Remove(&p)
}

func (m *DBServer) Update(p People) error {
	return db.C("peoples").UpdateId(p.ID, &p)
}
