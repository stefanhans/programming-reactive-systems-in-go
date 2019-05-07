package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

var (
	clientGui *gocui.Gui
)

// Configure and run the text-based UI
func runTUI() error {

	// Set function to manage all views and keybindings
	clientGui.SetManagerFunc(layout)

	// Bind keys with functions
	clientGui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	clientGui.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, send)

	// Start main event loop of the TUI
	return clientGui.MainLoop()
}

// Content to be displayed in the TUI
func layout(g *gocui.Gui) error {

	// Get rid of warnings
	_ = g

	// Get size of the terminal
	maxX, maxY := clientGui.Size()

	// Creates view "messages"
	messages, err := clientGui.SetView("messages", 0, 0, maxX-1, maxY-3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		messages.Autoscroll = true
		messages.Wrap = true
	}

	// Creates view "input"
	input, err := clientGui.SetView("input", 0, maxY-4, maxX-1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		input.Wrap = true
		input.Editable = true
	}

	// Set view "input" as the current view with focus and cursor
	_, err = clientGui.SetCurrentView("input")
	if err != nil {
		return err
	}

	// Show cursor
	clientGui.Cursor = true

	return nil
}

// Quit the TUI
func quit(g *gocui.Gui, v *gocui.View) error {

	// Get rid of warnings
	_, _ = g, v

	// Unsubscribe via PublisherClient
	//Unsubscribe()

	// Last entry in the individual log file
	log.Println("Session closing")

	// Last entry in the common log file
	//cLog.Printf("Session closing - details in %q\n", logfilename)

	return gocui.ErrQuit
}

// Send content from the view "input"
func send(g *gocui.Gui, inputView *gocui.View) error {

	// Get rid of warnings
	_ = g

	input := strings.Trim(inputView.Buffer(), "\n")

	switch {
	case strings.HasPrefix(input, "/"):

		// Interpret "input" as command
		displayText(strings.Trim(fmt.Sprintf("%s%v\n", prompt, input), "\n"))

		executeCommand(strings.TrimLeft(input, "/"))
	case strings.TrimSpace(input) == "":
		displayText(prompt)
	default:
		sendMessage(strings.Split(input, " "))
		displayText(strings.Trim(fmt.Sprintf("%s%v\n", prompt, input), "\n"))

	}

	// Clear the "input" and reset the cursor
	inputView.Clear()
	if err := inputView.SetCursor(0, 0); err != nil {
		log.Fatal(err)
	}
	return nil
}

func displayColoredMessages(msg string) {
	//log.Printf("msg: %v\n", msg)

	switch {
	case strings.Fields(msg)[0] == "<info>" ||
		strings.Fields(msg)[0] == "<joined>":
		displayGreenText(strings.Trim(msg, "\n"))
	case strings.Fields(msg)[0] == "<warn>" ||
		strings.Fields(msg)[0] == "<left>":
		displayYelloText(strings.Trim(msg, "\n"))
	case strings.Fields(msg)[0] == "<error>":
		displayRedText(strings.Trim(msg, "\n"))
	default:
		displayText(strings.Trim(fmt.Sprintf("\n%s", msg), "\n"))
	}
}

// Display text in "messages"
func displayText(txt string) error {

	// Update the "messages" view as soon as possible
	clientGui.Update(func(g *gocui.Gui) error {
		messagesView, err := clientGui.View("messages")
		if err != nil {
			return fmt.Errorf("could not display text: %v\n", err)
		}
		fmt.Fprintln(messagesView, txt)
		logGreen(txt)

		if *testMode {
			if sourceFilters, ok := testSourceFilters["messagesView"]; ok {

				for _, sourceFilter := range sourceFilters {
					filterEvent("messagesView", &sourceFilter, txt)
				}
			}
		}
		return nil
	})
	return nil
}

func displayGreenText(txt string) error {

	// Update the "messages" view as soon as possible
	clientGui.Update(func(g *gocui.Gui) error {
		messagesView, err := clientGui.View("messages")
		if err != nil {
			return fmt.Errorf("could not display text: %v\n", err)
		}

		fmt.Fprintln(messagesView, strings.Trim(fmt.Sprintf("\033[3%d;%dm%s\033[0m", 2, 1, txt), "\n"))
		logGreen(txt)
		return nil
	})
	return nil
}

func displayYelloText(txt string) error {

	// Update the "messages" view as soon as possible
	clientGui.Update(func(g *gocui.Gui) error {
		messagesView, err := clientGui.View("messages")
		if err != nil {
			return fmt.Errorf("could not display text: %v\n", err)
		}

		fmt.Fprintln(messagesView, strings.Trim(fmt.Sprintf("\033[3%d;%dm%s\033[0m", 3, 1, txt), "\n"))
		logYellow(txt)
		return nil
	})
	return nil
}

func displayRedText(txt string) error {

	// Update the "messages" view as soon as possible
	clientGui.Update(func(g *gocui.Gui) error {
		messagesView, err := clientGui.View("messages")
		if err != nil {
			return fmt.Errorf("could not display text: %v\n", err)
		}

		fmt.Fprintln(messagesView, strings.Trim(fmt.Sprintf("\033[3%d;%dm%s\033[0m", 1, 1, txt), "\n"))
		logRed(txt)
		return nil
	})
	return nil
}

// Send error to logfile and "messages"; txt is the text before ":"
func displayError(txt string, err ...error) error {

	var errorStr string

	if len(err) == 0 {
		errorStr = fmt.Sprintf("%s\n", txt)
	} else {
		errorStr = fmt.Sprintf("%s: %v\n", txt, err)
	}

	out := displayRedText(strings.Trim(errorStr, "\n"))
	displayText(prompt)

	return out
}

func displayJson(json []byte) error {
	return displayText(strings.Trim(fmt.Sprintf("%s\n%s", json,
		prompt), "\n"))
}
