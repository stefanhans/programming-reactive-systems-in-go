package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/stefanhans/programming-reactive-systems-in-go/cli-chat/chat-member"
)

var (
	chatSelf *chatmember.Member = &chatmember.Member{}
	chatStop bool
)

// startChat is the CLI function of 'chatjoin' and initiates the chat listener
func startChat(arguments []string) {

	// Get rid off warning
	_ = arguments

	if conf == nil {
		displayError("/chatjoin: could not start listener without memberlist configuration")
		return
	}

	if mlist == nil {
		displayError("/chatjoin: could not start listener without created memberlist")
		return
	}

	if broadcasts == nil {
		displayError("/chatjoin: could not start listener without broadcasting memberlist")
		return
	}

	chatStop = false

	go func() {

		// listen for TCP connections on localhost
		listener, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			displayError("could not listen to localhost:0", err)
			return
		}
		defer listener.Close()

		displayText(strings.Trim(fmt.Sprintf("listener.Addr(): %v\n%s", listener.Addr(),
			prompt), "\n"))

		chatSelf.MsgType = chatmember.Member_JOIN
		chatSelf.Name = conf.Name
		chatSelf.Sender = listener.Addr().String()

		// Set timestamp for joined peer
		jointime := ptypes.TimestampNow()
		chatSelf.Timestamp = jointime

		//fmt.Printf("chatSelf: %v\n", chatSelf)
		//fmt.Printf("chatSelf.MsgType: %v\n", chatSelf.MsgType)

		err = joiningChat(chatSelf)
		if err != nil {
			displayError("could not join the chat", err)
			return
		}

		// wait for connections
		for {
			// accept connection
			conn, err := listener.Accept()
			if err != nil {
				displayText(strings.Trim(fmt.Sprintf("could not accept connection: %v\n%s", err,
					prompt), "\n"))
				continue
			}

			if chatStop {
				conn.Close()
				return
			}

			// create a goroutine for connection
			go func(conn net.Conn) {

				// read and print the message
				msg, err := bufio.NewReader(conn).ReadString('\n')

				if err != nil {
					if err == io.EOF {
						conn.Close()
						return
					}
					displayError("could not read message", err)
					return
				}

				if strings.Fields(msg)[0] == "<left>" {
					if bootstrapApi == nil {
						initializeBootstrapApi()
						if bootstrapApi != nil {
							bootstrapData = bootstrapApi.Refill()
						}
					}
				}

				displayColoredMessages(msg)

				// close connection
				conn.Close()
			}(conn)
		}
	}()
}

func leaveChat(arguments []string) {

	// Get rid off warning
	_ = arguments

	// Stop chat listener
	chatStop = true

	// Empty list of chat member
	mtx.Lock()
	for k := range chatMembers {
		delete(chatMembers, k)
	}
	mtx.Unlock()

	// Send message about leaving via memberlist
	leavingChat(chatSelf)
}

func pingChat(arguments []string) {

	if len(arguments) > 0 &&
		arguments[0] != chatSelf.Name &&
		chatMembers[arguments[0]] != nil {

		pingMember(chatMembers[arguments[0]])
	}
}

func stopChat(arguments []string) {

	// Get rid off warning
	_ = arguments

	// Stop chat listener
	chatStop = true

	conn, _ := net.Dial("tcp", chatSelf.Sender)
	conn.Close()

	displayText(strings.Trim(fmt.Sprintf("stopped chat listener of %q\n%s", chatSelf.Sender,
		prompt), "\n"))
}
