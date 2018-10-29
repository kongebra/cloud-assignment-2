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

var startTime time.Time
var trackDB TrackDB
var webhookDB WebhookDB

func main() {
	go startClockTicker()

	trackDB = TrackDB{
		Addrs: []string{"ds133533.mlab.com:33533"},
		Database: "assignment-2",
		Username: "golang",
		Password: "golang1",
		Collection: "tracks",
	}

	trackDB.Init()

	webhookDB = WebhookDB{
		Addrs: []string{"ds133533.mlab.com:33533"},
		Database: "assignment-2",
		Username: "golang",
		Password: "golang1",
		Collection: "webhooks",
	}

	webhookDB.Init()

	startTime = time.Now()

	router := mux.NewRouter().StrictSlash(true)

	// TODO: Only for testing triggers
	router.HandleFunc("/webhook/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got a request to /webhook/test")
	}).Methods("GET")

	router.HandleFunc("/paragliding/api/discord/", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		fmt.Println(params["code"])
	})

	router.HandleFunc("/paragliding/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/paragliding/api/", http.StatusMovedPermanently)
	})

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

	log.Fatal(http.ListenAndServe(GetPort(), router))
}

func APIIndex(w http.ResponseWriter, r *http.Request) {
	var api = API{Info: "Service for Paragliding tracks.", Version: "v1"}

	api.CalculateUptime(int(time.Since(startTime).Seconds()))

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(api)
}

func TrackPOST(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	for _, track := range trackDB.GetAll() {
		if track.TrackSrcUrl == url {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}

	if url != "" {
		track, err := igc.ParseLocation(url)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		t := CreateTrackFromIGC(url, track)

		id, err := trackDB.Insert(t)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		type JSONID struct {
			Id bson.ObjectId `json:"id"`
		}

		for _, wh := range webhookDB.GetAll() {
			if wh.CheckTrigger() {
				wh.SendHook()
			}
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(JSONID{Id: id})
	}
}

func TrackGET(w http.ResponseWriter, r *http.Request) {
	var all = trackDB.GetAll()
	var ids []bson.ObjectId

	for _, value := range all {
		ids = append(ids, value.Id)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(ids)
}

func SingleTrackGET(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	response, found := trackDB.Get(id)

	if found != true {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

func SingleTrackFieldGET(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	track, found := trackDB.Get(id)

	if found != true {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	field := params["field"]

	response := track.GetField(field)

	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "%s", response)
}

func GetTicker(w http.ResponseWriter, r *http.Request) {
	start := time.Now().Unix()

	var ticker Ticker

	var all = trackDB.GetAll()
	ticker.Tracks = make([]bson.ObjectId, 0)

	for i, value := range all {
		if i < 5 {
			ticker.Tracks = append(ticker.Tracks, value.Id)
		}
	}

	ticker.TStart = all[0].Timestamp
	ticker.TLatest = all[len(all) - 1].Timestamp
	ticker.TStop = all[len(all) - 1].Timestamp

	end := time.Now().Unix() - start

	ticker.Processing = end

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(ticker)
}

func GetLatestTicker(w http.ResponseWriter, r *http.Request) {
	if trackDB.Count() < 1 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var ticker Ticker
	var tracks = trackDB.GetAll()

	ticker.TLatest = tracks[len(tracks) - 1].Timestamp

	w.Header().Set("Content-Type", "text/plain")

	json.NewEncoder(w).Encode(ticker.TLatest)
}

func GetTickerFromTimestamp(w http.ResponseWriter, r *http.Request) {
	start := time.Now().Unix()

	var ticker Ticker

	var params = mux.Vars(r)

	var timestamp, err = strconv.ParseInt(params["timestamp"], 10, 64)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var all = trackDB.GetAll()
	ticker.Tracks = make([]bson.ObjectId, 0)

	count := 0

	var result = make([]Track, 0)

	for _, value := range all {
		if count < 5 {
			if value.Timestamp > timestamp {
				result = append(result, value)
				count++
			}
		}
	}

	ticker.TStart = result[0].Timestamp
	ticker.TLatest = result[len(result) - 1].Timestamp
	ticker.TStop = result[len(result) - 1].Timestamp

	for _, t := range result {
		ticker.Tracks = append(ticker.Tracks, t.Id)
	}

	end := time.Now().Unix() - start

	ticker.Processing = end

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(ticker)
}

func WebhookNewTrack(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var hook Webhook

	err := json.NewDecoder(r.Body).Decode(&hook)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	hook.CheckTriggerValue()

	var id = bson.NewObjectId().Hex()
	hook.ID = id

	webhookDB.Insert(hook)
	//webhooks.Add(hook)

	fmt.Fprintf(w, "%s", id)
}

func WebhookNewTrackIdGET(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["webhook_id"]

	webhook, found := webhookDB.Get(id)

	if found != true {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(webhook)
}

func WebhookNewTrackIdDELETE(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)

	var id = params["webhook_id"]

	webhook, found := webhookDB.Get(id)

	if found != true {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	webhookDB.Delete(id)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(webhook)
}

func AdminTrackCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	count := trackDB.Count()

	if count == -1 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "%d", count)
}

func AdminDeleteTracks(w http.ResponseWriter, r *http.Request) {
	count := trackDB.Count()

	if count != -1 {
		err := trackDB.DeleteAll()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Deleted documents: %d", count)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}

