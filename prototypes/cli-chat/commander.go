package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	commands    = make(map[string]string)
	commandKeys []string

	tmpDebugfile *os.File
)

func commandsInit() {
	commands = make(map[string]string)

	// Memberlist
	commands["memberlistconfigure"] = "memberlistconfigure \n\t memberlistconfigure creates a memberlist configuration\n"
	commands["showconfig"] = "showconfig \n\t showconfig shows the memberlist configuration\n"
	commands["saveconfig"] = "saveconfig [file] \n\t saveconfig saves the memberlist configuration as JSON file\n"
	commands["loadconfig"] = "loadconfig file \n\t loadconfig load the memberlist configuration from JSON file\n"

	commands["memberlistcreate"] = "memberlistcreate \n\t memberlistcreate creates the memberlist specified by the configuration\n"
	commands["showmemberlist"] = "showmemberlist \n\t showmemberlist shows the memberlist\n"
	commands["localnode"] = "localnode \n\t localnode shows the local node's name and address\n"

	commands["memberlist"] = "memberlist \n\t memberlist lists all members\n"
	commands["memberlistjoin"] = "memberlistjoin [<members> ...] \n\t memberlistjoin add oneself or other member(s) to the memberlist\n"
	commands["memberlistleave"] = "memberlistleave [<timeout in seconds, default: 1 sec>] \n\t memberlistleave broadcasts leave message until finished or timeout is reached\n"
	commands["memberlistupdate"] = "memberlistupdate [<timeout in seconds, default: 1 sec>] \n\t memberlistupdate broadcasts re-advertising the local node message until finished or timeout is reached\n"

	commands["memberliststart"] = "memberliststart \n\t memberliststart starts broadcasting to the members\n"
	commands["memberlistshutdown"] = "memberlistshutdown \n\t memberlistshutdown stops broadcasting to the members\n"
	commands["memberlistshutdowntransport"] = "memberlistshutdowntransport \n\t memberlistshutdowntransport stops broadcasting transport to the members\n"
	commands["memberlisthealthscore"] = "memberlisthealthscore \n\t memberlisthealthscore shows the health score >= 0, lower numbers are better\n"
	commands["memberlistdelete"] = "memberlistdelete \n\t memberlistdelete sets memberlist = nil\n"

	commands["bootstrapjoin"] = "bootstrapjoin \n\t bootstrapjoin joins calling peer to bootstrap peers\n"
	commands["bootstrapleave"] = "bootstrapleave \n\t bootstrapleave leave calling peer from bootstrap peers\n"
	commands["bootstraprefill"] = "bootstraprefill \n\t bootstraprefill refill bootstrap peers with calling peer\n"
	commands["bootstraplist"] = "bootstraplist \n\t bootstraplist  list bootstrap peers from remote\n"
	commands["bootstrapreset"] = "bootstrapreset \n\t bootstrap joins peer to bootstrap peers\n"
	commands["bootstraplistlocal"] = "bootstraplistlocal \n\t bootstraplistlocal list bootstrap peers from local map\n"

	commands["broadcastadd"] = "broadcastadd <key> <message> \n\t broadcastadd updates a key/message at all members\n"
	commands["broadcastdel"] = "broadcastdel <key> <message> \n\t broadcastdel deletes a key at all members\n"
	commands["broadcastlist"] = "broadcastlist \n\t broadcastlist lists all local key/value pairs\n"

	// Chat
	commands["chatjoin"] = "chatjoin  \n\t chatjoin starts chat listener and broadcasts new chat member\n"
	commands["chatleave"] = "chatleave  \n\t chatleave broadcasts deletion of this chat member\n"
	commands["chatmemberlist"] = "chatmemberlist  \n\t chatmemberlist lists all chat members\n"
	commands["chatmemberping"] = "chatmemberping  \n\t chatmemberping ping a chat member\n"

	commands["ping"] = "Experimental: ping <chatmember>  \n\t ping pings a member of the chat\n"
	commands["msg"] = "msg  \n\t msg sends the rest of the line as message to all other members\n"

	// Script
	commands["execute"] = "execute file \n\t execute execute the commands in the file line by line, '#' is comment\n"
	commands["sleep"] = "sleep seconds \n\t sleep sleeps for seconds\n"
	commands["echo"] = "echo text_w/o_linebreak \n\t echo prints rest of line\n"

	// Internals
	commands["log"] = "log (on <filename>)|off \n\t log starts or stops writing logging output in the specified file\n"

	commands["quit"] = "quit  \n\t close the session and exit\n"

	// Developer
	commands["play"] = "play  \n\t for developer playing\n"

	// To store the keys in sorted order
	for commandKey := range commands {
		commandKeys = append(commandKeys, commandKey)
	}
	sort.Strings(commandKeys)
}

