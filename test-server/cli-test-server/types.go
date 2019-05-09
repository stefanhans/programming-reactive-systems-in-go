package server

import "sync"

var (
	mtx sync.RWMutex
)

// TestRun is the struct for the multi peer test
type TestRun struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Commands    []string `json:"queue,omitempty"`
	CurrentPeer string   `json:"currentpeer,omitempty"`
	Status      string   `json:"status,omitempty"`
}

// CommandResult is used to receive the result of a command and
// will be stored in the current test result.
type CommandResult struct {
	Peer   string `json:"peer,omitempty"`
	Status string `json:"status,omitempty"`
	Data   string `json:"data"`
}

// TestResult stores the results of a test run.
type TestResult struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	CommandResults []*CommandResult
}

// TestEventFilter stores information about defined filters
// and their matching events.
type TestEventFilter struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	Peer              string `json:"peer,omitempty"`
	Source            string `json:"source,omitempty"`
	Filter            string `json:"filter,omitempty"`
	Event             string `json:"event,omitempty"`
	NumExpectedEvents uint8  `json:"num_expected_events,omitempty"`
	NumReceivedEvents uint8  `json:"num_received_events"`
}

// TestEvaluation stores the evaluation of one command or event of the test run.
// An array of these test evaluations is the current test summary.
type TestEvaluation struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Peer    string `json:"peer,omitempty"`
	Status  string `json:"status,omitempty"`
	Kind    string `json:"kind,omitempty"`
	Test    string `json:"test,omitempty"`
	Result  string `json:"result,omitempty"`
	Comment string `json:"comment,omitempty"`
}

var (
	currentTestRun = TestRun{}

	currentTestEventFilter  = TestEventFilter{}
	currentTestEventFilters = make([]TestEventFilter, 0)

	currentTestResult = TestResult{}

	currentTestEvaluation = TestEvaluation{}
	currentTestSummary    = make([]TestEvaluation, 0)
)
