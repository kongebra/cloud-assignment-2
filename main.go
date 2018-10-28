package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/marni/goigc"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

var startTime time.Time
var tickerDB TickerDB
var trackDB TrackDB

var ticker Ticker

func main() {
	trackDB = TrackDB{
		Addrs: []string{"ds133533.mlab.com:33533"},
		Database: "assignment-2",
		Username: "golang",
		Password: "golang1",
		Collection: "tracks",
	}

	trackDB.Init()

	tickerDB = TickerDB{
		Addrs: []string{"ds133533.mlab.com:33533"},
		Database: "assignment-2",
		Username: "golang",
		Password: "golang1",
		Collection: "tickers",
	}

	tickerDB.Init()

	startTime = time.Now()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/paragliding/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/paragliding/api/", http.StatusMovedPermanently)
	})

	router.HandleFunc("/paragliding/api/", APIIndex).Methods("GET")

	router.HandleFunc("/paragliding/api/track/", TrackPOST).Methods("POST")
	router.HandleFunc("/paragliding/api/track/", TrackGET).Methods("GET")
	router.HandleFunc("/paragliding/api/track/{id}/", SingleTrackGET).Methods("GET")
	router.HandleFunc("/paragliding/api/track/{id}/{field}/", SingleTrackFieldGET).Methods("GET")

	router.HandleFunc("/paragliding/api/ticker", GetTicker).Methods("GET")
	router.HandleFunc("/paragliding/api/ticker/latest", GetLatestTicker).Methods("GET")
	router.HandleFunc("/paragliding/api/ticker/{timestamp}", GetTickerFromTimestamp).Methods("GET")


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

		if ticker.TStart == 0 {
			ticker.TStart = time.Now().Unix()
		}

		ticker.Timestamp()

		type JSONID struct {
			Id bson.ObjectId `json:"id"`
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
	w.Header().Set("Content-Type", "text/plain")

	json.NewEncoder(w).Encode(ticker.TLatest)
}

func GetTickerFromTimestamp(w http.ResponseWriter, r *http.Request) {

}