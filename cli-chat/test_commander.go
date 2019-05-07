package main

import (
	"fmt"
	"strconv"
	"strings"
)

func filterEvent(source string, sourceFilter *TestSourceFilter, event string) {

	for key, value := range testSourceFilters {

		logBlue(fmt.Sprintf("filterEvent: %s: %v\n", key, value))
	}

	if source == "messagesView" &&
		strings.TrimPrefix(event, sourceFilter.Filter) != event {

		logBlue("TESTFILTER: " + "-----------")
		logBlue("ID: " + currentTestEventFilter.ID)
		logBlue("Name: " + currentTestEventFilter.Name)
		logBlue("Peer: " + name)
		logBlue("Source: " + source)
		logBlue("Filter: " + sourceFilter.Filter)
		logBlue("Event: " + event)
		logBlue("TESTFILTER: " + "-----------")

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

		case "testfilter":
			return addTestFilter(commandFields[1:])

		case "testsummary":
			return showTestSummary(commandFields[1:])

		case "testevents":
			return showTestEvents(commandFields[1:])

		case "testfilters":
			return showTestSourceFilters(commandFields[1:])

		default:
			return noCommand(commandFields)
		}
	}
	return false
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

	logBlue(fmt.Sprintf("testSourceFilters: %v\n", testSourceFilters))

	return true
}

func showTestSummary(arguments []string) bool {

	//if len(arguments) < 2 {
	//	displayError("not enough arguments defined to addTestFilter")
	//	return false
	//}

	currentTestSummaries = make([]TestSummary, 0)

	if !callTestSummary() {
		displayRedText("test summary failed")
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

		logBlue(fmt.Sprintf("test: %v\n", test))

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

		/*

			<alice> 2019/05/03 11:25:00 logging.go:61: test 7: {3d9a6e86-6d85-11e9-9c25-8c8590a2bb7c testqueue bob OK command  testfilter messagesView <alice> test }

			<alice> 2019/05/03 11:25:00 logging.go:61: test 8: {3d9a6e86-6d85-11e9-9c25-8c8590a2bb7c testqueue bob OK event  testfilter messagesView <alice> test }
		*/

	}

	out += fmt.Sprintf("-------------------------------------")
	out += fmt.Sprintf("-------------------------------------\n")

	displayText(out)

	return true
}

//

func showTestEvents(arguments []string) bool {

	//if len(arguments) < 2 {
	//	displayError("not enough arguments defined to addTestFilter")
	//	return false
	//}

	if !callTestEvents() {
		displayRedText("call test events failed")
		return false
	}

	if len(currentTestEventFilters) == 0 {
		displayError("no test event filters available")
		return false
	}

	testID := currentTestEventFilters[0].ID
	testName := currentTestEventFilters[0].Name

	out := fmt.Sprintf("Event filter of %q (%s)\n", testName, testID)
	out += fmt.Sprintf("-------------------------------------")
	out += fmt.Sprintf("-------------------------------------\n")

	for _, filter := range currentTestEventFilters {

		logBlue(fmt.Sprintf("filter: %v\n", filter))

		if filter.ID != testID || filter.Name != testName {
			displayError("event filter inconsistent")
			return false
		}

		out += fmt.Sprintf("%s %s: Filter: %q Event: %q\n",
			filter.Peer, filter.Source, filter.Filter, filter.Event)

		/*filter: {ef87a344-6dbb-11e9-855d-8c8590a2bb7c testqueue bob messagesView <alice> test <alice> test}*/

	}

	out += fmt.Sprintf("-------------------------------------")
	out += fmt.Sprintf("-------------------------------------\n")

	displayText(out)

	return true
}

func showTestSourceFilters(arguments []string) bool {

	//if len(arguments) < 2 {
	//	displayError("not enough arguments defined to addTestFilter")
	//	return false
	//}

	if len(testSourceFilters) == 0 {
		displayError("no test event sources available")
		return false
	}

	out := fmt.Sprintf("Source filters\n")
	out += fmt.Sprintf("-------------------------------------")
	out += fmt.Sprintf("-------------------------------------\n")

	for source, filters := range testSourceFilters {

		logBlue(fmt.Sprintf("%s: %v\n", source, filters))

		filterStr := ""
		for _, filter := range filters {
			filterStr += fmt.Sprintf("%q ", filter)
		}

		out += fmt.Sprintf("%s %s\n", source, filterStr)
	}

	out += fmt.Sprintf("-------------------------------------")
	out += fmt.Sprintf("-------------------------------------\n")

	displayText(out)

	return true
}
