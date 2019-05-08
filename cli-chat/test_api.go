package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Todo: Make the API a library

// initTest sends a name to initialize an existing test
// to the test-sidecar-server running under http://localhost:8081
func initTest(testName string) bool {

	// Send request to service
	res, err := http.Post(testSidecarUrl+"/init",
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("%s",
			testName)))
	if err != nil {
		displayError("failed to post request to test sidecar", err)
		return false
	}
	if res.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(res.Body)

		displayError(fmt.Sprintf("Received no \"200 OK\" from test sidecar init: %q",
			strings.TrimSuffix(string(b), "\n")))
		return false

	}
	//fmt.Printf("Received reply from Join: %v\n", res.Status)

	// Read response body in JSON
	body, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		displayError("failed to read response from test sidecar", err)
		return false
	}

	err = json.Unmarshal(body, &currentTestRun)
	if err != nil {
		displayError("failed to unmarshall response from test sidecar", err)
		return false
	}

	currentTestEventFilter.ID = currentTestRun.ID
	currentTestEventFilter.Name = currentTestRun.Name

	displayGreenText(fmt.Sprintf("Test %s starting", currentTestRun.ID))

	//testRunJson, err := json.MarshalIndent(currentTestRun, "", "  ")
	//if err != nil {
	//	displayError("failed to marshall test run", err)
	//	return false
	//}
	//displayGreenText(string(testRunJson))

	return true
}

func refreshTest() bool {

	// Send request to service
	res, err := http.Post(testSidecarUrl+"/getcommand",
		"application/x-www-form-urlencoded",
		strings.NewReader(""))
	if err != nil {
		displayError("failed to post tests request to test sidecar server", err)
		return false
	}
	if res.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(res.Body)

		displayError(fmt.Sprintf("Received no \"200 OK\" from test sidecar server: %q",
			strings.TrimSuffix(string(b), "\n")))
		return false

	}

	// Read response body in JSON
	body, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		displayError("failed to read tests response from test sidecar server", err)
		return false
	}

	//displayGreenText(string(body))

	currentTestRun = TestRun{}

	err = json.Unmarshal(body, &currentTestRun)
	if err != nil {
		displayError("failed to unmarshall txestRun from test sidecar start", err)
		return false
	}

	return true
}

func sendTestCommandResult(status string) bool {

	// Send request to service
	res, err := http.Post(testSidecarUrl+"/putresult",
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("%s %s %s %s %s",
			currentTestRun.ID, currentTestRun.Name, name,
			status, strings.TrimSpace(strings.TrimLeft(currentTestRun.Commands[0], name)))))
	if err != nil {
		displayError("failed to post tests request to test sidecar server", err)
		return false
	}
	if res.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(res.Body)

		displayError(fmt.Sprintf("Received no \"200 OK\" from test sidecar server: %q",
			strings.TrimSuffix(string(b), "\n")))
		return false
	}

	//// Read response body in JSON
	//body, err := ioutil.ReadAll(res.Body)
	//res.Body.Close()
	//if err != nil {
	//	displayError("failed to read tests response from test sidecar server", err)
	//	return false
	//}

	//displayGreenText(string(body))

	return true
}

func sendTestFilterEvent(source string, event string) {

	currentTestEventFilter.Source = source
	currentTestEventFilter.Event = event

	testEventFilterJson, err := json.MarshalIndent(currentTestEventFilter, "", "  ")
	if err != nil {
		displayError("failed to marshall test event", err)
		return
	}

	// Send request to service
	res, err := http.Post(testSidecarUrl+"/putevent",
		"application/json",
		strings.NewReader(string(testEventFilterJson)))
	if err != nil {
		displayError("failed to put sidecar event", err)
	}
	if res.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(res.Body)

		displayError(fmt.Sprintf("Received no \"200 OK\" from put sidecar event: %q",
			strings.TrimSuffix(string(b), "\n")))

	}

	//// Read response body in JSON
	//body, err := ioutil.ReadAll(res.Body)
	//res.Body.Close()
	//if err != nil {
	//	displayError("failed to read sendTestFilterEvent response from test sidecar server", err)
	//}

	//displayGreenText(string(decodeJsonBytes(body)))
}

