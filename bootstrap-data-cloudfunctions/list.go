package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
)

// List returns all bootstrap data from Firestore
func List(w http.ResponseWriter, r *http.Request) {

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

	// Get document from Firestore
	document, err := client.Collection(collectionName).Doc(documentName).Get(ctx)
	if err != nil {
		// Return HTTP error code 500 Internal Server Error
		http.Error(w, fmt.Sprintf("failed to get document %q from collection %q: %s",
			documentName, collectionName, err), http.StatusInternalServerError)
	}

	// Get the JSON string and unmarshall it
	if v, ok := document.Data()[topic].(string); ok {

		// Send JSON as response
		_, err = fmt.Fprint(w, string(v))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to write response: %s", err), http.StatusInternalServerError)
		}

	} else {
		http.Error(w, fmt.Sprintf("no data found of field %q in document %q from collection %q",
			topic, documentName, collectionName), http.StatusInternalServerError)
	}
}

/*
DO NOT FORGET:

export GCP_PROJECT="bootstrap-data-cloudfunctions"
export GOOGLE_APPLICATION_CREDENTIALS="/Users/stefan/.secret/bootstrap-data-cloudfunctions-c628b7847572.json"

cd ~/go/src/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-cloudfunctions
gcloud functions deploy list --region europe-west1 --entry-point List --runtime go111 --trigger-http

curl https://europe-west1-bootstrap-data-cloudfunctions.cloudfunctions.net/list

*/
