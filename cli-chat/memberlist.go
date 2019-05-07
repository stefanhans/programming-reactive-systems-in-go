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
		log.Printf("arguments: %v\n", arguments)
		switch arguments[0] {
		case "BindPort":
			p, err := strconv.Atoi(arguments[1])
			if err != nil {
				displayText(strings.Trim(fmt.Sprintf("could not configure memberlist: %v\n%s", err,
					prompt), "\n"))
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
		displayText(strings.Trim(fmt.Sprintf("could not return hostname from OS: %v\n%s", err,
			prompt), "\n"))
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
	displayText(prompt)
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

		displayText(strings.Trim(fmt.Sprintf("%d:\t%v\n%s", i, w,
			prompt), "\n"))
	}

	displayText(prompt)
}

func saveMemberlistConfiguration(arguments []string) {

	// Marshal array of struct
	byteArray, err := json.Marshal(conf)
	if err != nil {
		displayText(strings.Trim(fmt.Sprintf("json.Marshal: %v\n%s", err,
			prompt), "\n"))
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

		err = ioutil.WriteFile(filename, append([]byte(str), byte('\n')), 0600)

		log.Printf("%s\n", filename)
	}
}

func loadMemberlistConfiguration(arguments []string) {

	if len(arguments) == 0 {
		displayText(strings.Trim(fmt.Sprintf("error: no filename to load specified\n%s",
			prompt), "\n"))
		return
	}

	b, err := ioutil.ReadFile(arguments[0])
	if err != nil {
		displayText(strings.Trim(fmt.Sprintf("ioutil.ReadFile: %v\n%s", err,
			prompt), "\n"))
		return
	}

	var c memberlist.Config
	err = json.Unmarshal(b, &c)
	if err != nil {
		displayText(strings.Trim(fmt.Sprintf("json.Unmarshal: %v\n%s", err,
			prompt), "\n"))
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

	if conf == nil {
		displayError("could not create memberlist without configuration")
		return
	}

	// Create will create a new Memberlist using the given configuration.
	// This will not connect to any other node (see Join) yet, but will start
	// all the listeners to allow other nodes to join this memberlist.
	// After creating a Memberlist, the configuration given should not be
	// modified by the user anymore.
	mlist, err = memberlist.Create(conf)
	if err != nil {
		displayError("failed to create memberlist", err)
		return
	}

	displayText(strings.Trim(fmt.Sprintf("local node name: %v\n%s", mlist.LocalNode().Name,
		prompt), "\n"))
}

func showMemberlist(arguments []string) {

	// Get rid off warning
	_ = arguments

	for i, w := range strings.Split(
		strings.Replace(
			strings.Replace(fmt.Sprintf("%#v\n", mlist), "}", "", 1), "{", ": ", 1), " ") {

		displayText(strings.Trim(fmt.Sprintf("%d:\t%v\n", i, w), "\n"))
	}
	displayText(prompt)
}

func showLocalNode(arguments []string) {

	// Get rid off warning
	_ = arguments

	// LocalNode is used to return the local Node
	node := mlist.LocalNode()

	//displayText(strings.Trim(fmt.Sprintf("localNode.Name: %v\n", node.Name), "\n"))
	displayText(strings.Trim(fmt.Sprintf("localNode.Name: %v\n"+
		"localNode.Address: %v\n%s",
		node.Name, node.Address(),
		prompt), "\n"))
}

func listMembers(arguments []string) {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		displayText(strings.Trim(fmt.Sprintf("no memberlist found\n%s", prompt), "\n"))
		return
	}

	if len(mlist.Members()) == 0 {
		displayText(strings.Trim(fmt.Sprintf("empty memberlist\n%s", prompt), "\n"))
		return
	}

	members := fmt.Sprintf("%s\n", prompt)
	for _, m := range mlist.Members() {
		members += fmt.Sprintf("Name: %v %v\n", m.Name, m.Address())
	}
	displayText(members + prompt)
}