func prepareTestSummary() bool {

	// Send request to service
	res, err := http.Post(testSidecarUrl+"/preparesummary",
		"application/json",
		strings.NewReader(name))
	if err != nil {
		displayError("failed to inform sidecar about end of test", err)
	}
	if res.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(res.Body)

		displayError(fmt.Sprintf("Received no \"200 OK\" from put sidecar event: %q",
			strings.TrimSuffix(string(b), "\n")))

	}

	//// Read response body in JSON
	//body, err := ioutil.ReadAll(res.Body)
	//res.Body.Close()
	//if err != nil {
	//	displayError("failed to read sendTestFilterEvent response from test sidecar server", err)
	//}
	//
	//displayGreenText(string(decodeJsonBytes(body)))

	return true
}

func callTestSummary() bool {

	// Send request to service
	res, err := http.Post(testSidecarUrl+"/getsummary",
		"application/json",
		strings.NewReader(""))
	if err != nil {
		displayError("failed to inform sidecar about end of test", err)
	}
	if res.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(res.Body)

		displayError(fmt.Sprintf("Received no \"200 OK\" from put sidecar event: %q",
			strings.TrimSuffix(string(b), "\n")))

	}

	// Read response body in JSON
	body, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		displayError("failed to read sendTestFilterEvent response from test sidecar server", err)
		return false
	}

	//displayGreenText(string(decodeJsonBytes(body)))

	err = json.Unmarshal(body, &currentTestSummaries)
	if err != nil {
		displayError("failed to unmarshall currentTestSummaries from test sidecar", err)
		return false
	}

	return true
}

func callTestEvents() bool {

	// Send request to service
	res, err := http.Post(testSidecarUrl+"/getevents",
		"application/json",
		strings.NewReader(""))
	if err != nil {
		displayError("failed to get test events from sidecar", err)
		return false
	}
	if res.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(res.Body)

		displayError(fmt.Sprintf("Received no \"200 OK\" from get sidecar events: %q",
			strings.TrimSuffix(string(b), "\n")))
		return false
	}

	// Read response body in JSON
	body, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		displayError("failed to read callTestEvents response from test sidecar server", err)
		return false
	}

	//displayGreenText(string(decodeJsonBytes(body)))

	err = json.Unmarshal(body, &currentTestEventFilters)
	if err != nil {
		displayError("failed to unmarshall currentTestEventFilters from test sidecar", err)
		return false
	}

	return true
}

// executeTestCommand sends a peer name to get next line, if appropriate
// to the test-sidecar-server running under http://localhost:8081
func executeTestCommand() bool {

	if len(currentTestRun.Commands) == 0 {
		logYellow("empty test run queue")
		testend = true
		return false
	}
	if strings.Split(currentTestRun.Commands[0], " ")[0] != name {
		logYellow("command for another peer")
		return false
	}

	input := "/" + strings.TrimSpace(strings.TrimLeft(currentTestRun.Commands[0], name))
	displayText(strings.Trim(fmt.Sprintf("%s%v\n", prompt, input), "\n"))

	//logBlue("??? executeCommand ???")
	logBlue(strings.TrimLeft(input, "/"))

	if executeCommand(strings.TrimLeft(input, "/")) {
		return true
	}
	return false
}

// removeTest sends a peer name to remove the next line, if appropriate
// to the test-sidecar-server running under http://localhost:8081
func removeTest() bool {

	// Send request to service
	res, err := http.Post(testSidecarUrl+"/removecommand",
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("%s",
			name)))
	if err != nil {
		displayError("failed to load sidecar test get", err)
		return false
	}
	if res.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(res.Body)

		displayError(fmt.Sprintf("Received no \"200 OK\" from test sidecar get: %q",
			strings.TrimSuffix(string(b), "\n")))
		return false

	}
	return true
}
