package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func List(w http.ResponseWriter, r *http.Request) {

	// Get rid off warning
	_ = r

	// Read bootstrap peers in JSON from file
	mtx.Lock()
	b, _ := ioutil.ReadFile(collectionFileName)
	mtx.Unlock()

	// Send as JSON response
	_, err := fmt.Fprint(w, string(b))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err), http.StatusInternalServerError)
	}
}

/*
curl http://localhost:8080/list
*/
