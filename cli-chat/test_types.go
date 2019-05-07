package main

// TestRun is the struct for the multi peer test
type TestRun struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Commands    []string `json:"queue,omitempty"`
	CurrentPeer string   `json:"currentpeer,omitempty"`
	Status      string   `json:"status,omitempty"`
}

type TestSourceFilter struct {
	Filter            string `json:"filter,omitempty"`
	NumExpectedEvents uint8  `json:"num_events,omitempty"`
}

type CommandResult struct {
	Peer    string `json:"peer,omitempty"`
	Status  string `json:"status,omitempty"`
	Comment string `json:"comment,omitempty"`
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

type TestSummary struct {
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
	currentTestRun    = TestRun{}
	currentTestResult = TestResult{}

	currentTestEventFilter  = TestEventFilter{}
	currentTestEventFilters = make([]TestEventFilter, 0)

	currentTestSummary   = TestSummary{}
	currentTestSummaries = make([]TestSummary, 0)

	currentSourceFilter = TestSourceFilter{}
	testSourceFilters   = map[string][]TestSourceFilter{
		"messagesView": make([]TestSourceFilter, 0),
	}

	testCommand string
	testend     = false

	// Todo: Use a environment variable
	TestUrl string = "http://localhost:8081"
)
