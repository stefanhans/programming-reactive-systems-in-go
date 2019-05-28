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
	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/liner-memberlist-chat/chat-member"
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
		fmt.Printf("broadcastadd <key> <message> \n\t broadcastadd updates a key/message at all members\n")
		return
	}

	mtx.Lock()
	items[arguments[0]] = strings.Join(arguments[1:], " ")
	mtx.Unlock()

	fmt.Printf("items: %v\n", items)

	b, err := json.Marshal([]*update{
		{
			Action: "add",
			Data: map[string]string{
				arguments[0]: strings.Join(arguments[1:], " "),
			},
		},
	})

	if err != nil {
		fmt.Printf("could not marshall update struct: %v\n", err)
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

		fmt.Printf("%d:\titems[%q] deleted\n", i, k)

		b, err := json.Marshal([]*update{
			{
				Action: "del",
				Data: map[string]string{
					k: "",
				},
			},
		})
		if err != nil {
			fmt.Printf("could not marshall update struct: %v\n", err)
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
	b, _ := json.MarshalIndent(chatMembers, "", "    ")
	fmt.Printf("%s\n", b)
}

func pingChatMember(arguments []string) {

	if len(arguments) > 0 &&
		arguments[0] != chatSelf.Name &&
		chatMembers[arguments[0]] != nil {

		conn, err := net.Dial("tcp", chatMembers[arguments[0]].Sender)
		if err != nil {
			fmt.Printf("could not dial to %v: %v\n", chatMembers[arguments[0]].Sender, err)
			return
		}

		// send message
		fmt.Fprintf(conn, "")

		// close connection
		conn.Close()
	}
}

func isChatMemberReachable(name string) (bool, error) {

	if name != "" &&
		name != chatSelf.Name &&
		chatMembers[name] != nil {

		conn, err := net.Dial("tcp", chatMembers[name].Sender)
		if err != nil {
			return false, nil
		}

		// send message
		fmt.Fprintf(conn, "")

		// close connection
		conn.Close()

		return true, nil
	}
	return false, fmt.Errorf("%q is not a valid chat member", name)
}

func joiningChat(chatMember *chatmember.Member) {

	mtx.Lock()
	chatMembers[chatMember.Name] = chatMember
	mtx.Unlock()

	b, err := json.Marshal(chatMember)
	if err != nil {
		fmt.Printf("could not marshall joining chat member: %v\n", err)
		return
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

}
func leavingChat(chatMember *chatmember.Member) {

	chatMember.MsgType = chatmember.Member_LEAVE

	b, err := json.Marshal(chatMember)
	if err != nil {
		fmt.Printf("could not marshall leaving chat member: %v\n", err)
		return
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

	// QueueBroadcast is used to enqueue a broadcast
	broadcasts.QueueBroadcast(&broadcast{
		msg:    b,
		notify: notifyChan,
	})

}

func pingMember(chatMember *chatmember.Member) {

	chatMember.MsgType = chatmember.Member_PING

	now, _ := ptypes.TimestampProto(time.Now())
	chatMember.Timestamp = now

	b, err := json.Marshal(chatMember)
	if err != nil {
		fmt.Printf("could not marshall chat member to ping: %v\n", err)
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
