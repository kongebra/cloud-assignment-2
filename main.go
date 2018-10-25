package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("listening...")

	err := http.ListenAndServe(GetPort(), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

}

func GetPort() string {
	var port = os.Getenv("PORT")

	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}

	return ":" + port
}