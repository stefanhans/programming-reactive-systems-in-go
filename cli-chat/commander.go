package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/memberlist"
	"github.com/pborman/uuid"
	"github.com/stefanhans/programming-reactive-systems-in-go/cli-chat/chat-member"
)

// Execute a command specified by the argument string
func executeCommand(commandline string) bool {

	// Trim prefix and split string by white spaces
	commandFields := strings.Fields(commandline)

	// Check for empty string without prefix
	if len(commandFields) > 0 {

		// Switch according to the first word and call appropriate function with the rest as arguments
		switch commandFields[0] {

		case "memberlistconfigure":
			return configureMemberlist(commandFields[1:])

		case "showconfig":
			showMemberlistConfiguration(commandFields[1:])
			return true

		case "saveconfig":
			return saveMemberlistConfiguration(commandFields[1:])

		case "loadconfig":
			return loadMemberlistConfiguration(commandFields[1:])

		case "memberlistcreate":
			return createMemberlist(commandFields[1:])

		case "showmemberlist":
			showMemberlist(commandFields[1:])
			return true

		case "localnode":
			showLocalNode(commandFields[1:])
			return true

		case "memberlist":
			return listMembers(commandFields[1:])

		case "memberlistjoin":
			return joinMemberlist(commandFields[1:])

		case "memberlistleave":
			return leaveMemberlist(commandFields[1:])

		case "memberlistupdate":
			return updateMemberlist(commandFields[1:])

		case "memberliststart":
			return startBroadcast(commandFields[1:])

		case "memberlistshutdown":
			return shutdownBroadcast(commandFields[1:])

		case "memberlistshutdowntransport":
			return shutdownBroadcastTransport(commandFields[1:])

		case "memberlisthealthscore":
			return getHealthScore(commandFields[1:])

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

		case "bootstrap":
			showBootstrapData(commandFields[1:])
			return true

		case "bootstrapconfig":
			listBootstrapConfig(commandFields[1:])
			return true

		case "bootstrappeers":
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

		case "chatstart":
			return startChat(commandFields[1:])

		case "chatleave":
			return leaveChat(commandFields[1:])

		case "chatmembers":
			listChatMembers(commandFields[1:])
			return true

		case "chatmemberping":
			pingChatMember(commandFields[1:])
			return true

		case "chatstop":
			stopChat(commandFields[1:])
			return true

		case "msg":
			sendMessage(commandFields[1:])
			return true

		case "execute":
			return executeScript(commandFields[1:])

		case "sleep":
			sleepScript(commandFields[1:])
			return true

		case "echo":
			echoScript(commandFields[1:])
			return true

		case "shell":
			return executeShellScript(commandFields[1:])

		case "ping":
			pingChat(commandFields[1:])
			return true

		case "log":
			cmdLogging(commandFields[1:])
			return true

		case "exit":
			exitCmdTool(commandFields[1:])
			return true

		case "help":
			showHelp(commandFields[1:])
			return true

		case "play":
			play(commandFields[1:])
			return true

		case "init":
			return initApplication(commandFields[1:])

		case "quit":
			quitApplication(commandFields[1:])
			return true

		default:

			if *testMode {
				return executeTestEventCommand(commandFields)
			}

			return noCommand(commandFields)
		}
	}
	return false
}

// noCommand says that it's not an available command
func noCommand(arguments []string) bool {
	displayError(fmt.Sprintf("cannot find command %q", arguments[0]))
	return false
}

func exitCmdTool(arguments []string) {

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
				displayText(strings.Trim(fmt.Sprintf("could not dial to %v: %v\n%s", v.Sender, err,
					prompt), "\n"))
				continue
			}

			// send message
			_, err = fmt.Fprintf(conn, "%s%s\n", prompt, strings.Join(arguments, " "))
			if err != nil {
				displayError("could not send message", err)
			}

			// close connection
			_ = conn.Close()
		}
	}
}

func play(arguments []string) {

	// Get rid off warning
	_ = arguments
}

