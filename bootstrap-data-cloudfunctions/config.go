package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
)

// ConfigUpdate update the configuration of the bootstrap service.
// It is only for manual administration via curl and not used by the API.
func ConfigUpdate(w http.ResponseWriter, r *http.Request) {

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

	bootstrapData.Config.MaxPeers = maxPeers
	bootstrapData.Config.MinRefillCandidates = minRefillCandidates

	// Marshall bootstrap data
	bootstrapDataJson, err := json.MarshalIndent(bootstrapData, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshall bootstrap data for %q: %v", collectionName, err),
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
gcloud functions deploy config --region europe-west1 --entry-point ConfigUpdate --runtime go111 --trigger-http

curl -d "2 1" https://europe-west1-bootstrap-data-cloudfunctions.cloudfunctions.net/config

*/
