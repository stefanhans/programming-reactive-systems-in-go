package chat

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

func configureMemberlist(arguments []string) bool {

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
				displayError("could not configure memberlist", err)
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

		displayError("could not return hostname from OS", err)
		return false
	}
	id := uuid.NewUUID().String()
	conf.Name = name + "-" + hostname + "-" + id

	// Prepare logfile for memberlist
	year, month, day := time.Now().Date()
	hour, minute, second := time.Now().Clock()
	mlogfile = fmt.Sprintf("memberlist-%s-%v%02d%02d%02d%02d%02d.log", name,
		year, int(month), int(day), int(hour), int(minute), int(second))

	// Set logger (with output to logfile parameter)
	w := io.Writer(os.Stderr)
	if len(mlogfile) > 0 {
		// Todo logdir as env variable
		f, err := os.Create("log/" + mlogfile)
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

	return true
}

func showMemberlistConfiguration(arguments []string) {

	// Get rid off warning
	_ = arguments

	for i, w := range strings.Split(
		strings.Replace(
			strings.Replace(fmt.Sprintf("%#v\n", conf), "}", "", 1), "{", ": ", 1), " ") {

		displayText(strings.Trim(fmt.Sprintf("%d:\t%v\n%s", i, w,
			prompt), "\n"))
	}

	displayText(prompt)
}

func saveMemberlistConfiguration(arguments []string) bool {

	// Marshal array of struct
	byteArray, err := json.Marshal(conf)
	if err != nil {
		displayError("could not marshall configuration", err)
		return false
	}

	var filename string

	if len(arguments) > 0 {
		filename = arguments[0]
	} else {
		filename = conf.Name + ".json"
	}

	//Formated JSON
	var out bytes.Buffer
	err = json.Indent(&out, byteArray, "", "\t")
	if err != nil {
		displayError("could not indent", err)
	}

	// Replace '{}' with 'null' in JSON string
	str := strings.Replace(string(out.Bytes()), "{}", "null", -1)

	err = ioutil.WriteFile(filename, append([]byte(str), byte('\n')), 0600)
	if err != nil {
		displayError("could not write to file", err)
	}
	log.Printf("%s\n", filename)

	return true
}