// Execute a command specified by the argument string
func executeCommand(commandline string) bool {

	// Trim prefix and split string by white spaces
	commandFields := strings.Fields(commandline)

	// Check for empty string without prefix
	if len(commandFields) > 0 {

		// Switch according to the first word and call appropriate function with the rest as arguments
		switch commandFields[0] {

		case "memberlistconfigure":
			configureMemberlist(commandFields[1:])
			return true

		case "showconfig":
			showMemberlistConfiguration(commandFields[1:])
			return true

		case "saveconfig":
			saveMemberlistConfiguration(commandFields[1:])
			return true

		case "loadconfig":
			loadMemberlistConfiguration(commandFields[1:])
			return true

		case "memberlistcreate":
			createMemberlist(commandFields[1:])
			return true

		case "showmemberlist":
			showMemberlist(commandFields[1:])
			return true

		case "localnode":
			showLocalNode(commandFields[1:])
			return true

		case "memberlist":
			listMembers(commandFields[1:])
			return true

		case "memberlistjoin":
			joinMemberlist(commandFields[1:])
			return true

		case "memberlistleave":
			leaveMemberlist(commandFields[1:])
			return true

		case "memberlistupdate":
			updateMemberlist(commandFields[1:])
			return true

		case "memberliststart":
			startBroadcast(commandFields[1:])
			return true

		case "memberlistshutdown":
			shutdownBroadcast(commandFields[1:])
			return true

		case "memberlistshutdowntransport":
			shutdownBroadcastTransport(commandFields[1:])
			return true

		case "memberlisthealthscore":
			getHealthScore(commandFields[1:])
			return true

		case "memberlistdelete":
			deleteMemberlist(commandFields[1:])
			return true

		case "bootstrapjoin":
			joinBootstrapPeers(commandFields[1:])
			return true

		case "bootstrapleave":
			leaveBootstrapPeers(commandFields[1:])
			return true

		case "bootstraprefill":
			refillBootstrapPeers(commandFields[1:])
			return true

		case "bootstraplist":
			listBootstrapPeers(commandFields[1:])
			return true

		case "bootstrapreset":
			resetBootstrapPeers(commandFields[1:])
			return true

		case "bootstraplistlocal":
			listLocalBootstrapPeers(commandFields[1:])
			return true

		case "broadcastadd":
			broadcastAddMessage(commandFields[1:])
			return true

		case "broadcastdel":
			broadcastDelMessage(commandFields[1:])
			return true

		case "chatjoin":
			listenStream(commandFields[1:])
			return true

		case "chatleave":
			leaveChat(commandFields[1:])
			return true

		case "chatmemberlist":
			listChatMembers(commandFields[1:])
			return true

		case "chatmemberping":
			pingChatMember(commandFields[1:])
			return true

		case "msg":
			sendMessage(commandFields[1:])
			return true

		case "execute":
			executeScript(commandFields[1:])
			return true

		case "sleep":
			sleepScript(commandFields[1:])
			return true

		case "echo":
			echoScript(commandFields[1:])
			return true

		case "ping":
			ping(commandFields[1:])
			return true

		case "log":
			cmdLogging(commandFields[1:])
			return true

		case "quit":
			quitCmdTool(commandFields[1:])
			return true

		case "play":
			play(commandFields[1:])
			return true

		default:
			usage()
			return false
		}
	}
	return false
}

