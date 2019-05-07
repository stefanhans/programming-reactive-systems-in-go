package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

type help struct {
	name        string
	group       string
	usage       string
	description string
	commands    map[string]*help
}

var (
	commandKeys      []string
	groupKeys        []string
	groupCommandKeys = make(map[string][]string)

	mainHelp *help = &help{
		name:  "help",
		usage: "help [ groups | <group> | <command> ]",
		description: "help shows: \n" +
			"- this help text (/help)\n" +
			"- all groups of commands (/help groups)\n" +
			"- all commands of a group (/help groups <group>)\n" +
			"- help text of a command (/help <command>)",
		commands: map[string]*help{
			"bootstrapjoin": {
				name:        "bootstrapjoin",
				group:       "bootstrap",
				usage:       "bootstrapjoin",
				description: "bootstrapjoin joins calling peer to bootstrap peers",
			},
			"bootstrapleave": {
				name:        "bootstrapleave",
				group:       "bootstrap",
				usage:       "bootstrapleave",
				description: "bootstrapleave leaves calling peer from bootstrap peers",
			},
			"bootstraprefill": {
				name:        "bootstraprefill",
				group:       "bootstrap",
				usage:       "bootstraprefill",
				description: "bootstraprefill refills bootstrap peers with calling peer",
			},
			"bootstrap": {
				name:        "bootstrap",
				group:       "bootstrap",
				usage:       "bootstrap",
				description: "bootstrap shows bootstrap data from remote",
			},
			"bootstrapconfig": {
				name:        "bootstrapconfig",
				group:       "bootstrap",
				usage:       "bootstrapconfig",
				description: "bootstrapconfig shows bootstrap configuration from remote",
			},
			"bootstrappeers": {
				name:        "bootstrappeers",
				group:       "bootstrap",
				usage:       "bootstrappeers",
				description: "bootstrappeers lists bootstrap peers from remote",
			},
			"bootstrapreset": {
				name:        "bootstrapreset",
				group:       "bootstrap",
				usage:       "bootstrapreset",
				description: "bootstrapreset resets the bootstrap peers at remote",
			},
			"bootstraplistlocal": {
				name:        "bootstraplistlocal",
				group:       "bootstrap",
				usage:       "bootstraplistlocal",
				description: "bootstraplistlocal lists bootstrap peers from local map",
			},
			"memberlistconfigure": {
				name:        "memberlistconfigure",
				group:       "memberlist",
				usage:       "memberlistconfigure",
				description: "memberlistconfigure creates a default memberlist configuration",
			},
			"showconfig": {
				name:        "showconfig",
				group:       "memberlist",
				usage:       "showconfig",
				description: "showconfig shows the memberlist configuration",
			},
			"saveconfig": {
				name:        "saveconfig",
				group:       "memberlist",
				usage:       "saveconfig [file]",
				description: "saveconfig saves the memberlist configuration as JSON file",
			},
			"loadconfig": {
				name:        "loadconfig",
				group:       "memberlist",
				usage:       "loadconfig file",
				description: "loadconfig loads the memberlist configuration from JSON file",
			},
			"memberlistcreate": {
				name:        "memberlistcreate",
				group:       "memberlist",
				usage:       "memberlistcreate",
				description: "memberlistcreate creates the memberlist specified by the configuration",
			},
			"showmemberlist": {
				name:        "showmemberlist",
				group:       "memberlist",
				usage:       "showmemberlist",
				description: "showmemberlist shows the memberlist",
			},
			"localnode": {
				name:        "localnode",
				group:       "memberlist",
				usage:       "localnode",
				description: "localnode shows the local node's name and address",
			},
			"memberlist": {
				name:        "memberlist",
				group:       "memberlist",
				usage:       "memberlist",
				description: "memberlist lists all members",
			},
			"memberlistjoin": {
				name:        "memberlistjoin",
				group:       "memberlist",
				usage:       "memberlistjoin",
				description: "memberlistjoin joins to memberlist",
			},
			"memberlistleave": {
				name:        "memberlistleave",
				group:       "memberlist",
				usage:       "memberlistleave [<timeout in seconds, default: 1 sec>]",
				description: "memberlistleave broadcasts leave message until finished or timeout is reached",
			},
			"memberlistupdate": {
				name:        "memberlistupdate",
				group:       "memberlist",
				usage:       "[<timeout in seconds, default: 1 sec>]",
				description: "memberlistupdate broadcasts re-advertising the local node message until finished or timeout is reached",
			},
			"memberliststart": {
				name:        "memberliststart",
				group:       "memberlist",
				usage:       "memberliststart",
				description: "memberliststart starts broadcasting to the members",
			},
			"memberlistshutdown": {
				name:        "memberlistshutdown",
				group:       "memberlist",
				usage:       "memberlistshutdown",
				description: "memberlistshutdown stops broadcasting to the members",
			},
			"memberlistshutdowntransport": {
				name:        "memberlistshutdowntransport",
				group:       "memberlist",
				usage:       "memberlistshutdowntransport",
				description: "memberlistshutdowntransport \n\t memberlistshutdowntransport stops broadcasting transport to the members",
			},
			"memberlisthealthscore": {
				name:        "memberlisthealthscore",
				group:       "memberlist",
				usage:       "memberlisthealthscore",
				description: "memberlisthealthscore shows the health score >= 0, lower numbers are better",
			},
			"memberlistdelete": {
				name:        "memberlistdelete",
				group:       "memberlist",
				usage:       "memberlistdelete",
				description: "memberlistdelete sets memberlist = nil",
			},
			"chatstart": {
				name:        "chatstart",
				group:       "chat",
				usage:       "chatstart",
				description: "chatstart starts chat listener and broadcasts the new chat member",
			},
			"chatleave": {
				name:        "chatleave",
				group:       "chat",
				usage:       "chatleave",
				description: "chatleave broadcasts deletion of this chat member",
			},
			"chatping": {
				name:        "chatping",
				group:       "chat",
				usage:       "chatping",
				description: "chatping pings a member of the chat via memberlist",
			},
			"chatmembers": {
				name:        "chatmembers",
				group:       "chat",
				usage:       "chatmembers <chatmember>",
				description: "chatmembers lists all chat members",
			},
			"chatmemberping": {
				name:        "chatmemberping",
				group:       "chat",
				usage:       "chatmemberping <member id>",
				description: "chatmemberping pings a chat member",
			},
			"chatstop": {
				name:        "chatstop",
				group:       "chat",
				usage:       "chatstop",
				description: "chatstop stops the chat listener",
			},
			"msg": {
				name:        "msg",
				group:       "chat",
				usage:       "msg <string>",
				description: "msg sends the rest of the line as message to all other members",
			},
			"broadcastadd": {
				name:        "broadcastadd",
				group:       "broadcast",
				usage:       "broadcastadd <key> <message>",
				description: "broadcastadd updates a key/message at all members",
			},
			"broadcastdel": {
				name:        "broadcastdel",
				group:       "broadcast",
				usage:       "broadcastdel <key>",
				description: "broadcastdel deletes a key at all members",
			},
			"broadcastlist": {
				name:        "broadcastlist",
				group:       "broadcast",
				usage:       "broadcastlist",
				description: "broadcastlist lists all local key/value pairs",
			},
			"execute": {
				name:        "execute",
				group:       "script",
				usage:       "execute file",
				description: "execute executes the commands in the file line by line; '#' is a comment line",
			},
			"sleep": {
				name:        "sleep",
				group:       "script",
				usage:       "sleep <seconds>",
				description: "sleep sleeps for seconds",
			},
			"echo": {
				name:        "echo",
				group:       "script",
				usage:       "echo <string>",
				description: "echo prints rest of line",
			},
			"lock": {
				name:        "lock",
				group:       "script",
				usage:       "lock <name>...",
				description: "lock creates a <name>.wait lockfile",
			},
			"unlock": {
				name:        "unlock",
				group:       "script",
				usage:       "unlock <name>...",
				description: "unlock removes a <name>.wait lockfile",
			},
			"wait": {
				name:        "wait",
				group:       "script",
				usage:       "wait",
				description: "wait until a <name>.wait lockfile is gone",
			},
			"shell": {
				name:        "shell",
				group:       "script",
				usage:       "shell <script>",
				description: "shell executes the shell script",
			},
			"log": {
				name:        "log",
				group:       "cli",
				usage:       "log (on <filename> | off)",
				description: "log starts or stops writing logging output in the specified file",
			},
			"exit": {
				name:        "exit",
				group:       "cli",
				usage:       "exit",
				description: "exit exits the application directly",
			},
			"play": {
				name:        "play",
				group:       "development",
				usage:       "play [...]",
				description: "play is for developer to play",
			},
			"init": {
				name:        "init",
				group:       "development",
				usage:       "init",
				description: "init simulates the init process of the application",
			},
			"quit": {
				name:        "quit",
				group:       "development",
				usage:       "quit",
				description: "quit simulates the exit process of the application",
			},
		},
	}
)

