package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/memberlist"
	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/cli-chat/chat-member"
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

	if err := json.Unmarshal(b, &newChatMember); err != nil {
		return
	}

	switch newChatMember.MsgType {
	case chatmember.Member_JOIN:

		if _, ok := chatMembers[newChatMember.Name]; ok {
			return
		}

		joiningChat(chatSelf)

		//for k, chatMember := range chatMembers {
		//
		//}

		mtx.Lock()
		chatMembers[newChatMember.Name] = newChatMember
		mtx.Unlock()

	case chatmember.Member_LEAVE:
		mtx.Lock()
		delete(chatMembers, newChatMember.Name)
		mtx.Unlock()

		// Todo: What about Member_PING?
	case chatmember.Member_PING:
		if newChatMember.Name == chatSelf.Name {
			pingMember(newChatMember)
			fmt.Printf("Sent ping back to %s\n", newChatMember.Name)
			return
		}

		now, _ := ptypes.TimestampProto(time.Now())
		secondsNeeded := now.GetSeconds() - newChatMember.Timestamp.GetSeconds()
		fmt.Printf("Ping returned after %d seconds: %v\n", secondsNeeded, newChatMember)
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
