package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func executeScript(arguments []string) {

	if len(arguments) == 0 {
		displayText(strings.Trim(fmt.Sprintf("no filename to execute specified\n%s",
			prompt), "\n"))
		return
	}

	b, err := ioutil.ReadFile(arguments[0])
	if err != nil {
		displayText(strings.Trim(fmt.Sprintf("could not read the file %q: %v\n%s",
			prompt, arguments[0], err), "\n"))
		return
	}

	prompt = fmt.Sprintf("<%s %q> ", name, arguments[0])

	lines := strings.Split(string(b), "\n")
	for i, line := range lines {
		log.Printf("EXECUTE %d: %q\n", i, line)
		if strings.TrimSpace(line) == "" ||
			strings.Split(strings.TrimSpace(line), "")[0] == "#" {
			continue
		}
		//echoScript(strings.Split(line, " "))

		displayText(prompt + line)
		executeCommand(line)

		//if _, ok := commands[strings.Split(strings.TrimLeft(line, "/"), " ")[0]]; ok {
		//	executeCommand(line)
		//} else {
		//	displayText(strings.Trim(fmt.Sprintf("%q is an unknown command\n%s", strings.Split(line, " ")[0],
		//		prompt), "\n"))
		//}
	}
	prompt = fmt.Sprintf("<%s> ", name)
	displayText(prompt)
}

func sleepScript(arguments []string) {

	var numSeconds int

	if len(arguments) == 0 {
		numSeconds = 1
	} else {
		numSeconds, err = strconv.Atoi(arguments[0])
	}

	time.Sleep(time.Second * time.Duration(numSeconds))
	displayText(prompt)
}

func echoScript(arguments []string) {

	displayText(strings.Trim(fmt.Sprintf("%s", strings.Join(arguments, " ")), "\n"))
}

func executeShellScript(arguments []string) {

	if len(arguments) == 0 {
		displayError("no shell script specified")
		//displayText(prompt)
		return
	}

	binary, lookErr := exec.LookPath(arguments[0])
	if lookErr != nil {
		displayError("not found", lookErr)
		return
	}

	var cmd *exec.Cmd
	if len(arguments) == 1 {
		cmd = exec.Command(binary)
	} else {
		switch len(arguments[1:]) {
		case 1:
			cmd = exec.Command(binary, arguments[1])
		case 2:
			cmd = exec.Command(binary, arguments[1], arguments[2])
		case 3:
			cmd = exec.Command(binary, arguments[1], arguments[2], arguments[3])
		case 4:
			cmd = exec.Command(binary, arguments[1], arguments[2], arguments[3], arguments[4])
		}
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		displayError("could not run", err)
		return
	}
	//fmt.Printf("in all caps: %q\n", out.String())
	displayText(strings.Trim(fmt.Sprintf("%s%s", out.String(), prompt), "\n"))

	//execErr := cmd.Run()
	//if execErr != nil {
	//	panic(execErr)
	//}
	//
	//out, err := cmd.Output()
	//if err != nil {
	//	displayError("output error", err)
	//}
	//
	//displayGreenText(string(out))
	//
	//syscall.Exec(arguments[0], []string{}, []string{})
}
