package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Leave(w http.ResponseWriter, r *http.Request) {

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read request: %s", err), http.StatusInternalServerError)
		return
	}

	// Get arguments from response body
	arguments := strings.Split(string(body), " ")

	// Read bootstrap peers in JSON from file
	mtx.Lock()
	b, err := ioutil.ReadFile(collectionFileName)
	mtx.Unlock()

	// Unmarshall JSON
	err = json.Unmarshal(b, &bootstrapData)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to unmarshal bootstrap data from %q: %s", collectionFileName, err),
			http.StatusInternalServerError)
		return
	}

	// Delete sent bootstrap peer
	delete(bootstrapData.Peers, arguments[0])

	// Save number of peers in configuration
	bootstrapData.Config.NumPeers = len(bootstrapData.Peers)

	// Marshall all remained bootstrap peers
	bootstrapDataJson, err := json.MarshalIndent(bootstrapData, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshall bootstrap data for %q: %v", collectionFileName, err),
			http.StatusInternalServerError)
		return
	}

	// Save JSON to file
	mtx.Lock()
	err = ioutil.WriteFile(collectionFileName, append(bootstrapDataJson, byte('\n')), 0600)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write bootstrap data for %q: %v", collectionFileName, err),
			http.StatusInternalServerError)
		return
	}
	mtx.Unlock()

	// Send JSON as response
	_, err = fmt.Fprintf(w, "%v", string(bootstrapDataJson))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err), http.StatusInternalServerError)
	}
}

/*
curl -d "uuid1" http://localhost:8080/leave
*/
