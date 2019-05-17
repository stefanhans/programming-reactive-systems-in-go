package main

import (
	"fmt"

	"github.com/hashicorp/memberlist"
)

// EventDelegate is a simpler delegate that is used only to receive
// notifications about members joining and leaving. The methods in this
// delegate may be called by multiple goroutines, but never concurrently.
// This allows you to reason about ordering.
type eventDelegate struct{}

// NotifyJoin is invoked when a node is detected to have joined.
// The Node argument must not be modified.
func (d *eventDelegate) NotifyJoin(node *memberlist.Node) {
	fmt.Printf("NotifyJoin: %v\n", node.Name)
}

// NotifyLeave is invoked when a node is detected to have left.
// The Node argument must not be modified.
func (d *eventDelegate) NotifyLeave(node *memberlist.Node) {
	fmt.Printf("NotifyLeave: %v %v\n", node.Name, node.Address())

	bootstrapData = bootstrapApi.List()

	var leftBootstrapNode string
	for k, v := range bootstrapData.Peers {
		if v.Name == node.Name {
			leftBootstrapNode = k
		}
	}

	if leftBootstrapNode != "" {
		bootstrapApi.Leave(leftBootstrapNode)
		delete(bootstrapData.Peers, leftBootstrapNode)
	}

	bootstrapData = bootstrapApi.Refill()
	fmt.Printf("%v\n", bootstrapData.Config)

	fmt.Printf("Ping chat member %v\n", node.Name)

	reachable, err := isChatMemberReachable(node.Name)
	if err != nil {
		fmt.Printf("could not check if reachable: %v\n", err)
		return
	}
	if !reachable {
		delete(chatMembers, node.Name)
	}
}

// NotifyUpdate is invoked when a node is detected to have
// updated, usually involving the meta data. The Node argument
// must not be modified.
func (d *eventDelegate) NotifyUpdate(node *memberlist.Node) {
	fmt.Printf("NotifyUpdate: %v\n", node.Name)
}

/*

// AliveDelegate is used to involve a client in processing
// a node "alive" message. When a node joins, either through
// a UDP gossip or TCP push/pull, we update the state of
// that node via an alive message. This can be used to filter
// a node out and prevent it from being considered a peer
// using application specific logic.
type AliveDelegate struct{}

func (d *AliveDelegate) NotifyAlive(node *memberlist.Node) error {
	fmt.Printf("NotifyAlive: %v\n", node.Name)

	//fmt.Printf("--------\n")
	//for i, node := range mlist.Members() {
	//	fmt.Printf("%d %s\n", i, node.Name)
	//}
	//fmt.Printf("--------\n")

	return nil
}


// MergeDelegate is used to involve a client in
// a potential cluster merge operation. Namely, when
// a node does a TCP push/pull (as part of a join),
// the delegate is involved and allowed to cancel the join
// based on custom logic. The merge delegate is NOT invoked
// as part of the push-pull anti-entropy.
type MergeDelegate  struct{}

// NotifyMerge is invoked when a merge could take place.
// Provides a list of the nodes known by the peer. If
// the return value is non-nil, the merge is canceled.
func (d *MergeDelegate) NotifyMerge(peers []*memberlist.Node) error {
	fmt.Printf("NotifyMerge: \n")

	fmt.Printf("--------\n")
	for i, node := range peers {
		fmt.Printf("MergeNodes %d %s\n", i, node.Name)
	}
	fmt.Printf("--------\n")


	fmt.Printf("--------\n")
	for i, node := range mlist.Members() {
		fmt.Printf("mlist.Members() %d %s\n", i, node.Name)
	}
	fmt.Printf("--------\n")

	return nil
}

*/
