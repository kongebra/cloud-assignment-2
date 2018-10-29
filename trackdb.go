package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type TrackDB struct {
	Addrs []string
	Database string
	Username string
	Password string
	Collection string
}

func (db *TrackDB) Dial() (*mgo.Session, error) {
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

func (db *TrackDB) Init() {
	session, err := db.Dial()

	defer session.Close()

	index := mgo.Index{
		Key:        []string{"track_src_url"},
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

func (db *TrackDB) Insert(t Track) (bson.ObjectId, error) {
	session, err := db.Dial()

	defer session.Close()

	id := bson.NewObjectId()
	t.Id = id

	err = session.DB(db.Database).C(db.Collection).Insert(t)

	if err != nil {
		fmt.Printf("error in Insert(): %v", err.Error())
		return "", err
	}

	return id, nil
}

func (db *TrackDB) Get(id string) (Track, bool) {
	session, err := db.Dial()

	defer session.Close()

	track := Track{}
	allWasGood := true

	err = session.DB(db.Database).C(db.Collection).Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&track)

	if err != nil {
		allWasGood = false
	}

	return track, allWasGood
}

func (db *TrackDB) Count() int {
	session, err := db.Dial()

	defer session.Close()

	count, err := session.DB(db.Database).C(db.Collection).Count()

	if err != nil {
		fmt.Printf("error in Count(): %v", err.Error())
		return -1
	}

	return count
}

func (db *TrackDB) GetAll() []Track {
	session, _ := db.Dial()

	defer session.Close()

	var all []Track

	err := session.DB(db.Database).C(db.Collection).Find(bson.M{}).All(&all)

	if err != nil {
		return []Track{}
	}

	return all
}

func (db *TrackDB) DeleteAll() error {
	session, _ := db.Dial()

	defer session.Close()

	_, err := session.DB(db.Database).C(db.Collection).RemoveAll(bson.M{})

	return err
}