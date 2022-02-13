package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

var pool *Pool
var currentState []byte

type State struct {
	Added      interface{} `json:"added"`
	Previously interface{} `json:"previously"`
}

// GameStateHandler
// When any data of a reported component changes, the game client will
// start the buffer timer, and after it elapses, or immediately if buffering
// is disabled, will report the new full game state, and delta states computed
// from payload with last successful HTTP 2XX response from the endpoint.
// For data that changed from the previous report a global block "previously"
// will contain old state of the changed game state, whereas for new data a global
// block "added" will contain the root JSON fields that are present in the new game state,
// but were absent in the old game state.
func GameStateHandler(_ http.ResponseWriter, request *http.Request) {
	var err error
	currentState, err = ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Printf("Error reading body: %s\n", err.Error())
	}

	buffer := new(bytes.Buffer)
	json.Compact(buffer, currentState)
	state := buffer.Bytes()

	pool.Broadcast <- state
}

func main() {
	pool = NewPool()
	go pool.Start()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
	router.HandleFunc("/gsi-feed", GameStateHandler).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./client/dist/")))

	log.Fatalln(http.ListenAndServe(":9090", router))
}
