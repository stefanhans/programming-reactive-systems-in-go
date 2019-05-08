package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/memberlist"
	"github.com/stefanhans/programming-reactive-systems-in-go/cli-chat/chat-member"
)

// Broadcast is something that can be broadcasted via gossip to
// the memberlist cluster.
type broadcast struct {
	msg    []byte
	notify chan<- struct{}
}

// Invalidates checks if enqueuing the current broadcast
// invalidates a previous broadcast
func (b *broadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

// Returns a byte form of the message
func (b *broadcast) Message() []byte {
	return b.msg
}

// Finished is invoked when the message will no longer
// be broadcast, either due to invalidation or to the
// transmit limit being reached
func (b *broadcast) Finished() {
	if b.notify != nil {
		close(b.notify)
	}
}

func broadcastAddMessage(arguments []string) {

	if len(arguments) < 2 {
		displayError("not enough arguments specified")
		return
	}

	mtx.Lock()
	items[arguments[0]] = strings.Join(arguments[1:], " ")
	mtx.Unlock()

	log.Printf("items: %v\n", items)

	b, err := json.Marshal([]*update{
		{
			Action: "add",
			Data: map[string]string{
				arguments[0]: strings.Join(arguments[1:], " "),
			},
		},
	})

	if err != nil {
		displayError("could not marshall update struct", err)
		return
	}

	// Channel listening for enqueued messages
	notifyChan := make(chan struct{})
	go func() {
		for {
			select {
			case <-notifyChan:
				log.Printf("Items after \"add\": %v\n", items)
				return
			}
		}
	}()

	// QueueBroadcast is used to enqueue a broadcast
	broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte("d"), b...),
		notify: notifyChan,
	})

}

func broadcastDelMessage(arguments []string) {

	for i, k := range arguments {

		mtx.Lock()
		delete(items, k)
		mtx.Unlock()

		log.Printf("%d:\titems[%q] deleted\n", i, k)

		b, err := json.Marshal([]*update{
			{
				Action: "del",
				Data: map[string]string{
					k: "",
				},
			},
		})
		if err != nil {
			displayError("could not marshall update struct", err)
			return
		}

		// Channel listening for enqueued messages
		notifyChan := make(chan struct{})
		go func() {
			for {
				select {
				case <-notifyChan:
					log.Printf("Items after \"del\": %v\n", items)
					return
				}
			}
		}()

		// QueueBroadcast is used to enqueue a broadcast
		broadcasts.QueueBroadcast(&broadcast{
			msg:    append([]byte("d"), b...),
			notify: notifyChan,
		})
	}
}

func listChatMembers(arguments []string) {

	// Get rid off warning
	_ = arguments

	b, _ := json.MarshalIndent(chatMembers, "", "    ")
	displayText(strings.Trim(fmt.Sprintf("%s\n%s", b,
		prompt), "\n"))
}

func pingChatMember(arguments []string) {

	if len(arguments) > 0 &&
		arguments[0] != chatSelf.Name &&
		chatMembers[arguments[0]] != nil {

		conn, err := net.Dial("tcp", chatMembers[arguments[0]].Sender)
		if err != nil {
			displayError(fmt.Sprintf("could not dial to %v", chatMembers[arguments[0]].Sender), err)
			return
		}

		// send message
		_, err = fmt.Fprintf(conn, "")
		if err != nil {
			displayError("could not send ping", err)
		}

		// close connection
		_ = conn.Close()
	}
}

// isChatMemberReachable send an empty request on the application layer
// after a NotifyLeave event from the memberlist layer
func isChatMemberReachable(name string) (bool, error) {

	// Does the request makes sense?
	if name != "" &&
		name != chatSelf.Name &&
		chatMembers[name] != nil {

		conn, err := net.Dial("tcp", chatMembers[name].Sender)
		if err != nil {
			return false, nil
		}

		// send message
		_, err = fmt.Fprintf(conn, "")
		if err != nil {
			displayError("could not send ping", err)
		}

		// close connection
		_ = conn.Close()

		return true, nil
	}
	return false, fmt.Errorf("%q is not a valid chat member to ping", name)
}

func joiningChat(chatMember *chatmember.Member) error {
	log.Printf("joiningChat(%v)\n", chatMember.Name)

	mtx.Lock()
	chatMembers[chatMember.Name] = chatMember
	mtx.Unlock()

	b, err := json.Marshal(chatMember)
	if err != nil {
		displayError("could not marshall joining chat member", err)
		return err
	}

	// Channel listening for enqueued messages
	notifyChan := make(chan struct{})
	go func() {
		for {
			select {
			case <-notifyChan:
				log.Printf("chatMembers after added %q: %v\n", chatMember.Name, chatMembers)
				return
			}
		}
	}()

	// QueueBroadcast is used to enqueue a broadcast
	broadcasts.QueueBroadcast(&broadcast{
		msg:    b,
		notify: notifyChan,
	})

	return nil
}

func leavingChat(chatMember *chatmember.Member) error {
	log.Printf("leavingChat(%v)\n", chatMember.Name)

	chatMember.MsgType = chatmember.Member_LEAVE

	b, err := json.Marshal(chatMember)
	if err != nil {
		displayError("could not marshall leaving chat member: %v\n%s", err)
		return err
	}

	// Channel listening for enqueued messages
	notifyChan := make(chan struct{})
	go func() {
		for {
			select {
			case <-notifyChan:
				log.Printf("chatMembers after deleted %q: %v\n", chatMember.Name, chatMembers)
				return
			}
		}
	}()

	if broadcasts == nil {
		logRed("no broadcasts defined")
		displayError("no broadcasts defined")
		displayText(prompt)
		return nil
	}

	// QueueBroadcast is used to enqueue a broadcast
	broadcasts.QueueBroadcast(&broadcast{
		msg:    b,
		notify: notifyChan,
	})

	return nil
}

func pingMember(chatMember *chatmember.Member) {

	chatMember.MsgType = chatmember.Member_PING

	now, _ := ptypes.TimestampProto(time.Now())
	chatMember.Timestamp = now

	b, err := json.Marshal(chatMember)
	if err != nil {
		displayError("could not marshall chat member", err)
		return
	}

	// Channel listening for enqueued messages
	notifyChan := make(chan struct{})
	go func() {
		for {
			select {
			case <-notifyChan:
				log.Printf("chatMembers after ping %q: %v\n", chatMember.Name, chatMembers)
				return
			}
		}
	}()

	// QueueBroadcast is used to enqueue a broadcast
	broadcasts.QueueBroadcast(&broadcast{
		msg:    b,
		notify: notifyChan,
	})
}
