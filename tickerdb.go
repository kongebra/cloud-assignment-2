package main

import (
	"gopkg.in/mgo.v2"
	"time"
)

type TickerDB struct {
	Addrs []string
	Database string
	Username string
	Password string
	Collection string
}

func (db *TickerDB) Dial() (*mgo.Session, error) {
	dialInfo := &mgo.DialInfo{
		Addrs: db.Addrs,
		Database: db.Database,
		Username: db.Username,
		Password: db.Password,
		Timeout: 60 * time.Second,
	}

	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		panic(err)
	}

	return session, err
}

func (db *TickerDB) Init() {
	session, err := db.Dial()

	defer session.Close()

	index := mgo.Index{
		Key:        []string{"foo"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = session.DB(db.Database).C(db.Collection).EnsureIndex(index)

	if err != nil {
		panic(err)
	}
}

