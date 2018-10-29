package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Clock ticker go routine
// Takes in an int-pointer
func clockTicker(count *int) {
	// Checks if the DB-count has changed
	if *count != trackDB.Count() {
		// Get the current time in milliseconds
		startTime := time.Now().Unix()
		// Declare the endTime
		var endTime int64

		// Get old count
		oldCount := *count
		// Get new count
		newCount := trackDB.Count()
		// Set count to new count
		*count = trackDB.Count()

		// Declare a message
		var message string

		// Checks if we deleted all the tracks
		if newCount == 0 {
			// set end time
			endTime = time.Now().Unix()
			// get time elapsed
			elapsed := endTime - startTime
			// arrange the message
			message = fmt.Sprintf("Latest timestamp: %d, all (amount: %d) tracks deleted. (Processing: %dms)", time.Now().Unix(), oldCount, elapsed)
		} else {
			// get how many new tracks are added
			diff := newCount - oldCount
			// get all tracks
			tracks := trackDB.GetAll()
			// create empty bson.ObjectId slice
			ids := make([]bson.ObjectId, 0)

			// loop from old through all tracks
			for i := oldCount; i < len(tracks); i++ {
				// Get ID's from new tracks
				ids = append(ids, tracks[i].Id)
			}

			// get timestamp from the last track
			timestamp := tracks[len(tracks) - 1].Timestamp

			// arrange message
			message = fmt.Sprintf("Latest timestamp: %d, %d new tracks are: ", timestamp, diff)

			// loop through all ID's
			for _, value := range ids {
				// append id on message
				message += value.Hex() + ", "
			}

			// set end time
			endTime = time.Now().Unix()
			// get time elapsed
			elapsed := endTime - startTime
			// append time elapsed to message
			message += fmt.Sprintf("(Processing: %dms)", elapsed)
		}

		// Send message to Discord
		SendDiscordLogEntry(message)
	}
}

// Start the clock ticker go routine
func startClockTicker() {
	// Get the track-count
	var count = trackDB.Count()

	// Loopie loop
	for {
		// Set routine timer
		<- time.After(10 * time.Minute)
		// Call goroutine
		clockTicker(&count)
	}
}
