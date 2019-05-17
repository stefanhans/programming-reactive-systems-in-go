package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/pborman/uuid"
)

var (
	conf  *memberlist.Config
	mlist *memberlist.Memberlist

	shortName string

	mlogfile string
)

func configureMemberlist(arguments []string) {

	// DefaultLocalConfig works like DefaultConfig, however it returns a configuration
	// that is optimized for a local loopback environments. The default configuration is
	// still very conservative and errs on the side of caution.
	conf = memberlist.DefaultLocalConfig()

	if len(arguments) > 1 {
		fmt.Printf("arguments: %v\n", arguments)
		switch arguments[0] {
		case "BindPort":
			p, err := strconv.Atoi(arguments[1])
			if err != nil {
				fmt.Printf("could not configure memberlist: %v\n", err)
				return
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
		fmt.Printf("could not return hostname from OS: %v\n", err)
		return
	}
	id := uuid.NewUUID().String()
	conf.Name = hostname + "-" + name + "-" + id

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
}

func showMemberlistConfiguration(arguments []string) {

	// Get rid off warning
	_ = arguments

	//fmt.Printf("Memberlist Configuration: \n")
	//fmt.Printf("--------------------------\n")
	//fmt.Printf("Name: %q\n", c.Name)
	//fmt.Printf("Transport: %v (interface)\n", c.Transport)
	//fmt.Printf("BindAddr: %q\n", c.BindAddr)
	//fmt.Printf("BindPort: %d\n", c.BindPort)

	for i, w := range strings.Split(
		strings.Replace(
			strings.Replace(fmt.Sprintf("%#v\n", conf), "}", "", 1), "{", ": ", 1), " ") {

		fmt.Printf("%d:\t%v\n", i, w)
	}
}

func saveMemberlistConfiguration(arguments []string) {

	// Marshal array of struct
	if byteArray, err := json.Marshal(conf); err != nil {
		fmt.Printf("json.Marshal: %v\n", err)
		return
	} else {
		var filename string

		if len(arguments) > 0 {
			filename = arguments[0]
		} else {
			filename = conf.Name + ".json"
		}

		//Formated JSON
		var out bytes.Buffer
		err = json.Indent(&out, byteArray, "", "\t")

		// Replace '{}' with 'null' in JSON string
		str := strings.Replace(string(out.Bytes()), "{}", "null", -1)

		ioutil.WriteFile(filename, append([]byte(str), byte('\n')), 0600)

		fmt.Printf("%s\n", filename)
	}
}

func loadMemberlistConfiguration(arguments []string) {

	if len(arguments) == 0 {
		fmt.Printf("error: no filename to load specified\n")
		return
	}

	b, err := ioutil.ReadFile(arguments[0])
	if err != nil {
		fmt.Printf("ioutil.ReadFile: %v\n", err)
		return
	}

	var c memberlist.Config
	err = json.Unmarshal(b, &c)
	if err != nil {
		fmt.Printf("json.Unmarshal: %v\n", err)
		return
	}
	conf = &c

	conf.Events = &eventDelegate{}

	mlogfile = "mlist" + name + ".log"

	// Set logger (with output to logfile parameter)
	w := io.Writer(os.Stderr)
	if len(mlogfile) > 0 {
		f, err := os.Create(mlogfile)
		if err == nil {
			w = io.Writer(f)
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
}

func createMemberlist(arguments []string) {

	// Get rid off warning
	_ = arguments

	// Create will create a new Memberlist using the given configuration.
	// This will not connect to any other node (see Join) yet, but will start
	// all the listeners to allow other nodes to join this memberlist.
	// After creating a Memberlist, the configuration given should not be
	// modified by the user anymore.
	mlist, err = memberlist.Create(conf)
	if err != nil {
		fmt.Printf("create memberlist failed: %v\n", err)
	}
	fmt.Printf("local node name: %v\n", mlist.LocalNode().Name)
}

func showMemberlist(arguments []string) {

	// Get rid off warning
	_ = arguments

	for i, w := range strings.Split(
		strings.Replace(
			strings.Replace(fmt.Sprintf("%#v\n", mlist), "}", "", 1), "{", ": ", 1), " ") {

		fmt.Printf("%d:\t%v\n", i, w)
	}
}

func showLocalNode(arguments []string) {

	// Get rid off warning
	_ = arguments

	// LocalNode is used to return the local Node
	node := mlist.LocalNode()
	fmt.Printf("localNode.Name: %v\n", node.Name)
	fmt.Printf("localNode.Address: %v\n", node.Address())
}

func listMembers(arguments []string) {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		fmt.Printf("no memberlist found\n")
		return
	}

	if len(mlist.Members()) == 0 {
		fmt.Printf("empty memberlist\n")
		return
	}
	for i, m := range mlist.Members() {
		fmt.Printf("%d:"+
			"\tName: %v\n"+
			"\tAddress: %v\n", i, m.Name, m.Address())
	}
}

func joinMemberlist(arguments []string) {

	if mlist == nil {
		fmt.Printf("no memberlist found\n")
		return
	}

	if len(arguments) > 0 {

		// Join arguments to slice of members
		n, err := mlist.Join(arguments)
		if err != nil {
			fmt.Printf("join the memberlist failed: %v\n", err)
		}
		fmt.Printf("%d host(s) successfully contacted\n", n)
	} else {

		// No bootstrap server
		if len(bootstrapData.Peers) == 0 {
			fmt.Printf("Nothing to join\n")
			return
		}

		// No other bootstrap server
		if len(bootstrapData.Peers) == 1 {
			if _, ok := bootstrapData.Peers[bootstrapApi.Self.ID]; !ok {
				fmt.Printf("Nothing to join\n")
				return
			}
		}

		// Join others from bootstrap peers to slice of members
		var bootstrapAddresses []string

		for _, v := range bootstrapData.Peers {
			bootstrapAddresses = append(bootstrapAddresses, fmt.Sprintf("%s:%s", v.Ip, v.Port))
		}

		n, err := mlist.Join(bootstrapAddresses)
		if err != nil {
			fmt.Printf("join memberlist failed: %v\n", err)
		}
		fmt.Printf("%d host(s) successfully contacted\n", n)
	}
	fmt.Printf("Known nodes: %v\n", mlist.Members())
}

func leaveMemberlist(arguments []string) {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		fmt.Printf("no memberlist found\n")
		return
	}

	// Leave will broadcast a leave message but will not shutdown the background
	// listeners, meaning the node will continue participating in gossip and state
	// updates.
	//
	// This will block until the leave message is successfully broadcasted to
	// a member of the cluster, if any exist or until a specified timeout
	// is reached.
	//
	// This method is safe to call multiple times, but must not be called
	// after the cluster is already shut down.
	err := mlist.Leave(time.Second)
	if err != nil {
		fmt.Printf("leaveMemberlist failed: %v\n", err)
	}
}

func updateMemberlist(arguments []string) {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		fmt.Printf("no memberlist found\n")
		return
	}

	// UpdateNode is used to trigger re-advertising the local node. This is
	// primarily used with a Delegate to support dynamic updates to the local
	// meta data.  This will block until the update message is successfully
	// broadcasted to a member of the cluster, if any exist or until a specified
	// timeout is reached.
	err := mlist.UpdateNode(time.Second)
	if err != nil {
		fmt.Printf("updateMemberlist failed: %v\n", err)
	}
}

func startBroadcast(arguments []string) {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		fmt.Printf("no memberlist found\n")
		return
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
	fmt.Printf("Local member %s:%d\n", node.Addr, node.Port)
}

func shutdownBroadcast(arguments []string) {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		fmt.Printf("no memberlist found\n")
		return
	}

	// Shutdown will stop any background maintanence of network activity
	// for this memberlist, causing it to appear "dead". A leave message
	// will not be broadcasted prior, so the cluster being left will have
	// to detect this node's shutdown using probing. If you wish to more
	// gracefully exit the cluster, call Leave prior to shutting down.
	//
	// This method is safe to call multiple times.
	err = mlist.Shutdown()
	if err != nil {
		fmt.Printf("mlist.Shutdown(): %v\n", err)
	}
}

func shutdownBroadcastTransport(arguments []string) {

	// Get rid off warning
	_ = arguments

	if conf == nil {
		fmt.Printf("nothing for memberlist configured\n")
		return
	}

	if conf.Transport == nil {
		fmt.Printf("no transport for memberlist configured\n")
		return
	}

	err := conf.Transport.Shutdown()
	if err != nil {
		fmt.Printf("Transport.Shutdown(): %v\n", err)
	}

}

func getHealthScore(arguments []string) {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		fmt.Printf("no memberlist found\n")
		return
	}

	fmt.Printf("Health Score: %d\n", mlist.GetHealthScore())
}

func deleteMemberlist(arguments []string) {

	// Get rid off warning
	_ = arguments

	mlist = nil
}
