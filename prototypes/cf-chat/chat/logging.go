package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	logfilename string
	logfile     *os.File

	cLog              *log.Logger
	commonLogfilename string
	commonLogfile     *os.File
)

func startLogging(name string) (*os.File, error) {

	// Config logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("DEBUG: ")

	// Prepare logfile for logging
	year, month, day := time.Now().Date()
	hour, minute, second := time.Now().Clock()
	logfilename = fmt.Sprintf("rudimentary-chat-tcp-%s-%v%02d%02d%02d%02d%02d.log", name,
		year, int(month), int(day), int(hour), int(minute), int(second))

	logfile, err := os.OpenFile(logfilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening logfile %v: %v", logfilename, err)
	}

	// Switch logging to logfile
	log.SetOutput(logfile)

	// First entry in the individual log file
	log.Printf("Session starting\n")

	return logfile, nil
}

func startCommonLogging(name, uuid string) (*log.Logger, *os.File, error) {

	id := fmt.Sprintf("%s %s: ", strings.Join(strings.Split(uuid, "")[:6], ""), name)
	cLog := log.New(os.Stderr, id, log.Ldate|log.Ltime|log.Lshortfile)

	// Prepare logfile for logging
	year, month, day := time.Now().Date()
	commonLogfilename = fmt.Sprintf("rudimentary-chat-tcp-%v%02d%02d.log",
		year, int(month), int(day))

	logfile, err := os.OpenFile(commonLogfilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening logfile %v: %v", logfilename, err)
	}

	// Switch logging to logfile
	cLog.SetOutput(logfile)

	// First entry in the common log file
	cLog.Printf("Session starting - details in %q\n", logfilename)

	return cLog, logfile, nil
}
