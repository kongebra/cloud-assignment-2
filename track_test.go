package main

import (
	"testing"
	"time"
)

func TestTrackInit(t *testing.T) {
	var timestamp = time.Now().Unix()
	var pilot = "Roger Hansen"
	var glider = "Glider 3000"
	var gliderId = "G3000"
	var length = 3.141592
	var url = "http://google.com/q?=meaning of life"

	var track = Track{
		Timestamp: timestamp,
		Pilot: pilot,
		Glider: glider,
		GliderID: gliderId,
		TrackLength: length,
		TrackSrcUrl: url,
	}

	if track.Timestamp != timestamp {
		t.Error("Track Timestamp Error")
	}

	if track.Pilot != pilot {
		t.Error("Track Pilot Error")
	}

	if track.Glider != glider {
		t.Error("Track Glider Error")
	}

	if track.GliderID != gliderId {
		t.Error("Track GliderID Error")
	}

	if track.TrackLength != length {
		t.Error("Track Length Error")
	}

	if track.TrackSrcUrl != url {
		t.Error("Track Source URL Error")
	}
}

func TestTrack_GetField(t *testing.T) {

}