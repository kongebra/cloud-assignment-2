package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func clockTicker(count *int) {
	if *count != trackDB.Count() {
		startTime := time.Now().Unix()
		var endTime int64

		oldCount := *count
		newCount := trackDB.Count()
		*count = trackDB.Count()

		var message string

		if newCount == 0 {
			endTime = time.Now().Unix()
			elapsed := endTime - startTime
			message = fmt.Sprintf("Latest timestamp: %d, all (amount: %d) tracks deleted. (Processing: %dms)", time.Now().Unix(), oldCount, elapsed)
		} else {
			diff := newCount - oldCount
			tracks := trackDB.GetAll()

			ids := make([]bson.ObjectId, 0)

			for i := oldCount; i < len(tracks); i++ {
				ids = append(ids, tracks[i].Id)
			}

			timestamp := tracks[len(tracks) - 1].Timestamp

			message = fmt.Sprintf("Latest timestamp: %d, %d new tracks are: ", timestamp, diff)

			for _, value := range ids {
				message += value.Hex() + ", "
			}

			endTime = time.Now().Unix()
			elapsed := endTime - startTime
			message += fmt.Sprintf("(Processing: %dms)", elapsed)
		}

		SendDiscordLogEntry(message)
	}
}

func startClockTicker() {
	var count = trackDB.Count()

	for {
		<- time.After(10 * time.Minute)
		clockTicker(&count)
	}
}
