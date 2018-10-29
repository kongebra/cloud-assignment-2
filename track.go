package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

// Track struct
type Track struct {
	Id 			bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Timestamp	int64 `json:"-" bson:"timestamp,omitempty"`

	HDate 		string `json:"H_date" bson:"H_date"`
	Pilot 		string `json:"pilot" bson:"pilot"`
	Glider 		string `json:"glider" bson:"glider"`
	GliderID 	string `json:"glider_id" bson:"glider_id"`
	TrackLength float64 `json:"track_length" bson:"track_length"`
	TrackSrcUrl string `json:"track_src_url" bson:"track_src_url"`

}

// Get a specific field from the Track as a string
func (track *Track) GetField(field string) string {
	// Response string
	var response string

	// Check if the field-parameter is valid
	switch field {
	case "pilot":
		response = track.Pilot
		break
	case "glider":
		response = track.Glider
		break
	case "glider_id":
		response = track.GliderID
		break
	case "track_length":
		response = fmt.Sprintf("%f", track.TrackLength)
		break
	case "H_date":
		response = track.HDate
		break
	case "track_src_url":
		response = track.TrackSrcUrl
		break
	default:
		response = ""
		break
	}

	return response
}