package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

var startTime time.Time

func main() {
	startTime = time.Now()

	router := httprouter.New()

	router.GET("/paragliding/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.Redirect(w, r, "/paragliding/api", http.StatusMovedPermanently)
	})

	router.GET("/paragliding/api/", APIIndex)


	log.Fatal(http.ListenAndServe(":8088", router))
}

func APIIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var api = API{Info: "Service for Paragliding tracks.", Version: "v1"}

	api.CalculateUptime(int(time.Since(startTime).Seconds()))

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(api)
}

