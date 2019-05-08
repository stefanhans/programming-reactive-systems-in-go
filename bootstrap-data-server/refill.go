package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Refill checks to add the requesting peer in the list of bootstrap peers.
// If a bootstrap peer has left, other peers do ask to fill the gap.
// Criterium is the number of needed peers, and,
// the peer with the oldest timestamp will succeed then.
func Refill(w http.ResponseWriter, r *http.Request) {

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

	// Create the candidate peer struct
	candidatePeer := &Peer{
		ID:       arguments[0],
		Name:     arguments[1],
		Ip:       arguments[2],
		Port:     arguments[3],
		Protocol: arguments[4],
		// todo get rid of unused argument status
		Status:    arguments[5],
		Timestamp: arguments[6],
	}

	// Read bootstrap peers in JSON from file
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

	// todo: ping all peers and clean list accordingly

	bootstrapData.Config.NumPeers = len(bootstrapData.Peers)

	// Candidate peer not already a bootstrap peer
	if _, ok := bootstrapData.Peers[candidatePeer.ID]; !ok {

		// Fill or replace?
		if bootstrapData.Config.NumPeers < bootstrapData.Config.MaxPeers {

			// Add new bootstrap peer
			bootstrapData.Peers[candidatePeer.ID] = candidatePeer
			bootstrapData.Config.NumPeers++
		} else {

			newestPeer := candidatePeer

			// Find the newest peer
			for _, peer := range bootstrapData.Peers {
				if newestPeer.Timestamp < peer.Timestamp {
					newestPeer = peer
				}
			}

			if newestPeer != candidatePeer {

				// Delete newest peer
				delete(bootstrapData.Peers, newestPeer.ID)

				// Add new bootstrap peer
				bootstrapData.Peers[candidatePeer.ID] = candidatePeer
			}
		}
	}

	// Marshall bootstrap data
	bootstrapDataJson, err := json.MarshalIndent(bootstrapData, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshall bootstrapData: %v", err), http.StatusInternalServerError)
		return
	}

	// Save JSON to file
	mtx.Lock()
	err = ioutil.WriteFile(collectionFileName, append(bootstrapDataJson, byte('\n')), 0600)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write JSON file: %s", err), http.StatusInternalServerError)
	}
	mtx.Unlock()

	// Send JSON as response
	_, err = fmt.Fprintf(w, "%v", string(bootstrapDataJson))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err), http.StatusInternalServerError)
	}
}

/*
curl -d "uuid1 memberlist 127.0.0.1 22365 tcp 0 0" http://localhost:8080/refill
*/