func joinMemberlist(arguments []string) {

	if mlist == nil {
		displayText(strings.Trim(fmt.Sprintf("no memberlist found\n%s", prompt), "\n"))
		return
	}

	if len(arguments) > 0 {

		// Join arguments to slice of members
		n, err := mlist.Join(arguments)
		if err != nil {
			displayText(strings.Trim(fmt.Sprintf("join the memberlist failed: %v\n%s", err, prompt), "\n"))
			return
		}
		displayText(strings.Trim(fmt.Sprintf("%d host(s) successfully contacted\n%s", n, prompt), "\n"))
	} else {

		if bootstrapData == nil {
			displayText(strings.Trim(fmt.Sprintf("no bootstrap data found\n%s", prompt), "\n"))
			return
		}

		// No bootstrap server
		if len(bootstrapData.Peers) == 0 {
			displayText(strings.Trim(fmt.Sprintf("Nothing to join\n%s", prompt), "\n"))
			return
		}

		// No other bootstrap server
		if len(bootstrapData.Peers) == 1 {
			if _, ok := bootstrapData.Peers[bootstrapApi.Self.ID]; !ok {
				displayText(strings.Trim(fmt.Sprintf("Nothing to join\n5s", prompt), "\n"))
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
			displayText(strings.Trim(fmt.Sprintf("join bootstrap peers to members failed: %v\n%s", err, prompt), "\n"))
			return
		}
		displayText(strings.Trim(fmt.Sprintf("%d host(s) successfully contacted\n%s", n, prompt), "\n"))
	}
	displayText(strings.Trim(fmt.Sprintf("Known nodes: %v\n%s", mlist.Members(), prompt), "\n"))
}

func leaveMemberlist(arguments []string) {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		displayText(strings.Trim(fmt.Sprintf("no memberlist found\n%s", prompt), "\n"))
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
		displayText(strings.Trim(fmt.Sprintf("leaveMemberlist failed: %v\n%s", err, prompt), "\n"))
		return
	}
	displayText(prompt)
}

// ToDo: Is UpdateNode needed?
func updateMemberlist(arguments []string) {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		displayText(strings.Trim(fmt.Sprintf("no memberlist found\n%s", prompt), "\n"))
		return
	}

	// UpdateNode is used to trigger re-advertising the local node. This is
	// primarily used with a Delegate to support dynamic updates to the local
	// meta data.  This will block until the update message is successfully
	// broadcasted to a member of the cluster, if any exist or until a specified
	// timeout is reached.
	err := mlist.UpdateNode(time.Second)
	if err != nil {
		displayText(strings.Trim(fmt.Sprintf("updateMemberlist failed: %v\n%s", err, prompt), "\n"))
	}
}

func startBroadcast(arguments []string) {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		displayError("cannot start broadcasting without memberlist")
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
	displayText(strings.Trim(fmt.Sprintf("Local member %s:%d\n%s", node.Addr, node.Port, prompt), "\n"))
}

func shutdownBroadcast(arguments []string) {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		displayText(strings.Trim(fmt.Sprintf("no memberlist found\n%s", prompt), "\n"))
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
		displayText(strings.Trim(fmt.Sprintf("failed mlist.Shutdown(): %v\n%s", err, prompt), "\n"))
		return
	}

	displayText(prompt)
}

func shutdownBroadcastTransport(arguments []string) {

	// Get rid off warning
	_ = arguments

	if conf == nil {
		displayText(strings.Trim(fmt.Sprintf("nothing for memberlist configured\n%s", prompt), "\n"))
		return
	}

	if conf.Transport == nil {
		displayText(strings.Trim(fmt.Sprintf("no transport for memberlist configured\n%s", prompt), "\n"))
		return
	}

	err := conf.Transport.Shutdown()
	if err != nil {
		displayText(strings.Trim(fmt.Sprintf("failed Transport.Shutdown(): %v\n%s", err, prompt), "\n"))
	}

}

func getHealthScore(arguments []string) {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		displayText(strings.Trim(fmt.Sprintf("no memberlist found\n"), "\n"))
		return
	}

	displayText(strings.Trim(fmt.Sprintf("Health Score: %d\n", mlist.GetHealthScore()), "\n"))
}

func deleteMemberlist(arguments []string) {

	// Get rid off warning
	_ = arguments

	mlist = nil

	displayText(prompt)
}
