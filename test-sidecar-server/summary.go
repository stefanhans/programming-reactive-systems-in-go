package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func PrepareTestSummary(w http.ResponseWriter, r *http.Request) {

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read request: %s", err),
			http.StatusInternalServerError)
		return
	}

	// Get arguments from response body
	arguments := strings.Split(string(body), " ")
	if len(arguments) < 1 {
		http.Error(w, fmt.Sprintf("error: request has not enough arguments"),
			http.StatusInternalServerError)
		return
	}

	testPeer := arguments[0]

	for _, currentTestSummary := range currentTestSummaries {
		if currentTestSummary.Peer == testPeer {
			http.Error(w, fmt.Sprintf("error: summary of %s already processed\n", testPeer),
				http.StatusInternalServerError)
			return
		}
	}

	fmt.Printf("GetTestSummary: %s\n", string(decodeJsonBytes(body)))

	fmt.Printf("currentTestResult.ID: %s\n", currentTestResult.ID)
	fmt.Printf("currentTestResult.Name: %s\n", currentTestResult.Name)

	for i, result := range currentTestResult.CommandResults {
		if result.Peer == testPeer {

			currentTestSummary.ID = currentTestResult.ID
			currentTestSummary.Name = currentTestResult.Name
			currentTestSummary.Peer = testPeer
			currentTestSummary.Kind = "command"
			currentTestSummary.Status = result.Status
			currentTestSummary.Test = result.Data
			currentTestSummary.Result = result.Data

			currentTestSummaries = append(currentTestSummaries, currentTestSummary)

			if strings.Split(result.Data, " ")[0] == "testfilter" {

				source := strings.Split(result.Data, " ")[1]
				filter := strings.Join(strings.Split(result.Data, " ")[3:], " ")

				currentTestSummary.ID = currentTestResult.ID
				currentTestSummary.Name = currentTestResult.Name
				currentTestSummary.Peer = testPeer
				currentTestSummary.Kind = "event"
				currentTestSummary.Test = result.Data

				currentTestSummary.Status = "FAILED"
				for _, eventFilter := range currentTestEventFilters {

					fmt.Printf("\nXXX\n eventFilter: %s\n\n", eventFilter)

					fmt.Printf("NumExpectedEvents %d == NumReceivedEvents %d\n",
						eventFilter.NumExpectedEvents, eventFilter.NumReceivedEvents)

					fmt.Printf("eventFilter.Peer %q == testPeer %q\n",
						eventFilter.Peer, testPeer)
					fmt.Printf("eventFilter.Source %q == source %q\n",
						eventFilter.Source, source)
					fmt.Printf("eventFilter.Filter %q == filter %q\n",
						eventFilter.Filter, filter)

					currentTestSummary.Comment = fmt.Sprintf("NumExpectedEvents: %d, NumReceivedEvents: %d",
						eventFilter.NumExpectedEvents, eventFilter.NumReceivedEvents)

					if eventFilter.Peer == testPeer &&
						eventFilter.Source == source &&
						eventFilter.Filter == filter &&
						eventFilter.NumExpectedEvents == eventFilter.NumReceivedEvents {

						currentTestSummary.Status = "OK"
					}
				}
				currentTestSummary.Result = result.Data
				currentTestSummaries = append(currentTestSummaries, currentTestSummary)
			}

			fmt.Printf("%d result: %s\n", i, result)
		}
	}

	_, err = fmt.Fprintf(w, "")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err),
			http.StatusInternalServerError)
	}
}

/*
curl -d "alice" http://localhost:8081/preparesummary
*/

func GetTestSummary(w http.ResponseWriter, r *http.Request) {

	// Marshal array of struct
	currentTestSummariesJson, err := json.MarshalIndent(currentTestSummaries, "", " ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write marshal currentTestSummary: %s", err),
			http.StatusInternalServerError)
	}

	_, err = fmt.Fprintf(w, string(currentTestSummariesJson))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err),
			http.StatusInternalServerError)
	}
}

/*
curl http://localhost:8081/getsummary
*/
