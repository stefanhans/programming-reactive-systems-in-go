package main

import (
	"flag"
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	name   string
	prompt string

	err error

	logfile *string

	testMode     *bool
	testfilename *string
)

func checkCommandlineArgs() {

	// test switches on testing
	testMode = flag.Bool("test", false,
		"switches on test mode")

	// testfilename is the file in TEST_DIR to read test commands
	testfilename = flag.String("testfile", "",
		"file to load test commands")

	// logfile is the file to write loggging output
	logfile = flag.String("logfile", "",
		"file to write logging output to; use /dev/null to suppress logging")

	// Parse input and check arguments
	flag.Parse()
	if flag.NArg() < 1 {
		_, _ = fmt.Fprintln(os.Stderr,
			fmt.Sprintf("usage: "+
				"\t ./cli-chat [-test [-testfile=<filename>]] "+
				"[-logfile=<filename> | -logfile=/dev/null] <name>"))
		os.Exit(1)
	}

	name = flag.Arg(0)
	prompt = fmt.Sprintf("<%s> ", name)

}

func doTesting() {

	// http://localhost:8081
	testSidecarUrl = os.Getenv("TEST_SIDECAR_SERVER")
	_, err := url.ParseRequestURI(testSidecarUrl)
	if err != nil {
		fmt.Printf("environment variable TEST_SIDECAR_SERVER is not a valid URL")
		return
	}

	if *testfilename != "" {
		testName = strings.Split(*testfilename, ".")[0]
	} else {
		// Todo: testing with multiple tests
		testName = "testqueue"
	}

	cmdLogging(strings.Split(fmt.Sprintf("on %s.log", testName), " "))

	if !initTest(testName) {
		displayError("test load failed")
	}

	for !testend {

		if executeTestCommand() {

			logGreen("command execution successful")

			sendTestCommandResult("OK")

			if removeTest() {
				logGreen("command remove successful")
			} else {
				logYellow("command remove error")
			}
		} else {
			logYellow("command execution failed")
		}

		if !refreshTest() {
			displayError("command load failed")
			break
		}

		if len(currentTestRun.Commands) != 0 && strings.Split(currentTestRun.Commands[0], " ")[0] == name {
			continue
		}

		time.Sleep(time.Second * 1)
	}

	displayGreenText(fmt.Sprintf("Test %s finished", currentTestRun.ID))

	if !prepareTestSummary() {
		displayError("test summary preparation failed")
	}

	cmdLogging([]string{"off"})
}

func main() {

	checkCommandlineArgs()

	// Start logging to specified or default logfile
	file, err := startLogging(*logfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Current logfilename
	fmt.Printf("Start logging to %q\n", logfilename)

	// First entry in the logfile
	log.Printf("Session starting\n")

	// Create the TUI
	clientGui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalf("could not create tui: %v\n", err)
	}
	defer clientGui.Close()

	// Initialize help
	helpInit()

	// Start text-based UI
	go func() {
		err = runTUI()
		if err != nil {
			log.Fatalf("runTUI: %v", err)
		}
		os.Exit(0)
	}()

	if *testMode {
		doTesting()
	} else {
		if !executeCommand("init") {
			displayError("could not initialize chat")
		}
	}

	select {}
}
