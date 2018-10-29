package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/marni/goigc"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Start time of the application
var startTime time.Time

// Track database
var trackDB TrackDB

// Webhook database
var webhookDB WebhookDB

// Main function
func main() {
	// Start the go routine of the clock ticker
	go startClockTicker()

	// Track Database Settings
	trackDB = TrackDB{
		Addrs: []string{"ds133533.mlab.com:33533"},
		Database: "assignment-2",
		Username: "golang",
		Password: "golang1",
		Collection: "tracks",
	}

	// Initialize the Track Collection
	trackDB.Init()

	// Webhook Database Settings
	webhookDB = WebhookDB{
		Addrs: []string{"ds133533.mlab.com:33533"},
		Database: "assignment-2",
		Username: "golang",
		Password: "golang1",
		Collection: "webhooks",
	}

	// Initialize the Webhook Collection
	webhookDB.Init()

	// Record the start time for the application
	startTime = time.Now()

	// Create a mux-router
	router := mux.NewRouter().StrictSlash(true)

	// Redirects all request to URL/paragliding/ to URL/paragliding/api/
	router.HandleFunc("/paragliding/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/paragliding/api/", http.StatusMovedPermanently)
	})

	// Handle all routes for the application and send them to their respective functions
	router.HandleFunc("/paragliding/api/", APIIndex).Methods("GET")
	router.HandleFunc("/paragliding/api/track/", TrackPOST).Methods("POST")
	router.HandleFunc("/paragliding/api/track/", TrackGET).Methods("GET")
	router.HandleFunc("/paragliding/api/track/{id}/", SingleTrackGET).Methods("GET")
	router.HandleFunc("/paragliding/api/track/{id}/{field}/", SingleTrackFieldGET).Methods("GET")
	router.HandleFunc("/paragliding/api/ticker/", GetTicker).Methods("GET")
	router.HandleFunc("/paragliding/api/ticker/latest/", GetLatestTicker).Methods("GET")
	router.HandleFunc("/paragliding/api/ticker/{timestamp}/", GetTickerFromTimestamp).Methods("GET")
	router.HandleFunc("/paragliding/api/webhook/new_track/", WebhookNewTrack).Methods("POST")
	router.HandleFunc("/paragliding/api/webhook/new_track/{webhook_id}/", WebhookNewTrackIdGET).Methods("GET")
	router.HandleFunc("/paragliding/api/webhook/new_track/{webhook_id}/", WebhookNewTrackIdDELETE).Methods("DELETE")
	router.HandleFunc("/paragliding/admin/api/track_count/", AdminTrackCount).Methods("GET")
	router.HandleFunc("/paragliding/admin/api/tracks/", AdminDeleteTracks).Methods("DELETE")

	// Logs the listen and serve for the server
	log.Fatal(http.ListenAndServe(GetPort(), router))
}

// GET /paragliding/api/
// Handles the base path for the API
// Shows information about the API (uptime, information & version)
// Response: JSON
func APIIndex(w http.ResponseWriter, _ *http.Request) {
	// Set the information and version for the API
	var api = API{Info: "Service for Paragliding tracks.", Version: "v1"}

	// Calculate the uptime and convert it to ISO-8601 (duration)
	api.CalculateUptime(int(time.Since(startTime).Seconds()))

	// Set header content-type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode and displays the API-information as JSON
	json.NewEncoder(w).Encode(api)
}

