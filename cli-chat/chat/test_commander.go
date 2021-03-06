package chat

import (
	"fmt"
	"strconv"
	"strings"
)

func filterEvent(source string, sourceFilter *TestSourceFilter, event string) {

	if source == "messagesView" &&
		strings.TrimPrefix(event, sourceFilter.Filter) != event {

		currentTestEventFilter.Peer = name
		currentTestEventFilter.Filter = sourceFilter.Filter
		currentTestEventFilter.NumExpectedEvents = sourceFilter.NumExpectedEvents

		sendTestFilterEvent(source, event)
	}
}

func executeTestEventCommand(commandFields []string) bool {

	// Check for empty string without prefix
	if len(commandFields) > 0 {

		// Switch according to the first word and call appropriate function with the rest as arguments
		switch commandFields[0] {

		// Todo complete the help for new test commands
		// testrun, testclose, testreset

		case "testrun":
			return getTestRun()

		case "testreset":
			return resetTestRun()

		case "testfilter":
			return addTestFilter(commandFields[1:])

		case "testlocalfilters":
			return showLocalTestFilters(commandFields[1:])

		case "testfilters":
			return showTestFilters(commandFields[1:])

		case "testsummary":
			return showTestSummary(commandFields[1:])

		case "testsummaryprepare":
			return prepareTestSummary()

		default:
			return noCommand(commandFields)
		}
	}
	return false
}

func getTestRun() bool {

	if !refreshTest() {
		displayError("test refresh failed")
		return false
	}

	testID := currentTestRun.ID
	testName := currentTestRun.Name

	out := fmt.Sprintf("Test Run %q (%s)\n", testName, testID)
	out += fmt.Sprintf("-------------------------------------")
	out += fmt.Sprintf("-------------------------------------\n")

	for _, testCommand := range currentTestRun.Commands {
		out += fmt.Sprintf("%q\n", testCommand)
	}

	out += fmt.Sprintf("-------------------------------------")
	out += fmt.Sprintf("-------------------------------------\n")

	displayText(out)

	return true
}

func addTestFilter(arguments []string) bool {

	if len(arguments) < 3 {
		displayError("not enough arguments defined to addTestFilter")
		return false
	}

	if _, ok := testSourceFilters[arguments[0]]; !ok {
		displayError("unknown source to addTestFilter")
		return false
	}

	numExpectedEvents, err := strconv.Atoi(arguments[1])
	if err != nil {
		displayError(fmt.Sprintf("%q is not a number", arguments[1]))
		return false
	}

	currentSourceFilter.NumExpectedEvents = uint8(numExpectedEvents)
	currentSourceFilter.Filter = strings.Join(arguments[2:], " ")

	testSourceFilters[arguments[0]] = append(testSourceFilters[arguments[0]], currentSourceFilter)

	return true
}

func showTestSummary(arguments []string) bool {

	// Get rid off warning
	_ = arguments

	currentTestSummaries = make([]TestSummary, 0)

	if !callTestSummary() {
		displayError("test summary failed")
		return false
	}

	if len(currentTestSummaries) == 0 {
		displayError("no test summary available")
		return false
	}

	testID := currentTestSummaries[0].ID
	testName := currentTestSummaries[0].Name

	out := fmt.Sprintf("Summary of %q (%s)\n", testName, testID)
	out += fmt.Sprintf("-------------------------------------")
	out += fmt.Sprintf("-------------------------------------\n")

	for _, test := range currentTestSummaries {

		if test.ID != testID || test.Name != testName {
			displayError("summary inconsistent")
			return false
		}

		if test.Status == "OK" {
			out += fmt.Sprintf("%s\t %q %s\n",
				test.Kind,
				fmt.Sprintf("%s %s", test.Peer, test.Test),
				fmt.Sprintf("\033[3%d;%dm%s\033[0m", 2, 1, test.Status))
		} else {
			out += fmt.Sprintf("%s from %s: %q %s\n",
				test.Kind, test.Peer, test.Test,
				fmt.Sprintf("\033[3%d;%dm%s: %q\033[0m", 1, 1,
					test.Status, test.Result))

		}
	}

	out += fmt.Sprintf("-------------------------------------")
	out += fmt.Sprintf("-------------------------------------\n")

	displayText(out)

	return true
}

func showTestFilters(arguments []string) bool {

	// Get rid off warning
	_ = arguments

	if !callTestEvents() {
		displayError("call test events failed")
		return false
	}

	if len(currentTestEventFilters) == 0 {
		displayError("no test event filters available")
		return false
	}

	testID := currentTestEventFilters[0].ID
	testName := currentTestEventFilters[0].Name

	out := fmt.Sprintf("Event filters of %q (%s)\n", testName, testID)
	out += fmt.Sprintf("-------------------------------------")
	out += fmt.Sprintf("-------------------------------------\n")

	for _, filter := range currentTestEventFilters {

		if filter.ID != testID || filter.Name != testName {
			displayError("event filter inconsistent")
			return false
		}

		out += fmt.Sprintf("%s %s: Filter: %q (%d) Event: %q (%d)\n",
			filter.Peer, filter.Source,
			filter.Filter, filter.NumExpectedEvents,
			filter.Event, filter.NumReceivedEvents)
	}

	out += fmt.Sprintf("-------------------------------------")
	out += fmt.Sprintf("-------------------------------------\n")

	displayText(out)

	return true
}

func showLocalTestFilters(arguments []string) bool {

	// Get rid off warning
	_ = arguments

	if len(testSourceFilters) == 0 {
		displayError("no test event sources available")
		return false
	}

	out := fmt.Sprintf("Source filters\n")
	out += fmt.Sprintf("-------------------------------------")
	out += fmt.Sprintf("-------------------------------------\n")

	for source, filters := range testSourceFilters {

		filterStr := ""
		for _, filter := range filters {
			filterStr += fmt.Sprintf("\t%q (%d events expected)\n", filter.Filter, filter.NumExpectedEvents)
		}

		out += fmt.Sprintf("%s\n%s\n", source, filterStr)
	}

	out += fmt.Sprintf("-------------------------------------")
	out += fmt.Sprintf("-------------------------------------\n")

	displayText(out)

	return true
}
