package main

import (
	"log"
	"net/http"
	"sync"
)

var (
	mtx sync.RWMutex
)

func main() {

	// TestRun related
	http.HandleFunc("/init", InitRun)
	http.HandleFunc("/getrun", GetRun)

	// TestRun.Commands related
	http.HandleFunc("/getcommand", GetCommand)
	http.HandleFunc("/removecommand", RemoveCommand)

	// TestEventFilter related
	http.HandleFunc("/putevent", PutTestEvent)
	http.HandleFunc("/getevents", GetTestEvents)

	// TestResult related
	http.HandleFunc("/putresult", PutTestResult)
	http.HandleFunc("/getresults", GetTestResults)

	// TestSummary related
	http.HandleFunc("/preparesummary", PrepareTestSummary)
	http.HandleFunc("/getsummary", GetTestSummary)

	// TestEventFilter related - curl only, i.e. not used by cli-chat
	http.HandleFunc("/getfilters", GetTestSourceFilters)

	log.Fatal(http.ListenAndServe(":8081", nil))

}
