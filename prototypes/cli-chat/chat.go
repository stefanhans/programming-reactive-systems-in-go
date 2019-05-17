package main

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/golang/protobuf/ptypes"
	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/cli-chat/chat-member"
)

var (
	chatSelf *chatmember.Member = &chatmember.Member{}
	chatStop bool
)

func listenStream(arguments []string) {

	chatStop = false

	go func() {

		// listen for TCP connections on localhost
		listener, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			fmt.Printf("could not listen to localhost:0: %v\n", err)
			return
		}
		defer listener.Close()

		fmt.Printf("listener.Addr(): %v\n", listener.Addr())

		chatSelf.MsgType = chatmember.Member_JOIN
		chatSelf.Name = conf.Name
		chatSelf.Sender = listener.Addr().String()

		// Set timestamp for joined peer
		jointime := ptypes.TimestampNow()
		chatSelf.Timestamp = jointime

		fmt.Printf("chatSelf: %v\n", chatSelf)
		fmt.Printf("chatSelf.MsgType: %v\n", chatSelf.MsgType)

		joiningChat(chatSelf)

		// wait for connections
		for {
			// accept connection
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("could not accept connection: %v\n", err)
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
					fmt.Printf("could not read message: %v\n", err)
					return
				}
				fmt.Printf("\n%s%s", msg, prompt())

				// send reply
				//conn.Write([]byte(fmt.Sprintf("Message accepted: %s", msg)))

				// close connection
				conn.Close()
			}(conn)
		}
	}()
}

func leaveChat(arguments []string) {

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
