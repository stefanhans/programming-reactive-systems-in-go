package main

var (
	// Directory for the test data
	testDir, currentTestDir string

	// Portnumber the service is listening
	testPort string
)

// TestRun is the struct for the multi peer test
type TestRun struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Commands    []string `json:"queue,omitempty"`
	CurrentPeer string   `json:"currentpeer,omitempty"`
	Status      string   `json:"status,omitempty"`
}

type CommandResult struct {
	Peer   string `json:"peer,omitempty"`
	Status string `json:"status,omitempty"`
	Data   string `json:"data"`
}

type TestResult struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	CommandResults []*CommandResult
}

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
