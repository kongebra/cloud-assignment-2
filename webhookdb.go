package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Webhook Database
type WebhookDB struct {
	Addrs []string
	Database string
	Username string
	Password string
	Collection string
}

// Dial up
func (db *WebhookDB) Dial() (*mgo.Session, error) {
	dialInfo := &mgo.DialInfo{
		Addrs: db.Addrs,
		Database: db.Database,
		Username: db.Username,
		Password: db.Password,
		Timeout: 60 * time.Second,
	}

	session, err := mgo.DialWithInfo(dialInfo)

	// check for errors
	if err != nil {
		panic(err)
	}

	return session, err
}

// Initialize database
func (db *WebhookDB) Init() {
	// dial
	session, err := db.Dial()

	// clean up
	defer session.Close()

	// set indexes
	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	// ensure indexes
	err = session.DB(db.Database).C(db.Collection).EnsureIndex(index)

	// check for errors
	if err != nil {
		panic(err)
	}
}

// Insert Webhook to DB
func (db *WebhookDB) Insert(wh Webhook) (string, error) {
	// dial
	session, err := db.Dial()

	// clean up
	defer session.Close()

	// get ID
	id := wh.ID

	// insert webhook
	err = session.DB(db.Database).C(db.Collection).Insert(wh)

	// check for errors
	if err != nil {
		fmt.Printf("error in Insert(): %v", err.Error())
		return "", err
	}

	// return id
	return id, nil
}

// Get webhook from DB
func (db *WebhookDB) Get(id string) (Webhook, bool) {
	// DIAL
	session, err := db.Dial()

	// CLEAN UP
	defer session.Close()

	// make Webhook variable
	webhook := Webhook{}
	// ALL GOOD
	allWasGood := true

	// Get webhook from DB
	err = session.DB(db.Database).C(db.Collection).Find(bson.M{"id": id}).One(&webhook)

	// Check for errors
	if err != nil {
		// NOT GOOD
		allWasGood = false
	}

	return webhook, allWasGood
}

// Get all webhooks from the DB
func (db *WebhookDB) GetAll() []Webhook {
	// dial
	session, _ := db.Dial()

	// clean
	defer session.Close()

	// webhooks
	var all []Webhook

	// get all webhooks
	err := session.DB(db.Database).C(db.Collection).Find(bson.M{}).All(&all)

	// error
	if err != nil {
		// return empty slice
		return []Webhook{}
	}

	// return all
	return all
}

// Count webhooks in the database
func (db *WebhookDB) Count() int {
	// Dial
	session, err := db.Dial()

	// Clean up
	defer session.Close()

	// get count
	count, err := session.DB(db.Database).C(db.Collection).Count()

	// ERROR
	if err != nil {
		fmt.Printf("error in Count(): %v", err.Error())
		return -1
	}

	return count
}

// Delete webhook from the database
func (db *WebhookDB) Delete(id string) error {
	// Dial
	session, err := db.Dial()

	// Clean up
	defer session.Close()

	// Remove from db
	err = session.DB(db.Database).C(db.Collection).Remove(bson.M{"id": id})

	return err
}