func loadMemberlistConfiguration(arguments []string) bool {

	if len(arguments) == 0 {
		displayError("no filename to load specified", err)
		return false
	}

	b, err := ioutil.ReadFile(arguments[0])
	if err != nil {
		displayError("could not read file", err)
		return false
	}

	var c memberlist.Config
	err = json.Unmarshal(b, &c)
	if err != nil {
		displayError("could not unmarshall configuration", err)
		return false
	}
	conf = &c

	conf.Events = &eventDelegate{}

	// Prepare logfile for memberlist
	year, month, day := time.Now().Date()
	hour, minute, second := time.Now().Clock()
	mlogfile = fmt.Sprintf("memberlist-%s-%v%02d%02d%02d%02d%02d.log", name,
		year, int(month), int(day), int(hour), int(minute), int(second))

	// Set logger (with output to logfile parameter)
	w := io.Writer(os.Stderr)
	if len(mlogfile) > 0 {
		// Todo logdir as env variable
		f, err := os.Create("log/" + mlogfile)
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
	return true
}

func createMemberlist(arguments []string) bool {

	// Get rid off warning
	_ = arguments

	if conf == nil {
		displayError("could not create memberlist without configuration")
		return false
	}

	// Create will create a new Memberlist using the given configuration.
	// This will not connect to any other node (see Join) yet, but will start
	// all the listeners to allow other nodes to join this memberlist.
	// After creating a Memberlist, the configuration given should not be
	// modified by the user anymore.
	mlist, err = memberlist.Create(conf)
	if err != nil {
		displayError("failed to create memberlist", err)
		return false
	}

	displayText(strings.Trim(fmt.Sprintf("local node name: %v\n%s", mlist.LocalNode().Name,
		prompt), "\n"))
	return true
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

func listMembers(arguments []string) bool {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		displayError("no memberlist found", err)
		return false
	}

	if len(mlist.Members()) == 0 {
		displayError("empty memberlist", err)
		return false
	}

	members := fmt.Sprintf("%s\n", prompt)
	for _, m := range mlist.Members() {
		members += fmt.Sprintf("Name: %v %v\n", m.Name, m.Address())
	}
	displayText(members + prompt)
	return true
}

func joinMemberlist(arguments []string) bool {

	if mlist == nil {
		displayError("no memberlist found", err)
		return false
	}

	if len(arguments) > 0 {

		// Join arguments to slice of members
		n, err := mlist.Join(arguments)
		if err != nil {
			displayError("could not join the memberlist", err)
			return false
		}
		displayText(strings.Trim(fmt.Sprintf("%d host(s) successfully contacted\n%s", n, prompt), "\n"))
	} else {

		if bootstrapData == nil {
			displayError("no bootstrap data found", err)
			return false
		}

		// No bootstrap server
		if len(bootstrapData.Peers) == 0 {
			displayError("nothing to join", err)
			return false
		}

		// No other bootstrap server
		if len(bootstrapData.Peers) == 1 {
			if _, ok := bootstrapData.Peers[bootstrapApi.Self.ID]; !ok {
				displayError("nothing to join", err)
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
		displayText(strings.Trim(fmt.Sprintf("%d host(s) successfully contacted\n%s", n, prompt), "\n"))
	}
	displayText(strings.Trim(fmt.Sprintf("Known nodes: %v\n%s", mlist.Members(), prompt), "\n"))

	return true
}

func leaveMemberlist(arguments []string) bool {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		displayError("no memberlistfound", err)
		return false
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
		displayError("failed to leave the memberlist", err)
		return false
	}
	displayText(prompt)

	return true
}

// ToDo: Is UpdateNode needed?
func updateMemberlist(arguments []string) bool {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		displayError("no memberlist found", err)
		return false
	}

	// UpdateNode is used to trigger re-advertising the local node. This is
	// primarily used with a Delegate to support dynamic updates to the local
	// meta data.  This will block until the update message is successfully
	// broadcasted to a member of the cluster, if any exist or until a specified
	// timeout is reached.
	err := mlist.UpdateNode(time.Second)
	if err != nil {
		displayError("failed to update the memberlist", err)
		return false
	}
	return true
}

func startBroadcast(arguments []string) bool {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		displayError("cannot start broadcasting without memberlist")
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
	displayText(strings.Trim(fmt.Sprintf("Local member %s:%d\n%s", node.Addr, node.Port, prompt), "\n"))

	return true
}

func shutdownBroadcast(arguments []string) bool {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		displayError("no memberlist found", err)
		return false
	}

	// Shutdown will stop any background maintenance of network activity
	// for this memberlist, causing it to appear "dead". A leave message
	// will not be broadcasted prior, so the cluster being left will have
	// to detect this node's shutdown using probing. If you wish to more
	// gracefully exit the cluster, call Leave prior to shutting down.
	//
	// This method is safe to call multiple times.
	err = mlist.Shutdown()
	if err != nil {
		displayError("failed to shutdown the memberlist", err)
		return false
	}

	displayText(prompt)

	return true
}

func shutdownBroadcastTransport(arguments []string) bool {

	// Get rid off warning
	_ = arguments

	if conf == nil {
		displayError("no configuration for the memberlist found", err)
		return false
	}

	if conf.Transport == nil {
		displayError("no configuration for memberlist transport found", err)
		return false
	}

	err := conf.Transport.Shutdown()
	if err != nil {
		displayError("failed to shutdown memberlist transport", err)
		return false
	}
	return true
}

func getHealthScore(arguments []string) bool {

	// Get rid off warning
	_ = arguments

	if mlist == nil {
		displayError("no memberlist found", err)
		return false
	}

	displayText(strings.Trim(fmt.Sprintf("Health Score: %d\n", mlist.GetHealthScore()), "\n"))
	return true
}

func deleteMemberlist(arguments []string) {

	// Get rid off warning
	_ = arguments

	mlist = nil

	displayText(prompt)
}
