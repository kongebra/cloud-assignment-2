package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// TrackDB struct
type TrackDB struct {
	Addrs []string
	Database string
	Username string
	Password string
	Collection string
}

// Dial up the database
func (db *TrackDB) Dial() (*mgo.Session, error) {
	dialInfo := &mgo.DialInfo{
		Addrs: db.Addrs,
		Database: db.Database,
		Username: db.Username,
		Password: db.Password,
		Timeout: 60 * time.Second,
	}

	session, err := mgo.DialWithInfo(dialInfo)

	// Check for errors
	if err != nil {
		panic(err)
	}

	return session, err
}

// Initialize the DB
func (db *TrackDB) Init() {
	// Dial up
	session, err := db.Dial()

	// Close the session after we are done
	defer session.Close()

	// Set indexes for the DB
	index := mgo.Index{
		Key:        []string{"track_src_url"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	// Set the indexes
	err = session.DB(db.Database).C(db.Collection).EnsureIndex(index)

	// Check for errrors
	if err != nil {
		panic(err)
	}
}

// Insert track to the DB
func (db *TrackDB) Insert(t Track) (bson.ObjectId, error) {
	// Dial up
	session, err := db.Dial()

	// Tidy up after ourself
	defer session.Close()

	// Get new ObjectId
	id := bson.NewObjectId()
	// Set ID
	t.Id = id

	// Insert track to DB
	err = session.DB(db.Database).C(db.Collection).Insert(t)

	// Check for errors
	if err != nil {
		fmt.Printf("error in Insert(): %v", err.Error())
		return "", err
	}

	return id, nil
}

// Get track from DB based on ID
func (db *TrackDB) Get(id string) (Track, bool) {
	// Dial up
	session, err := db.Dial()

	// Clean up
	defer session.Close()

	// make empty track
	track := Track{}
	// All is GOOD
	allWasGood := true

	// get track and set data to track
	err = session.DB(db.Database).C(db.Collection).Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&track)

	// check for errors
	if err != nil {
		// all is not good anymore
		allWasGood = false
	}

	return track, allWasGood
}

// Get count of tracks in database
func (db *TrackDB) Count() int {
	// Dial up
	session, err := db.Dial()

	// Clean me :)
	defer session.Close()

	// Get count from database
	count, err := session.DB(db.Database).C(db.Collection).Count()

	// Check for errors
	if err != nil {
		fmt.Printf("error in Count(): %v", err.Error())
		return -1
	}

	return count
}

// Get all tracks from DB
func (db *TrackDB) GetAll() []Track {
	// RING RING
	session, _ := db.Dial()

	// GET OUT OFF HERE!
	defer session.Close()

	// Declare track-slice
	var all []Track

	// Get all tracks
	err := session.DB(db.Database).C(db.Collection).Find(bson.M{}).All(&all)

	// Check for errors
	if err != nil {
		return []Track{}
	}

	return all
}

// Delete all tracks
func (db *TrackDB) DeleteAll() error {
	// Hello?
	session, _ := db.Dial()

	// Bye!
	defer session.Close()

	// Delete all tracks
	_, err := session.DB(db.Database).C(db.Collection).RemoveAll(bson.M{})

	return err
}