// Display the usage of all available commands
func usage() {
	for _, key := range commandKeys {
		fmt.Printf("%v\n", commands[key])
	}

}

func quitCmdTool(arguments []string) {

	// Get rid of warnings
	_ = arguments

	os.Exit(0)
}

func sendMessage(arguments []string) {

	for k, v := range chatMembers {

		// Do only send to others, not to yourself
		if k != conf.Name {

			// create TCP connection to recipient
			conn, err := net.Dial("tcp", v.Sender)
			if err != nil {
				fmt.Printf("could not dial to %v: %v\n", v.Sender, err)
				return
			}

			// send message
			fmt.Fprintf(conn, "%s%s\n", prompt(), strings.Join(arguments, " "))

			// close connection
			conn.Close()
		}
	}
}

func scriptPrompt(scriptname string) string {
	return fmt.Sprintf("<%s %q> ", time.Now().Format("Jan 2 15:04:05.000"), scriptname)
}

func executeScript(arguments []string) {

	if len(arguments) == 0 {
		fmt.Printf("error: no filename to execute specified\n")
		return
	}

	b, err := ioutil.ReadFile(arguments[0])
	if err != nil {
		fmt.Printf("ioutil.ReadFile: %v\n", err)
		return
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		//fmt.Printf("EXECUTE %d: %q\n", i, line)
		if strings.TrimSpace(line) == "" ||
			strings.Split(strings.TrimSpace(line), "")[0] == "#" {
			continue
		}
		echoScript(strings.Split(scriptPrompt(arguments[0])+line, " "))
		if _, ok := commands[strings.Split(line, " ")[0]]; ok {
			executeCommand(line)
		} else {
			fmt.Printf("error: %q is an unknown command\n", strings.Split(line, " ")[0])
		}
	}

}

func sleepScript(arguments []string) {

	var numSeconds int

	if len(arguments) == 0 {
		numSeconds = 1
	} else {
		numSeconds, err = strconv.Atoi(arguments[0])
	}

	time.Sleep(time.Second * time.Duration(numSeconds))
}

func echoScript(arguments []string) {

	fmt.Printf("%s\n", strings.Join(arguments, " "))
}

func ping(arguments []string) {

	if len(arguments) > 0 &&
		arguments[0] != chatSelf.Name &&
		chatMembers[arguments[0]] != nil {

		pingMember(chatMembers[arguments[0]])
	}
}

func cmdLogging(arguments []string) {

	if len(arguments) == 0 ||
		(len(arguments) == 1 && arguments[0] != "off") {
		fmt.Printf("Error: wrong input. Usage: \n\t 'log (on <filename>) | off\n")

		return
	}

	if arguments[0] == "on" && len(arguments) > 1 {
		log.Printf("Switch to logging by command to %q\n", arguments[1])
		tmpDebugfile, err = startLogging(arguments[1])
		if err != nil {
			fmt.Printf("Error: startLogging: %v\n", err)
		} else {
			log.Printf("Start logging by command to %q\n", arguments[1])
		}

		return
	}

	if arguments[0] == "off" {
		log.Printf("Stop logging by command")
		_ = tmpDebugfile.Close()

		// Start debugging to file, if switched on or filename specified
		if *debug || len(*debugfilename) > 0 {

			_, err := startLogging(*debugfilename)
			if err != nil {
				fmt.Printf("could not start logging: %v\n", err)
				return
			}
			log.Printf("Switch from logging by command to %q\n", debugfilename)
		}
	}
}

func play(arguments []string) {

	conn, err := net.Dial("tcp", chatMembers[arguments[0]].Sender)
	if err != nil {
		fmt.Printf("could not dial to %v: %v\n", chatMembers[arguments[0]].Sender, err)
		return
	}
	// send message
	fmt.Fprintf(conn, "")
	// close connection
	conn.Close()

	//fmt.Printf("chatStop: %v\n", chatStop)

}
