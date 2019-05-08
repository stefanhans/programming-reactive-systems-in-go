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

// Leave get information of the leaving peer via http and deletes it in Firestore
func Leave(w http.ResponseWriter, r *http.Request) {

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
	if len(arguments) < 1 {
		http.Error(w, fmt.Sprintf("error: request has not enough arguments"),
			http.StatusInternalServerError)
		return
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

	// Delete sent bootstrap peer
	delete(bootstrapData.Peers, arguments[0])

	// Save number of peers in configuration
	bootstrapData.Config.NumPeers = len(bootstrapData.Peers)

	// Marshall all bootstrap data
	bootstrapDataJson, err := json.MarshalIndent(bootstrapData, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshall bootstrap peers for %q: %v", collectionName, err),
			http.StatusInternalServerError)
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

export GCP_PROJECT="chat-bootstrap-peers"
export GOOGLE_APPLICATION_CREDENTIALS="/Users/stefan/.secret/chat-bootstrap-peers-c4fe2a951411.json"

cd ~/go/src/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-cloudfunctions
gcloud functions deploy leave --region europe-west1 --entry-point Leave --runtime go111 --trigger-http

curl -d "uuid" https://europe-west1-bootstrap-data-cloudfunctions.cloudfunctions.net/leave

*/
