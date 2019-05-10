package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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

	// Config logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Prepare logfile for logging
	year, month, day := time.Now().Date()
	hour, minute, second := time.Now().Clock()
	logfilename := fmt.Sprintf("cli-test-server%v%02d%02d%02d%02d%02d.log",
		year, int(month), int(day), int(hour), int(minute), int(second))
	// Todo logdir as env variable
	logfile, err := os.OpenFile("log/"+logfilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening logfile %v: %v", logfilename, err)
		return
	}

	// Switch logging to logfile
	log.SetOutput(logfile)

	// TestRun related handler functions
	http.HandleFunc("/init", InitRun)
	http.HandleFunc("/getrun", GetRun)
	http.HandleFunc("/resetrun", ResetRun)

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
