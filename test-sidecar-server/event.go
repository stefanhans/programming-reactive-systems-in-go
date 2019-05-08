package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func decodeJsonBytes(txt []byte) []byte {
	txt = bytes.Replace(txt, []byte(`\u003c`), []byte("<"), -1)
	txt = bytes.Replace(txt, []byte(`\u003e`), []byte(">"), -1)
	return bytes.Replace(txt, []byte(`\u0026`), []byte("&"), -1)
}

func PutTestEvent(w http.ResponseWriter, r *http.Request) {

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read request: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("PutTestEvent: " + string(decodeJsonBytes(body)))

	mtx.Lock()

	err = json.Unmarshal(body, &currentTestEventFilter)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to unmarshall currentTestEventFilter: %s", err),
			http.StatusInternalServerError)
	}

	fmt.Printf("\n--------\ncurrentTestEventFilter: %v\n", currentTestEventFilter)

	isNew := true
	for i, testEventFilter := range currentTestEventFilters {

		fmt.Printf("\nLOOP testEventFilter: %v\n", testEventFilter)

		if testEventFilter.Source == currentTestEventFilter.Source &&
			testEventFilter.Peer == currentTestEventFilter.Peer &&
			testEventFilter.Filter == currentTestEventFilter.Filter {
			currentTestEventFilters[i].NumReceivedEvents++
			isNew = false
			fmt.Printf("\nUPDATED\ncurrentTestEventFilter: %v\n", currentTestEventFilter)

			continue

		}
	}

	if isNew {
		currentTestEventFilter.NumReceivedEvents = 1
		currentTestEventFilters = append(currentTestEventFilters, currentTestEventFilter)
		fmt.Printf("\nAPPENDED\ncurrentTestEventFilters: %v\n", currentTestEventFilters)
	}

	testEventFilterJson, err := json.MarshalIndent(currentTestEventFilter, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshall currentTestEventFilter: %s", err),
			http.StatusInternalServerError)
		mtx.Unlock()
	}
	mtx.Unlock()

	_, err = fmt.Fprintf(w, string(testEventFilterJson))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err),
			http.StatusInternalServerError)
	}
}

/*
curl -d "testID testname alice OK" http://localhost:8081/putresult
curl -d "testID testname alice OK just an comment" http://localhost:8081/putresult
*/

func GetTestEvents(w http.ResponseWriter, r *http.Request) {

	// Get rid of warning
	_ = r

	// Marshal array of struct
	testEventFiltersJson, err := json.MarshalIndent(currentTestEventFilters, "", " ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write marshal currentTestEventFilter: %s", err),
			http.StatusInternalServerError)
	}

	_, err = fmt.Fprintf(w, string(testEventFiltersJson))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err),
			http.StatusInternalServerError)
	}
}

/*
curl http://localhost:8081/getevents
*/
