package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	// Directory for the test data
	testDir, currentTestDir string

	// Portnumber the service is listening
	testPort string
)

// Run starts the service.
func Run() {

	// The environment variable "TEST_DIR"
	// defines the directory for the test data
	testDir = os.Getenv("TEST_DIR")
	dir, err := os.Stat(testDir)
	if err != nil {
		fmt.Printf("TEST_DIR environment variable invalid: %v\n", err)
		return
	}
	if !dir.IsDir() {
		fmt.Printf("TEST_DIR is not a directory\n")
		return
	}

	// The environment variable "TEST_PORT"
	// defines the portnumber the test service is listening
	testPort = os.Getenv("TEST_PORT")
	if testPort == "" {
		fmt.Printf("TEST_PORT environment variable unset or missing\n")
		return
	}

	// TestRun related handler functions
	http.HandleFunc("/init", InitRun)
	http.HandleFunc("/getrun", GetRun)

	// TestRun.Commands related handler functions
	http.HandleFunc("/getcommand", GetCommand)
	http.HandleFunc("/removecommand", RemoveCommand)

	// TestEventFilter related handler functions
	http.HandleFunc("/putevent", PutTestEvent)
	http.HandleFunc("/getevents", GetTestEvents)

	// TestResult related handler functions
	http.HandleFunc("/putresult", PutTestResult)
	http.HandleFunc("/getresults", GetTestResults)

	// TestSummary related handler functions
	http.HandleFunc("/preparesummary", PrepareTestSummary)
	http.HandleFunc("/getsummary", GetTestSummary)

	// Start the service
	log.Fatal(http.ListenAndServe(":"+testPort, nil))
}
