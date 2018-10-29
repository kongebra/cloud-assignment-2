package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type WebhookDB struct {
	Addrs []string
	Database string
	Username string
	Password string
	Collection string
}

func (db *WebhookDB) Dial() (*mgo.Session, error) {
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

func (db *WebhookDB) Init() {
	session, err := db.Dial()

	defer session.Close()

	index := mgo.Index{
		Key:        []string{"id"},
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

func (db *WebhookDB) Insert(wh Webhook) (string, error) {
	session, err := db.Dial()

	defer session.Close()

	id := wh.ID

	err = session.DB(db.Database).C(db.Collection).Insert(wh)

	if err != nil {
		fmt.Printf("error in Insert(): %v", err.Error())
		return "", err
	}

	return id, nil
}

func (db *WebhookDB) Get(id string) (Webhook, bool) {
	session, err := db.Dial()

	defer session.Close()

	webhook := Webhook{}
	allWasGood := true

	err = session.DB(db.Database).C(db.Collection).Find(bson.M{"id": id}).One(&webhook)

	if err != nil {
		allWasGood = false
	}

	return webhook, allWasGood
}

func (db *WebhookDB) GetAll() []Webhook {
	session, _ := db.Dial()

	defer session.Close()

	var all []Webhook

	err := session.DB(db.Database).C(db.Collection).Find(bson.M{}).All(&all)

	if err != nil {
		return []Webhook{}
	}

	return all
}

func (db *WebhookDB) Count() int {
	session, err := db.Dial()

	defer session.Close()

	count, err := session.DB(db.Database).C(db.Collection).Count()

	if err != nil {
		fmt.Printf("error in Count(): %v", err.Error())
		return -1
	}

	return count
}

func (db *WebhookDB) Delete(id string) error {
	session, err := db.Dial()

	defer session.Close()

	err = session.DB(db.Database).C(db.Collection).Remove(bson.M{"id": id})

	return err
}