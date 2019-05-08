package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/memberlist"
	"github.com/stefanhans/programming-reactive-systems-in-go/cli-chat/chat-member"
)

var (
	mtx   sync.RWMutex
	items = map[string]string{}

	chatMembers = map[string]*chatmember.Member{}

	broadcasts *memberlist.TransmitLimitedQueue
)

type update struct {
	Action string // add, del
	Data   map[string]string
}

// Delegate is the interface that clients must implement if they want to hook
// into the gossip layer of Memberlist. All the methods must be thread-safe,
// as they can and generally will be called concurrently.
type delegate struct{}

// NodeMeta is used to retrieve meta-data about the current node
// when broadcasting an alive message. It's length is limited to
// the given byte size. This metadata is available in the Node structure.
func (d *delegate) NodeMeta(limit int) []byte {
	return []byte{}
}

// NotifyMsg is called when a user-data message is received.
// Care should be taken that this method does not block, since doing
// so would block the entire UDP packet receive loop. Additionally, the byte
// slice may be modified after the call returns, so it should be copied if needed
func (d *delegate) NotifyMsg(b []byte) {

	if len(b) == 0 {
		return
	}

	var newChatMember *chatmember.Member

	err := json.Unmarshal(b, &newChatMember)
	if err != nil {
		logRed(fmt.Sprintf("could not unmarshall message %q: %v", string(b), err))
		return
	}

	logYellow(fmt.Sprintf("NotifyMsg: %v\n%q\n", newChatMember.MsgType, string(b)))

	switch newChatMember.MsgType {
	case chatmember.Member_JOIN:

		if _, ok := chatMembers[newChatMember.Name]; ok {
			return
		}

		err = joiningChat(chatSelf)
		if err != nil {
			logRed(fmt.Sprintf("could not join chat: %v", err))
			return
		}

		mtx.Lock()
		chatMembers[newChatMember.Name] = newChatMember
		mtx.Unlock()

	case chatmember.Member_LEAVE:
		mtx.Lock()
		delete(chatMembers, newChatMember.Name)
		mtx.Unlock()

		displayYelloText(fmt.Sprintf("Member_LEAVE: %s", newChatMember.Name))

		log.Printf("Ping chat member %q to decide removal from chat\n", newChatMember.Name)

		// Remove the peer from the chat if not reachable
		reachable, err := isChatMemberReachable(newChatMember.Name)
		if err != nil {
			log.Printf("could not check if reachable: %v\n", err)
			return
		}
		if !reachable {
			delete(chatMembers, newChatMember.Name)
			log.Printf("%q not reachable and deleted\n", newChatMember.Name)
		}

		// Refresh bootstrap data
		bootstrapData = bootstrapApi.List()

		for k, v := range bootstrapData.Peers {

			// The node which left is a bootstrap peer?
			if v.Name == newChatMember.Name {

				// Remove it from bootstrap peers
				bootstrapApi.Leave(k)

				// Try to refill bootstrap peers with this peer
				bootstrapData = bootstrapApi.Refill()
			}
		}

		// Todo: What about Member_PING?
	case chatmember.Member_PING:
		if newChatMember.Name == chatSelf.Name {
			pingMember(newChatMember)
			displayText(strings.Trim(fmt.Sprintf("Sent ping back to %s\n%s", newChatMember.Name,
				prompt), "\n"))
			return
		}

		now, _ := ptypes.TimestampProto(time.Now())
		secondsNeeded := now.GetSeconds() - newChatMember.Timestamp.GetSeconds()
		displayText(strings.Trim(fmt.Sprintf("Ping returned after %d seconds: %v\n%s", secondsNeeded, newChatMember,
			prompt), "\n"))
	}
}

// GetBroadcasts is called when user data messages can be broadcast.
// It can return a list of buffers to send. Each buffer should assume an
// overhead as provided with a limit on the total byte size allowed.
// The total byte size of the resulting data to send must not exceed
// the limit. Care should be taken that this method does not block,
// since doing so would block the entire UDP packet receive loop.
func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return broadcasts.GetBroadcasts(overhead, limit)
}

// LocalState is used for a TCP Push/Pull. This is sent to
// the remote side in addition to the membership information. Any
// data can be sent here. See MergeRemoteState as well. The `join`
// boolean indicates this is for a join instead of a push/pull.
func (d *delegate) LocalState(join bool) []byte {
	mtx.RLock()
	m := chatMembers
	mtx.RUnlock()
	b, _ := json.Marshal(m)
	return b
}

// MergeRemoteState is invoked after a TCP Push/Pull. This is the
// state received from the remote side and is the result of the
// remote side's LocalState call. The 'join'
// boolean indicates this is for a join instead of a push/pull.
func (d *delegate) MergeRemoteState(buf []byte, join bool) {
	if len(buf) == 0 {
		return
	}
	if !join {
		return
	}
	var m map[string]*chatmember.Member
	if err := json.Unmarshal(buf, &m); err != nil {
		return
	}
	mtx.Lock()
	for k, v := range m {
		chatMembers[k] = v
	}
	mtx.Unlock()
}
