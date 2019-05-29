package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/peterh/liner"
)

// Define a set of commands
var commands = []string{"ls", "pwd", "date", "hostname", "quit"}

func prompt() string {
	return fmt.Sprintf("<%s %s> ", time.Now().Format("Jan 2 15:04:05.000"), "user")
}

func executeShellScript(cmdline string) bool {

	if cmdline == "quit" {
		os.Exit(0)
	}

	words := strings.Fields(cmdline)

	if len(words) == 0 {
		fmt.Print("")
		return false
	}

	binary, lookErr := exec.LookPath(words[0])
	if lookErr != nil {
		fmt.Printf("%v\n", lookErr)
		return false
	}

	var cmd *exec.Cmd
	if len(words) == 1 {
		cmd = exec.Command(binary)
	} else {
		fmt.Printf("command with arguments not supported\n")
		return false
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		fmt.Printf("could not run: %v\n", err)
		return false
	}

	fmt.Printf(out.String())

	return true
}

func main() {
	state := liner.NewLiner()
	state.SetTabCompletionStyle(liner.TabPrints)
	state.SetCompleter(func(line string) (ret []string) {
		for _, c := range commands {
			if strings.HasPrefix(c, line) {
				ret = append(ret, c)
			}
		}
		return
	})
	defer state.Close()

	for {
		p, err := state.Prompt(prompt())
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		if executeShellScript(p) {
			state.AppendHistory(p)
		}
	}
}
