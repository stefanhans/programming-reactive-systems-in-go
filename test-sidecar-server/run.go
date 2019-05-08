package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/pborman/uuid"
)

func InitRun(w http.ResponseWriter, r *http.Request) {

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read request: %s", err), http.StatusInternalServerError)
		return
	}

	// Get arguments from response body
	arguments := strings.Split(string(body), " ")
	if len(arguments) < 1 {
		http.Error(w, fmt.Sprintf("error: request has not enough arguments"),
			http.StatusInternalServerError)
		return
	}

	if currentTestRun.ID == "" {
		currentTestRun.ID = uuid.NewUUID().String()
		currentTestRun.Name = arguments[0]

		b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.cmd", testDir, currentTestRun.Name))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read command file: %s", err),
				http.StatusInternalServerError)
		}

		fmt.Printf("b: %s\n", string(b))

		mtx.Lock()

		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				currentTestRun.Commands = append(currentTestRun.Commands, line)
			}
		}

		fmt.Printf("currentTestRun.Commands: %v\n", currentTestRun.Commands)
		mtx.Unlock()

		currentTestRun.Status = "READY"
	}

	// Marshal array of struct
	testRunJson, err := json.MarshalIndent(currentTestRun, "", " ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write marshal testRun: %s", err),
			http.StatusInternalServerError)
	}

	currentTestDir = fmt.Sprintf("%s/%s", testDir, currentTestRun.ID)
	err = os.MkdirAll(currentTestDir, os.ModePerm)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create current test directory: %s", err),
			http.StatusInternalServerError)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/testRun.json", currentTestDir), append(testRunJson, byte('\n')), 0600)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write testrun to test directory: %s", err),
			http.StatusInternalServerError)
	}

	_, err = fmt.Fprintf(w, string(testRunJson))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err),
			http.StatusInternalServerError)
	}
}

/*
curl -d "testqueue" http://localhost:8081/init
*/

func GetCommand(w http.ResponseWriter, r *http.Request) {

	// Get rid of warning
	_ = r

	// Marshal array of struct
	testRunJson, err := json.MarshalIndent(currentTestRun, "", " ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write marshal testRun: %s", err),
			http.StatusInternalServerError)
	}

	_, err = fmt.Fprintf(w, string(testRunJson))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err),
			http.StatusInternalServerError)
	}
}

/*
curl http://localhost:8081/get
*/

func RemoveCommand(w http.ResponseWriter, r *http.Request) {

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read request: %s", err), http.StatusInternalServerError)
		return
	}

	// Get arguments from response body
	arguments := strings.Split(string(body), " ")
	if len(arguments) < 1 {
		http.Error(w, fmt.Sprintf("error: request has not enough arguments"),
			http.StatusInternalServerError)
		return
	}

	mtx.Lock()

	if len(currentTestRun.Commands) == 0 {
		currentTestRun.Status = "EMPTY"
		mtx.Unlock()
		return
	}

	//fmt.Printf("Before rm %q: %v\n", arguments[0], currentTestRun.Commands)
	if strings.Split(currentTestRun.Commands[0], " ")[0] == arguments[0] {
		currentTestRun.Commands = currentTestRun.Commands[1:]
	}
	//fmt.Printf("After rm %q: %v\n", arguments[0], currentTestRun.Commands)

	mtx.Unlock()

	_, err = fmt.Fprintf(w, "")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err),
			http.StatusInternalServerError)
	}
}

func GetRun(w http.ResponseWriter, r *http.Request) {

	// Get rid of warning
	_ = r

	mtx.Lock()

	// Marshal array of struct
	testRunJson, err := json.MarshalIndent(currentTestRun, "", " ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write marshal currentTestRun: %s", err),
			http.StatusInternalServerError)
	}

	mtx.Unlock()

	_, err = fmt.Fprintf(w, string(testRunJson))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err),
			http.StatusInternalServerError)
	}
}

/*
curl http://localhost:8081/getrun
*/