// POST /paragliding/api/track/
// Handles the insertion of new tracks, and returns the ID for the track
// Body: JSON {"url": <url>}
// Response: JSON {"id": <id>}
func TrackPOST(w http.ResponseWriter, r *http.Request) {
	// Get the form value "url"
	url := r.FormValue("url")

	// Loop through all tracks in the database
	// Checking if the url is already inserted
	// Keeps computation time low, so it does not need to parse the IGC-file
	for _, track := range trackDB.GetAll() {
		// Checks if URL exits
		if track.TrackSrcUrl == url {
			// Show an 409, Conflict error
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}

	// Checks if the URL is not blank
	if url != "" {
		// Parse the IGC-url into a object
		track, err := igc.ParseLocation(url)

		// Checks if we got an error parsing the file
		if err != nil {
			// Shows an 400, Bad Request error
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Create a track from the IGC-object
		t := CreateTrackFromIGC(url, track)

		// Insert the track to the database
		// Gets the ID of the inserted track
		id, err := trackDB.Insert(t)

		// Check if we got an error inserting track to database
		if err != nil {
			// Shows an 400, Bad Request error
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Define a temporary struct for displaying ID as JSON
		type JSONID struct {
			Id bson.ObjectId `json:"id"`
		}

		// Loop through all Webhooks
		for _, wh := range webhookDB.GetAll() {
			// Check if the Webhook is triggered
			if wh.CheckTrigger() {
				// Sends the webhook
				wh.SendHook()
			}
		}

		// Set header content-type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Encode and displays the ID as JSON
		json.NewEncoder(w).Encode(JSONID{Id: id})
	}
}

// GET /paragliding/api/track/
// Handles the displaying of all Track-ID's
// Response: JSON [id0, id1, id2, ..., idn]
func TrackGET(w http.ResponseWriter, _ *http.Request) {
	// Get all tracks from the database
	var all = trackDB.GetAll()
	// Create a bson.ObjectID slice
	var ids []bson.ObjectId

	// Loop through all tracks
	for _, track := range all {
		// Append the track ID to the ID-slice
		ids = append(ids, track.Id)
	}

	// Set header content-type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode and displays all the track-ID's
	json.NewEncoder(w).Encode(ids)
}

// GET /paragliding/api/track/<id>/
// Handles the displaying of a single Track
// Response: JSON {Track}
func SingleTrackGET(w http.ResponseWriter, r *http.Request) {
	// Get the parameters with mux
	params := mux.Vars(r)

	// Get the <id> parameter
	id := params["id"]

	// Try to fetch the track from the database
	track, found := trackDB.Get(id)

	// Check if the track was not found
	if found != true {
		// Shows an 404, Not Found error
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Set header content-type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode and display the track as JSON
	json.NewEncoder(w).Encode(track)
}

// GET /paragliding/api/<id>/<field>/
// Handles the displaying of a field of a single track
// Response: text/plain <field>
func SingleTrackFieldGET(w http.ResponseWriter, r *http.Request) {
	// Get the parameters with mux
	params := mux.Vars(r)

	// Get the <id> parameter
	id := params["id"]

	// Get the track from the database
	track, found := trackDB.Get(id)

	// Check if the track was not found
	if found != true {
		// Show an 404, Not Found error
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Get the <field> parameter
	field := params["field"]

	// Get the field from the track
	response := track.GetField(field)

	// Set header content-type to plain text
	w.Header().Set("Content-Type", "text/plain")

	// Display the field-value
	fmt.Fprintf(w, "%s", response)
}

// GET /paragliding/api/ticker/
// Handles the displaying of the Ticker
// Shows timestamp of the first and last Track added
// Shows 5 latest Tracks
// Shows processing time
// Response: JSON {Ticker}
func GetTicker(w http.ResponseWriter, _ *http.Request) {
	// Get current time in milliseconds
	start := time.Now().Unix()

	// Create a Ticker object
	var ticker Ticker

	// Get all tracks from the database
	var all = trackDB.GetAll()

	// Creates a empty track-ID bson.ObjectId slice
	ticker.Tracks = make([]bson.ObjectId, 0)

	// Loop through all tracks
	for index, track := range all {
		// Checks if we are under 5 tracks
		if index < 5 {
			// Append track to ticker.Tracks
			ticker.Tracks = append(ticker.Tracks, track.Id)
		}
	}

	// Get the first tracks timestamp
	ticker.TStart = all[0].Timestamp

	var last int

	// Get the last index of the tracks
	if last = len(all) - 1; last > 4 {
		last = 4
	}

	// Get the last tracks timestamp
	ticker.TLatest = all[last].Timestamp
	// Get the last tracks timestamp
	ticker.TStop = all[last].Timestamp

	// Get time elapsed from the start of the function in milliseconds
	end := time.Now().Unix() - start

	// Set the Ticker.Processing to the time elapsed
	ticker.Processing = end

	// Set header content-type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode and display the ticker as JSON
	json.NewEncoder(w).Encode(ticker)
}

// GET /paragliding/api/ticker/latest/
// Handles the fetching and displaying of the latest inserted track
// Response: text/plain <timestamp>
func GetLatestTicker(w http.ResponseWriter, _ *http.Request) {
	// Check if the have any tracks in the database
	if trackDB.Count() < 1 {
		// Show an 404, Not Found error
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Declare a ticker variable
	var ticker Ticker
	// Get all tracks from the database
	var tracks = trackDB.GetAll()

	// Get the timestamp of the latest track
	ticker.TLatest = tracks[len(tracks) - 1].Timestamp

	// Set header content-type to plain text
	w.Header().Set("Content-Type", "text/plain")

	// Encode and display the latest timestamp
	json.NewEncoder(w).Encode(ticker.TLatest)
}

// GET /paragliding/api/ticker/<timestamp>/
// Handles the ticker of tracks with a higher timestamp then the parameter
// Response: JSON {Ticker}
func GetTickerFromTimestamp(w http.ResponseWriter, r *http.Request) {
	// Get the current time in milliseconds
	start := time.Now().Unix()

	// Declare a Ticker variable
	var ticker Ticker

	// Get parameters with mux
	var params = mux.Vars(r)

	// Convert the timestamp to a int64
	var timestamp, err = strconv.ParseInt(params["timestamp"], 10, 64)

	// Check if there was an error in the convertings
	if err != nil {
		// Show 400, Bad Request error
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Get all tracks from the database
	var all = trackDB.GetAll()
	// Create an empty slice with bson.ObjectId
	ticker.Tracks = make([]bson.ObjectId, 0)

	// Set count to 0
	count := 0

	// Create an empty slice of Tracks
	var result = make([]Track, 0)

	// Loop through all tracks
	for _, track := range all {
		// Check if we are under 5 counts
		if count < 5 {
			// Check that the timestamp is bigger than the parameter
			if track.Timestamp > timestamp {
				// Append the track
				result = append(result, track)
				// Increment the count
				count++
			}
		}
	}

	// set the first timestamp
	ticker.TStart = result[0].Timestamp
	// set the last timestamp
	ticker.TLatest = result[len(result) - 1].Timestamp
	ticker.TStop = result[len(result) - 1].Timestamp

	// loop through the results
	for _, track := range result {
		// Append the track-ID
		ticker.Tracks = append(ticker.Tracks, track.Id)
	}

	// Get the time elapsed
	end := time.Now().Unix() - start

	// set the time elapsed
	ticker.Processing = end

	// set header content-type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode and display the ticker as JSON
	json.NewEncoder(w).Encode(ticker)
}

// POST /paragliding/api/webhook/new_track
// Handles the insertion of a new Webhook
// Body: {"webhookURL": <url>, "minTriggerValue": <value>
func WebhookNewTrack(w http.ResponseWriter, r *http.Request) {
	// Check if the request body is empty
	if r.Body == nil {
		// Show 400, Bad Request error
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Declare a Webhook
	var hook Webhook

	// Decode the JSON from the body
	err := json.NewDecoder(r.Body).Decode(&hook)

	// Check if there was an error in decoding
	if err != nil {
		// Show 400, Bad Request error
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Checks the minTriggerValue
	hook.CheckTriggerValue()

	// Create a new ID for the webhook
	var id = bson.NewObjectId().Hex()
	// Set the id to the Webhook
	hook.ID = id

	// Insert webhook to the database
	webhookDB.Insert(hook)

	// Prints the ID of the webhook
	fmt.Fprintf(w, "%s", id)
}

// GET /paragliding/api/webhook/<webhook_id>
// Handles the displaying of information about a webhook
// Response: JSON {Webhook}
func WebhookNewTrackIdGET(w http.ResponseWriter, r *http.Request) {
	// Get the parameters with mux
	params := mux.Vars(r)

	// Get the <webhook_id>
	id := params["webhook_id"]

	// Get the webhook from the database
	webhook, found := webhookDB.Get(id)

	// Check if the webhook was not found
	if found != true {
		// Show 500, Internal Server Error
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Set header content-type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode and display the webhook as JSON
	json.NewEncoder(w).Encode(webhook)
}

// DELETE /paragliding/api/webhook/<webhook_id>
// Handles the deleting of a webhook
// Response: JSON {Webhook}
func WebhookNewTrackIdDELETE(w http.ResponseWriter, r *http.Request) {
	// Get the parameter with mux
	var params = mux.Vars(r)

	// Get the <webhook_id>
	var id = params["webhook_id"]

	// Get the webhook from the database
	webhook, found := webhookDB.Get(id)

	// Check if the webhook was not found
	if found != true {
		// Show 500, Internal Server Error
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Delete the webhook from the database
	webhookDB.Delete(id)

	// Set header content-type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode and display the deleted webhook
	json.NewEncoder(w).Encode(webhook)
}

// GET /paragliding/admin/api/track_count
// Retrieves the count of Tracks in the database
// Response: text/plain <count>
func AdminTrackCount(w http.ResponseWriter, _ *http.Request) {
	// Set content-type to plain text
	w.Header().Set("Content-Type", "text/plain")

	// Get the count from the database
	count := trackDB.Count()

	// Check if count is -1
	if count == -1 {
		// Show 404, Not Found Error
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Print the count
	fmt.Fprintf(w, "%d", count)
}

// DELETE /paragliding/admin/api/tracks
// Deletes all tracks from the database
// Response: text/plain "<count> tracks deleted"
func AdminDeleteTracks(w http.ResponseWriter, _ *http.Request) {
	// Get the count of tracks in the database
	count := trackDB.Count()

	// Check if the count is not -1
	if count != -1 {
		// Delete all tracks
		err := trackDB.DeleteAll()

		// Check if there was an error
		if err != nil {
			// Show 500, Internal Server Error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set content-type to plain text
		w.Header().Set("Content-Type", "text/plain")

		// Print the feedback
		fmt.Fprintf(w, "Deleted documents: %d", count)
	} else {
		// Show 404, Not Found Error
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}