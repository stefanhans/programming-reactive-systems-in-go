package bootstrap_data_cloudfunctions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
)

// Reset deletes all documents of the collection from Firestore
func Reset(w http.ResponseWriter, r *http.Request) {

	// Get rid of warnings
	_ = r

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

	// Reset the bootstrap data
	for k := range bootstrapData.Peers {
		delete(bootstrapData.Peers, k)
	}
	bootstrapData.Config.NumPeers = 0

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

export GCP_PROJECT="bootstrap-data-cloudfunctions"
export GOOGLE_APPLICATION_CREDENTIALS="/Users/stefan/.secret/bootstrap-data-cloudfunctions-c628b7847572.json"

cd ~/go/src/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-cloudfunctions
gcloud functions deploy reset --region europe-west1 --entry-point Reset --runtime go111 --trigger-http

curl https://europe-west1-bootstrap-data-cloudfunctions.cloudfunctions.net/reset

*/
