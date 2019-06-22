package chat

// TestRun is the struct for the multi peer test
type TestRun struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Commands    []string `json:"commands,omitempty"`
	CurrentPeer string   `json:"currentpeer,omitempty"`
	Status      string   `json:"status,omitempty"`
}

// TestSourceFilter is used to send filter to the service via 'testfilter'
type TestSourceFilter struct {
	Filter            string `json:"filter,omitempty"`
	NumExpectedEvents uint8  `json:"num_events,omitempty"`
}

// CommandResult is used to send the result of a command and
// will be stored in the current test result.
type CommandResult struct {
	Peer    string `json:"peer,omitempty"`
	Status  string `json:"status,omitempty"`
	Comment string `json:"comment,omitempty"`
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

// TestSummary stores the summary of a test run.
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
	currentTestRun = TestRun{}

	currentTestEventFilter  = TestEventFilter{}
	currentTestEventFilters = make([]TestEventFilter, 0)

	currentTestSummaries = make([]TestSummary, 0)

	currentSourceFilter = TestSourceFilter{}
	testSourceFilters   = map[string][]TestSourceFilter{
		"messagesView": make([]TestSourceFilter, 0),
	}

	testend = false

	testName       string
	testSidecarUrl string
)
