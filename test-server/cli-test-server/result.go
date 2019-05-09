package cli_test_server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// PutTestResult receives the result of a test command
func PutTestResult(w http.ResponseWriter, r *http.Request) {

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read request: %s", err), http.StatusInternalServerError)
		return
	}

	// Get arguments from response body
	arguments := strings.Split(string(body), " ")
	if len(arguments) < 4 {
		http.Error(w, fmt.Sprintf("error: request has not enough arguments"),
			http.StatusInternalServerError)
		return
	}

	testId := arguments[0]
	testName := arguments[1]
	testPeer := arguments[2]
	testStatus := arguments[3]

	testComment := ""
	if len(arguments) > 4 {
		testComment = strings.Join(arguments[4:], " ")
	}

	mtx.Lock()

	if currentTestResult.ID == "" {
		currentTestResult.ID = testId
		currentTestResult.Name = testName
	}
	currentTestResult.CommandResults = append(currentTestResult.CommandResults, &CommandResult{
		Peer:   testPeer,
		Status: testStatus,
		Data:   testComment,
	})

	mtx.Unlock()

	// Marshal array of struct
	testResultJson, err := json.MarshalIndent(currentTestResult, "", " ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write marshal testRun: %s", err),
			http.StatusInternalServerError)
	}

	_, err = fmt.Fprintf(w, string(testResultJson))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err),
			http.StatusInternalServerError)
	}
}

/*
curl -d "testID testname alice OK" http://localhost:8081/putresult
curl -d "testID testname alice OK just an comment" http://localhost:8081/putresult
*/

// GetTestResults sends back all test results
func GetTestResults(w http.ResponseWriter, r *http.Request) {

	// Get rid of warnings
	_ = r

	mtx.Lock()

	// Marshal array of struct
	testResultJson, err := json.MarshalIndent(currentTestResult, "", " ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write marshal testRun: %s", err),
			http.StatusInternalServerError)
	}

	mtx.Unlock()

	_, err = fmt.Fprintf(w, string(testResultJson))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err),
			http.StatusInternalServerError)
	}
}

/*
curl http://localhost:8081/getresults
*/
