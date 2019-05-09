package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func saveTestJsonData(w http.ResponseWriter) {

	// Marshal array of struct
	testEventFiltersJson, err := json.MarshalIndent(currentTestEventFilters, "", " ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write marshal currentTestEventFilter: %s", err),
			http.StatusInternalServerError)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/testEventFilters.json", currentTestDir),
		append(testEventFiltersJson, byte('\n')), 0600)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write testeventfilters to test directory: %s", err),
			http.StatusInternalServerError)
	}

	currentTestSummaryJson, err := json.MarshalIndent(currentTestSummary, "", " ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal currentTestSummary: %s", err),
			http.StatusInternalServerError)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/testSummary.json", currentTestDir),
		append(currentTestSummaryJson, byte('\n')), 0600)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write testsummary to test directory: %s", err),
			http.StatusInternalServerError)
	}
}

// PrepareTestSummary prepares the summary after the test run.
// It saves the data of the run, the filters, and the summary in JSON
// files in the directory of the test run.
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

	for _, currentTestSummary := range currentTestSummary {
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

			currentTestEvaluation.ID = currentTestResult.ID
			currentTestEvaluation.Name = currentTestResult.Name
			currentTestEvaluation.Peer = testPeer
			currentTestEvaluation.Kind = "command"
			currentTestEvaluation.Status = result.Status
			currentTestEvaluation.Test = result.Data
			currentTestEvaluation.Result = result.Data

			currentTestSummary = append(currentTestSummary, currentTestEvaluation)

			if strings.Split(result.Data, " ")[0] == "testfilter" {

				source := strings.Split(result.Data, " ")[1]
				filter := strings.Join(strings.Split(result.Data, " ")[3:], " ")

				currentTestEvaluation.ID = currentTestResult.ID
				currentTestEvaluation.Name = currentTestResult.Name
				currentTestEvaluation.Peer = testPeer
				currentTestEvaluation.Kind = "event"
				currentTestEvaluation.Test = result.Data

				currentTestEvaluation.Status = "FAILED"
				for _, eventFilter := range currentTestEventFilters {

					fmt.Printf("\nXXX\n eventFilter: %v\n\n", eventFilter)

					fmt.Printf("NumExpectedEvents %d == NumReceivedEvents %d\n",
						eventFilter.NumExpectedEvents, eventFilter.NumReceivedEvents)

					fmt.Printf("eventFilter.Peer %q == testPeer %q\n",
						eventFilter.Peer, testPeer)
					fmt.Printf("eventFilter.Source %q == source %q\n",
						eventFilter.Source, source)
					fmt.Printf("eventFilter.Filter %q == filter %q\n",
						eventFilter.Filter, filter)

					currentTestEvaluation.Comment = fmt.Sprintf("NumExpectedEvents: %d, NumReceivedEvents: %d",
						eventFilter.NumExpectedEvents, eventFilter.NumReceivedEvents)

					if eventFilter.Peer == testPeer &&
						eventFilter.Source == source &&
						eventFilter.Filter == filter &&
						eventFilter.NumExpectedEvents == eventFilter.NumReceivedEvents {

						currentTestEvaluation.Status = "OK"
					}
				}
				currentTestEvaluation.Result = result.Data
				currentTestSummary = append(currentTestSummary, currentTestEvaluation)
			}

			fmt.Printf("%d result: %s\n", i, result)
		}
	}

	saveTestJsonData(w)

	_, err = fmt.Fprintf(w, "")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err),
			http.StatusInternalServerError)
	}
}

/*
curl -d "alice" http://localhost:8081/preparesummary
*/

// GetTestSummary sends back the summary of the last test run.
func GetTestSummary(w http.ResponseWriter, r *http.Request) {

	// Get rid of warning
	_ = r

	// Marshal array of struct
	currentTestSummariesJson, err := json.MarshalIndent(currentTestSummary, "", " ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write marshal currentTestEvaluation: %s", err),
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
