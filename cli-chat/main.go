package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

var (
	name   string
	prompt string

	err error

	debug         *bool
	debugfilename *string
	initfilename  *string

	testMode  *bool
	testStart *bool
)

func main() {

	// todo -nolog and refactor debug -> log

	// test switches on testing
	testMode = flag.Bool("test", false, "switches on test mode")

	// start test
	testStart = flag.Bool("start", false, "last peer starts the test run")

	// debug switches on debugging
	debug = flag.Bool("debug", true, "switches on debugging")

	// debugfilename is the file to write debugging output
	debugfilename = flag.String("debugfile", "", "file to write debugging output to, use /dev/null to suppress debugging")

	// Parse input and check arguments
	flag.Parse()
	if flag.NArg() < 1 {
		_, _ = fmt.Fprintln(os.Stderr, "missing or wrong parameter: <name>")
		// todo usage
		os.Exit(1)
	}
	name = flag.Arg(0)
	prompt = fmt.Sprintf("<%s> ", name)

	// Start debugging to file, if switched on or filename specified
	if *debug || len(*debugfilename) > 0 {

		debugfile, err := startLogging(*debugfilename)
		if err != nil {
			panic(err)
		}
		defer debugfile.Close()

		// Current logfilename
		fmt.Printf("Start logging to %q\n", logfilename)

		// First entry in the logfile
		log.Printf("Session starting\n")
	}

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

		cmdLogging(strings.Split("on testqueue.log", " "))

		//logBlue(fmt.Sprintf("Empty: %v\n", currentTestRun.Queue))

		if !initTest("testqueue") {
			displayRedText("test load failed")
		}
		//logBlue(fmt.Sprintf("Init: %v\n", currentTestRun.Queue))

		for !testend {

			if executeTestCommand() {

				displayYelloText("command execution successful")

				sendTestCommandResult("OK")

				if removeTest() {
					displayYelloText("command remove successful")
				} else {
					displayRedText("command remove error")
				}
			} else {
				displayRedText("command execution failed")

			}

			//logBlue(fmt.Sprintf("After: %v\n", currentTestRun.Queue))

			if !refreshTest() {
				displayRedText("command load failed")
				break
			}

			if len(currentTestRun.Commands) != 0 && strings.Split(currentTestRun.Commands[0], " ")[0] == name {
				continue
			}

			//logBlue(fmt.Sprintf("Commands: %v\n", currentTestRun.Queue))

			time.Sleep(time.Second * 1)
		}

		displayGreenText(fmt.Sprintf("Test %s finished", currentTestRun.ID))

		if !prepareTestSummary() {
			displayRedText("test summary preparation failed")
		}

		cmdLogging([]string{"off"})
	}

	select {}
}
