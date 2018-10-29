package main

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Ticker struct
type Ticker struct {
	TLatest int64 `json:"t_latest"`
	TStart int64 `json:"t_start"`
	TStop int64 `json:"t_stop"`
	Tracks []bson.ObjectId `json:"tracks"`
	Processing int64 `json:"processing"`
}

// Get current timestamp
func (t *Ticker) Timestamp() {
	t.TLatest = time.Now().Unix()
}