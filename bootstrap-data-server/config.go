package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// ConfigUpdate update the configuration of the bootstrap service.
// It is only for manual administration via curl and not used by the API.
func ConfigUpdate(w http.ResponseWriter, r *http.Request) {

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read request: %s", err), http.StatusInternalServerError)
		return
	}

	// Get arguments from response body
	arguments := strings.Split(string(body), " ")
	if len(arguments) < 2 {
		http.Error(w, fmt.Sprintf("error: request has not enough arguments"),
			http.StatusInternalServerError)
		return
	}

	maxPeers, err := strconv.Atoi(arguments[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("first argument %q is not a number", arguments[0]),
			http.StatusInternalServerError)
		return
	}

	minRefillCandidates, err := strconv.Atoi(arguments[1])
	if err != nil {
		http.Error(w, fmt.Sprintf("second argument %q is not a number", arguments[1]),
			http.StatusInternalServerError)
		return
	}

	//if maxPeers < minRefillCandidates {
	//	http.Error(w, fmt.Sprintf("wrong relation between maxPeers and minPeers: %d < %d", maxPeers, minPeers),
	//		http.StatusInternalServerError)
	//	return
	//}

	// Read bootstrap data in JSON from file
	mtx.Lock()
	b, err := ioutil.ReadFile(collectionFileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read bootstrap peers from %q: %s", collectionFileName, err),
			http.StatusInternalServerError)
		return
	}
	mtx.Unlock()

	// Unmarshall JSON
	err = json.Unmarshal(b, &bootstrapData)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to unmarshal bootstrap peers from %q: %s", collectionFileName, err),
			http.StatusInternalServerError)
		return
	}

	bootstrapData.Config.MaxPeers = maxPeers
	bootstrapData.Config.MinRefillCandidates = minRefillCandidates

	// Marshall all bootstrap peers
	bootstrapPeersJson, err := json.MarshalIndent(bootstrapData, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshall bootstrap peers for %q: %v", collectionFileName, err),
			http.StatusInternalServerError)
		return
	}

	// Save JSON to file
	mtx.Lock()
	err = ioutil.WriteFile(collectionFileName, append(bootstrapPeersJson, byte('\n')), 0600)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write JSON file: %s", err), http.StatusInternalServerError)
	}
	mtx.Unlock()

	// Send JSON as response
	_, err = fmt.Fprintf(w, "%v", string(bootstrapPeersJson))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err), http.StatusInternalServerError)
	}
}

/*
curl -d "3 1" http://localhost:8080/config
*/
