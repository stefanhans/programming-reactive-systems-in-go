package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
)

// Refill checks to add the requesting peer in the list of bootstrap peers.
// If a bootstrap peer has left, other peers do ask to fill the gap.
// Criterium is the number of needed peers, and,
// the peer with the oldest timestamp will succeed then.
func Refill(w http.ResponseWriter, r *http.Request) {

	// Sets your Google Cloud Platform project ID.
	projectId := os.Getenv("GCP_PROJECT")
	if projectId == "" {
		http.Error(w, fmt.Sprintf("GCP_PROJECT environment variable unset or missing"), http.StatusInternalServerError)
	}

	// Get a Firestore client.
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create client: %v", err), http.StatusInternalServerError)
		return
	}
	defer client.Close()

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

	// Get document from Firestore
	document, err := client.Collection(collectionName).Doc(documentName).Get(ctx)
	if err != nil {
		// Return HTTP error code 500 Internal Server Error
		http.Error(w, fmt.Sprintf("failed to get document %q from collection %q: %s",
			documentName, collectionName, err), http.StatusInternalServerError)
	}

	// Local variable BootstrapData
	var bootstrapData BootstrapData

	// Get the JSON string and unmarshall it
	if v, ok := document.Data()[topic].(string); ok {
		err = json.Unmarshal([]byte(v), &bootstrapData)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to unmarshall JSON from collection %q: %s", collectionName, err), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, fmt.Sprintf("failed to get JSON from collection %q: %v", collectionName, err), http.StatusInternalServerError)
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

	// Save the JSON as string in collection of documents of topics
	_, err = client.Collection(collectionName).Doc(documentName).Set(ctx,
		map[string]interface{}{
			topic: fmt.Sprintf("%v", string(bootstrapDataJson)),
		})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to Reset %v: %s", collectionName, err), http.StatusInternalServerError)
		return
	}

	// Send JSON as response
	_, err = fmt.Fprintf(w, "%v", string(bootstrapDataJson))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err), http.StatusInternalServerError)
	}
}

/*
DO NOT FORGET:

export GCP_PROJECT="bootstrap-data-cloudfunctions"
export GOOGLE_APPLICATION_CREDENTIALS="/Users/stefan/.secret/bootstrap-data-cloudfunctions-c628b7847572.json"

cd ~/go/src/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-cloudfunctions
gcloud functions deploy refill --region europe-west1 --entry-point Refill --runtime go111 --trigger-http

curl -d "uuid memberlist 127.0.0.1 22365 tcp test 123" https://europe-west1-bootstrap-data-cloudfunctions.cloudfunctions.net/refill
curl -d "uuid2 memberlist 127.0.0.1 22366 tcp test 124" https://europe-west1-bootstrap-data-cloudfunctions.cloudfunctions.net/refill
curl -d "uuid3 memberlist 127.0.0.1 22367 tcp test 125" https://europe-west1-bootstrap-data-cloudfunctions.cloudfunctions.net/refill
curl -d "uuid4 memberlist 127.0.0.1 22367 tcp test 126" https://europe-west1-bootstrap-data-cloudfunctions.cloudfunctions.net/refill

*/
