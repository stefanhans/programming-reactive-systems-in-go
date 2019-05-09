package chat

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	logfilename  string
	tmpDebugfile *os.File
)

func startLogging(logname string) (*os.File, error) {

	// Config logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if len(logname) == 0 {

		// Prepare logfile for logging
		year, month, day := time.Now().Date()
		hour, minute, second := time.Now().Clock()
		logfilename = fmt.Sprintf("cli-chat-%s-%v%02d%02d%02d%02d%02d.log", name,
			year, int(month), int(day), int(hour), int(minute), int(second))
	} else {
		logfilename = logname
		log.SetPrefix(fmt.Sprintf("<%s> ", name))
	}
	logfile, err := os.OpenFile(logfilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening logfile %v: %v", logfilename, err)
	}

	// Switch logging to logfile
	log.SetOutput(logfile)

	return logfile, nil
}

func logYellow(msg string) {

	log.Println(fmt.Sprintf("\033[3%d;%dm%s\033[0m", 3, 1, msg))
}

func logRed(msg string) {

	log.Println(fmt.Sprintf("\033[3%d;%dm%s\033[0m", 1, 1, msg))
}

func logGreen(msg string) {

	log.Println(fmt.Sprintf("\033[3%d;%dm%s\033[0m", 2, 1, msg))
}

func logBlue(msg string) {

	log.Println(fmt.Sprintf("\033[3%d;%dm%s\033[0m", 6, 1, msg))
}

func cmdLogging(arguments []string) {

	if len(arguments) == 0 ||
		(len(arguments) == 1 && arguments[0] != "off") {
		displayText(strings.Trim(fmt.Sprintf("wrong input. Usage: \n\t 'log (on <filename>) | off\n%s",
			prompt), "\n"))

		return
	}

	if arguments[0] == "on" && len(arguments) > 1 {
		displayText(strings.Trim(fmt.Sprintf("Switch to logging by command to %q\n%s", arguments[1],
			prompt), "\n"))
		tmpDebugfile, err = startLogging(arguments[1])
		if err != nil {
			displayError("could not start logging", err)
		} else {
			displayText(strings.Trim(fmt.Sprintf("Start logging by command to %q\n%s", arguments[1],
				prompt), "\n"))
		}

		return
	}

	if arguments[0] == "off" {
		log.Printf("Stop logging by command")
		_ = tmpDebugfile.Close()

		// Start logging to file
		_, err := startLogging(*logfile)
		if err != nil {
			displayError("could not start logging", err)
			return
		}
		log.Printf("Switch from logging by command to %q\n", logfile)
	}
}