func helpInit() {

	// To store the commands in sorted order
	for _, command := range mainHelp.commands {
		commandKeys = append(groupKeys, command.name)
	}
	sort.Strings(commandKeys)

	// To store the groups in sorted order
	var exists bool
	for _, command := range mainHelp.commands {

		exists = false
		for _, group := range groupKeys {
			if group == command.group {
				exists = true
			}
		}
		if !exists {
			groupKeys = append(groupKeys, command.group)
		}
	}
	sort.Strings(groupKeys)

	for _, helpGroup := range groupKeys {
		var helpGroupCommands []string
		for _, helpGroupCommand := range mainHelp.commands {
			if helpGroupCommand.group == helpGroup {
				helpGroupCommands = append(helpGroupCommands, helpGroupCommand.name)
			}
		}
		sort.Strings(helpGroupCommands)

		groupCommandKeys[helpGroup] = helpGroupCommands
	}
}

//
func showHelp(arguments []string) {

	switch len(arguments) {
	case 0:
		displayHelp(mainHelp)
	case 1:
		displayHelp(mainHelp, arguments[0])
	case 2:
		displayHelp(mainHelp, arguments[0], arguments[1])
	}
}

func displayHelp(help *help, cmd ...string) error {

	log.Printf("displayHelp: %v\n", cmd)

	// Shows main help
	if len(cmd) == 0 {
		return displayText(strings.Trim(fmt.Sprintf("Usage: %s\nDescription: %s\n%s", help.usage, help.description,
			prompt), "\n"))
	}

	// commands related?
	if cmd[0] == "commands" {
		log.Printf("in commands\n")

		// Shows available commands
		var helpCommands string
		for _, k := range commandKeys {
			helpCommands += fmt.Sprintf("%v ", k)
		}
		return displayText(strings.Trim(fmt.Sprintf("%s\n%s", helpCommands, prompt), "\n"))
	}

	// groups related?
	if cmd[0] == "groups" {
		log.Printf("in groups\n")

		// Shows available groups
		if len(cmd) == 1 {
			var helpGroups string
			for _, k := range groupKeys {
				helpGroups += fmt.Sprintf("%v ", k)
			}
			return displayText(strings.Trim(fmt.Sprintf("%s\n%s", helpGroups, prompt), "\n"))
		}

		// Shows commands for a group
		helpGroup := cmd[1]
		var groupCommands string
		for _, command := range groupCommandKeys[helpGroup] {

			groupCommands += fmt.Sprintf("%v ", command)
		}
		return displayText(strings.Trim(fmt.Sprintf("%s\n%s", groupCommands, prompt), "\n"))
	}

	// Shows help for command
	if commandHelp, ok := help.commands[cmd[0]]; ok {
		return displayText(strings.Trim(fmt.Sprintf("Usage: %s\nDescription: %s\n%s", commandHelp.usage, commandHelp.description,
			prompt), "\n"))
	}
	return displayError(fmt.Sprintf("%q is not a known command", cmd[0]))

}
