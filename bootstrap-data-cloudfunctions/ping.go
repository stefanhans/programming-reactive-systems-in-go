package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
)

func Ping(w http.ResponseWriter, r *http.Request) {

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

	// Send response
	_, err = fmt.Fprintf(w, "OK")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err), http.StatusInternalServerError)
	}
}

/*
DO NOT FORGET:

export GCP_PROJECT="bootstrap-data-cloudfunctions"
export GOOGLE_APPLICATION_CREDENTIALS="/Users/stefan/.secret/bootstrap-data-cloudfunctions-c628b7847572.json"

cd ~/go/src/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-cloudfunctions
gcloud functions deploy ping --region europe-west1 --entry-point Ping --runtime go111 --trigger-http

curl https://europe-west1-bootstrap-data-cloudfunctions.cloudfunctions.net/ping

*/
