package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Join(w http.ResponseWriter, r *http.Request) {

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read request: %s", err), http.StatusInternalServerError)
		return
	}

	// Get arguments from response body
	arguments := strings.Split(string(body), " ")
	if len(arguments) < 7 {
		http.Error(w, fmt.Sprintf("error: request has not enough arguments"),
			http.StatusInternalServerError)
		return
	}

	// Create the new peer struct
	newPeer := Peer{
		ID:       arguments[0],
		Name:     arguments[1],
		Ip:       arguments[2],
		Port:     arguments[3],
		Protocol: arguments[4],
		// todo get rid of unused argument status
		Status:    arguments[5],
		Timestamp: arguments[6],
	}

	// Read bootstrap data in JSON from file
	mtx.Lock()
	b, err := ioutil.ReadFile(collectionFileName)
	mtx.Unlock()

	// Unmarshall JSON
	err = json.Unmarshal(b, &bootstrapData)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to unmarshal bootstrap peers from %q: %s", collectionFileName, err),
			http.StatusInternalServerError)
		return
	}

	// If not enough bootstrap peers exist
	if len(bootstrapData.Peers) < bootstrapData.Config.MaxPeers {

		// Add new bootstrap peer
		bootstrapData.Peers[newPeer.ID] = &newPeer

		// Adjust number of peers
		bootstrapData.Config.NumPeers = len(bootstrapData.Peers)

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

		return
	}

	// Send as JSON response all the already existing bootstrap peers
	bootstrapPeersJson, err := json.MarshalIndent(bootstrapData, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshall bootstrap peers for %q: %v", collectionFileName, err),
			http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintf(w, "%v", string(bootstrapPeersJson))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err), http.StatusInternalServerError)
	}
}

/*

curl -d "uuid1 memberlist 127.0.0.1 22365 tcp 0 0" http://localhost:8080/join
curl -d "uuid2 memberlist 127.0.0.1 22365 tcp 0 0" http://localhost:8080/join
curl -d "uuid3 memberlist 127.0.0.1 22365 tcp 0 0" http://localhost:8080/join
curl -d "uuid4 memberlist 127.0.0.1 22365 tcp 0 0" http://localhost:8080/join
curl -d "uuid5 memberlist 127.0.0.1 22365 tcp 0 0" http://localhost:8080/join

*/