func initApplication(arguments []string) bool {

	// Get rid off warning
	_ = arguments

	/*
		memberlistconfigure
	*/

	// DefaultLocalConfig works like DefaultConfig, however it returns a configuration
	// that is optimized for a local loopback environments. The default configuration is
	// still very conservative and errs on the side of caution.
	conf = memberlist.DefaultLocalConfig()

	if len(arguments) > 1 {
		log.Printf("arguments: %v\n", arguments)
		switch arguments[0] {
		case "BindPort":
			p, err := strconv.Atoi(arguments[1])
			if err != nil {
				displayText(strings.Trim(fmt.Sprintf("could not configure memberlist: %v\n%s", err,
					prompt), "\n"))
				return false
			}
			conf.BindPort = p
		}
	}
	conf.BindAddr = "127.0.0.1"

	// Delegate and Events are delegates for receiving and providing
	// data to memberlist via callback mechanisms. For Delegate, see
	// the Delegate interface. For Events, see the EventDelegate interface.
	//
	// The DelegateProtocolMin/Max are used to guarantee protocol-compatibility
	// for any custom messages that the delegate might do (broadcasts,
	// local/remote state, etc.). If you don't set these, then the protocol
	// versions will just be zero, and version compliance won't be done.
	conf.Delegate = &delegate{}

	// EventDelegate is a simpler delegate that is used only to receive
	// notifications about members joining and leaving. The methods in this
	// delegate may be called by multiple goroutines, but never concurrently.
	// This allows you to reason about ordering.
	conf.Events = &eventDelegate{}

	// NotifyMerge is invoked when a merge could take place.
	// Provides a list of the nodes known by the peer. If
	// the return value is non-nil, the merge is canceled.
	//conf.Merge = &MergeDelegate{}

	//conf.Alive = &AliveDelegate{}

	// The name of this node. This must be unique in the cluster.
	hostname, err := os.Hostname()
	if err != nil {
		displayText(strings.Trim(fmt.Sprintf("could not return hostname from OS: %v\n%s", err,
			prompt), "\n"))
		return false
	}
	id := uuid.NewUUID().String()
	conf.Name = name + "-" + hostname + "-" + id

	mlogfile = "mlist" + name + ".log"

	// Set logger (with output to logfile parameter)
	w := io.Writer(os.Stderr)
	if len(mlogfile) > 0 {
		f, err := os.Create(mlogfile)
		if err == nil {
			w = io.Writer(f)
			shortName = fmt.Sprintf("<%s-%s-*> ", hostname, id[:7])
		}
	}
	lg := log.New(w, shortName, log.LstdFlags|log.Lshortfile)

	conf.Logger = lg
	nc := &memberlist.NetTransportConfig{
		BindAddrs: []string{"127.0.0.1"},
		BindPort:  0,
		Logger:    lg,
	}
	if nt, err := memberlist.NewNetTransport(nc); err == nil {
		conf.Transport = nt
	}

	/*
		memberlistcreate
	*/

	// Create will create a new Memberlist using the given configuration.
	// This will not connect to any other node (see Join) yet, but will start
	// all the listeners to allow other nodes to join this memberlist.
	// After creating a Memberlist, the configuration given should not be
	// modified by the user anymore.
	mlist, err = memberlist.Create(conf)
	if err != nil {
		log.Printf("create memberlist failed: %v\n", err)
	}
	log.Printf("local node name: %v\n", mlist.LocalNode().Name)

	/*
		bootstrapjoin
	*/
	if bootstrapApi == nil {
		initializeBootstrapApi()
		if bootstrapApi == nil {
			return false
		}
	}

	bootstrapData = bootstrapApi.Join()
	//displayBootstrapData(bootstrapData)

	/*
		memberliststart
	*/
	if mlist == nil {
		log.Printf("no memberlist found\n%s", prompt)
		return false
	}

	// TransmitLimitedQueue is used to queue messages to broadcast to
	// the cluster (via gossip) but limits the number of transmits per
	// message. It also prioritizes messages with lower transmit counts
	// (hence newer messages).
	broadcasts = &memberlist.TransmitLimitedQueue{
		// NumNodes returns the number of nodes in the cluster. This is
		// used to determine the retransmit count, which is calculated
		// based on the log of this.
		NumNodes: func() int {
			return mlist.NumMembers()
		},

		// RetransmitMult is the multiplier used to determine the maximum
		// number of retransmissions attempted.
		RetransmitMult: 3,
	}

	// LocalNode is used to return the local Node
	node := mlist.LocalNode()
	log.Printf("Local member %s:%d\n", node.Addr, node.Port)

	/*
		memberlistjoin
	*/

	if bootstrapData == nil {
		displayError("no bootstrap data found")
		return false
	}

	// No bootstrap server
	if len(bootstrapData.Peers) == 0 {
		log.Printf("No bootstrap server\n")
		return false
	}

	// No other bootstrap server
	if len(bootstrapData.Peers) == 1 {
		if _, ok := bootstrapData.Peers[bootstrapApi.Self.ID]; !ok {
			log.Printf("No other bootstrap server\n")
			return false
		}
	}

	// Join others from bootstrap peers to slice of members
	var bootstrapAddresses []string

	for _, v := range bootstrapData.Peers {
		bootstrapAddresses = append(bootstrapAddresses, fmt.Sprintf("%s:%s", v.Ip, v.Port))
	}

	n, err := mlist.Join(bootstrapAddresses)
	if err != nil {
		displayError("failed to join bootstrap peers to members", err)
		return false
	}
	log.Printf("%d host(s) successfully contacted\n", n)

	log.Printf("Known nodes: %v\n", mlist.Members())

	/*
		chatstart
	*/

	chatStop = false

	go func() {

		// listen for TCP connections on localhost
		listener, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			displayError("could not listen to localhost:0", err)
			return
		}
		defer listener.Close()

		log.Printf("listener.Addr(): %v\n", listener.Addr())

		chatSelf.MsgType = chatmember.Member_JOIN
		chatSelf.Name = conf.Name
		chatSelf.Sender = listener.Addr().String()

		// Set timestamp for joined peer
		jointime := ptypes.TimestampNow()
		chatSelf.Timestamp = jointime

		log.Printf("chatSelf: %v\n", chatSelf)
		log.Printf("chatSelf.MsgType: %v\n", chatSelf.MsgType)

		err = joiningChat(chatSelf)
		if err != nil {
			displayError("could not join chat", err)
			return
		}

		// wait for connections
		for {
			// accept connection
			conn, err := listener.Accept()
			if err != nil {
				displayError("could not accept connection", err)
				continue
			}

			if chatStop {
				_ = conn.Close()
				return
			}

			// create a goroutine for connection
			go func(conn net.Conn) {

				// read and print the message
				msg, err := bufio.NewReader(conn).ReadString('\n')

				if err != nil {
					if err == io.EOF {
						_ = conn.Close()
						return
					}
					displayError("could not read message", err)
					return
				}

				displayColoredMessages(msg)

				if strings.Fields(msg)[0] == "<left>" {
					if bootstrapApi == nil {
						initializeBootstrapApi()
						if bootstrapApi != nil {
							bootstrapData = bootstrapApi.Refill()
						}
					} else {
						bootstrapData = bootstrapApi.Refill()
						log.Printf("%v", bootstrapData)
					}
				}

				// close connection
				_ = conn.Close()
			}(conn)
		}
	}()

	// Send leave message to all chat members
	for k, v := range chatMembers {

		// Do only send to others, not to yourself
		if k != conf.Name {

			// create TCP connection to recipient
			conn, err := net.Dial("tcp", v.Sender)
			if err != nil {
				displayText(strings.Trim(fmt.Sprintf("could not dial to %v: %v\n%s", v.Sender, err,
					prompt), "\n"))
				continue
			}

			// Send message tagged as joined
			_, err = fmt.Fprintf(conn, "<joined> %s has joined\n", name)
			if err != nil {
				displayError("could not send joined message", err)
			}

			// close connection
			_ = conn.Close()
		}
	}
	return true
}

func quitApplication(arguments []string) {

	// Get rid off warning
	_ = arguments

	// Send message about leaving the bootstrap peers
	// before a 'left' tagged message triggers bootstrap refill requests
	if bootstrapApi != nil {
		bootstrapApi.Leave(bootstrapApi.Self.ID)
	}

	// Send leave message to all chat members
	for k, v := range chatMembers {

		// Do only send to others, not to yourself
		if k != conf.Name {

			// create TCP connection to recipient
			conn, err := net.Dial("tcp", v.Sender)
			if err != nil {
				displayError(fmt.Sprintf("could not dial to %v", v.Sender), err)
				continue
			}

			// Send message tagged as left
			_, err = fmt.Fprintf(conn, "<left> %s is leaving\n", name)
			if err != nil {
				displayError("could not send left message", err)
			}

			// close connection
			_ = conn.Close()
		}
	}

	// Send message about leaving via memberlist
	if chatSelf.Name != "" {
		err = leavingChat(chatSelf)
		if err != nil {
			displayError("could not leave chat", err)
		}
	}

	// Last entry in the logfile
	log.Printf("Session finished\n")

	// Exit the application
	os.Exit(0)
}
