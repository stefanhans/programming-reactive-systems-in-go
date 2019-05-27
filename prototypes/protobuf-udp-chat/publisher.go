package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"syscall"

	"github.com/golang/protobuf/proto"
	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/protobuf-udp-chat/chat-group"
)

// Start publisher service to provide member registration and message publishing
func startPublisher() error {

	// Create publishingListener
	publishingListener, err := net.ListenPacket("udp", publishingService)
	log.Printf("LISTENER: %q\n", publishingService)

	if err != nil {

		// Check if publisher error is "address already in use"
		if strings.Contains(err.Error(), syscall.EADDRINUSE.Error()) {

			// Subscribe at already running Publisher
			err = Subscribe()
			if err != nil {
				log.Fatalf("failed to subscribe at already running Publisher: %v", err)
			}

			// Append text messages in "messages" view of subscriber
			displayText(fmt.Sprintf("<%s (%s:%s) has joined>", selfMember.Name, selfMember.Ip, selfMember.Port))

			return nil
		}

		// Exit on unexpected error
		log.Fatalf("could not listen to %q: %v\n", publishingService, err)
	}
	defer publishingListener.Close()

	log.Printf("Started publishing service listening on %q\n", publishingService)

	// Append text messages in "messages" view of publisher
	displayText(fmt.Sprintf("<publishing service running: %s (%s:%s)>", selfMember.Name, serverIp, serverPort))

	// Resolve Ip string and update accordingly
	addr, err := net.ResolveIPAddr("ip", selfMember.Ip)
	if err != nil {
		fmt.Printf("no valid ip address of client %q for publishing service: %v\n", selfMember.Ip, err.Error())
		os.Exit(1)
	}
	selfMember.Ip = addr.String()

	// Subscribe directly at started publishing service
	selfMember.Leader = true
	cgMember = append(cgMember, selfMember)
	log.Printf("Subscribed directly at started publishing service: %v\n", cgMember[0])

	// Append text messages in "messages" view of publisher
	displayText(fmt.Sprintf("<%s (%s:%s) has joined>", selfMember.Name, selfMember.Ip, selfMember.Port))

	buffer := make([]byte, bufferSize)

	// Endless loop in foreground of goroutine
	for {
		n, addr, err := publishingListener.ReadFrom(buffer)

		if err != nil {
			log.Printf("cannot read from buffer: %v\n", err)
		} else {
			//log.Printf("Read %v bytes from %v: %v\n", n, addr, buffer)
			go func(buffer []byte, addr net.Addr) {
				handlePublisherRequest(buffer, addr)

			}(buffer[:n], addr)
		}
	}

	return nil
}

// Read all incoming data, take the message type,
// and use the appropriate handler for the rest
func handlePublisherRequest(data []byte, addr net.Addr) {

	// Unmarshall message
	var msg chatgroup.Message
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		fmt.Printf("could not unmarshall message: %v", err)
	}

	// Fetch the handler from a map by the message type and call it accordingly
	if requestAction, ok := requestActionMap[msg.MsgType]; ok {
		log.Printf("%v\n", msg)
		err := requestAction(&msg, addr)
		if err != nil {
			fmt.Printf("could not handle %v from %v: %v", msg.MsgType, addr, err)
		}
	} else {
		log.Printf("publisher: unknown message type %v\n", msg.MsgType)
	}
}

func handleSubscribeRequest(message *chatgroup.Message, addr net.Addr) error {

	// Update remote Ip address, if changed
	updateRemoteIP(message, addr)

	// Check subscriber for uniqueness
	for _, recipient := range cgMember {
		if recipient.Name == message.Sender.Name {
			return fmt.Errorf("name %q already used", message.Sender.Name)
		}
		if recipient.Ip == message.Sender.Ip && recipient.Port == message.Sender.Port {
			return fmt.Errorf("address %s:%s already used by %s", recipient.Ip, recipient.Port, recipient.Name)
		}
	}

	// Add subscriber
	log.Printf("Add subscriber: %v\n", message.Sender)
	cgMember = append(cgMember, message.Sender)
	log.Printf("Current members registered: %v\n", cgMember)

	err := publishMessage(message, chatgroup.Message_SUBSCRIBE_REPLY)
	if err != nil {
		return fmt.Errorf("Failed to publish Message_SUBSCRIBE_REPLY: %v\n", err)
	}

	return nil
}

func handleUnsubscribeRequest(message *chatgroup.Message, addr net.Addr) error {

	log.Printf("Unregister: %v\n", message.Sender)

	// Remove subscriber
	for i, s := range cgMember {
		if s.Name == message.Sender.Name {
			cgMember = append(cgMember[:i], cgMember[i+1:]...)
			break
		}
	}
	log.Printf("Current members registered: %v\n", cgMember)

	err := publishMessage(message, chatgroup.Message_UNSUBSCRIBE_REPLY)
	if err != nil {
		return fmt.Errorf("Failed to publish Message_UNSUBSCRIBE_REPLY: %v\n", err)
	}

	return nil
}

func handlePublishRequest(message *chatgroup.Message, addr net.Addr) error {

	// Update remote Ip address, if changed
	updateRemoteIP(message, addr)

	log.Printf("Publish from %v: %q\n", message.Sender.Name, message.Text)

	err := publishMessage(message, chatgroup.Message_PUBLISH_REPLY)
	if err != nil {
		return fmt.Errorf("Failed to publish Message_Message_PUBLISH_REPLY: %v\n", err)
	}

	return nil
}

func updateRemoteIP(msg *chatgroup.Message, addr net.Addr) {

	// Check remote Ip address change of message
	if msg.Sender.Ip != strings.Split(addr.String(), ":")[0] {
		log.Printf("Remote Ip address update from %v to %v\n", msg.Sender.Ip, strings.Split(addr.String(), ":")[0])
		msg.Sender.Ip = strings.Split(addr.String(), ":")[0]
	}
}

// Publish a message to all members except the sender
func publishMessage(message *chatgroup.Message, msgType chatgroup.Message_MessageType) error {

	// Set the reply message type
	message.MsgType = msgType

	// Forward message to other chat group members
	for _, recipient := range cgMember {

		// Exclude sender
		if recipient.Name != message.Sender.Name {

			// Send message to recipient
			log.Printf("From %s to %s (%s:%s): %q\n",
				message.Sender.Name, recipient.Name, recipient.Ip, recipient.Port, message.Sender)
			err := sendMessage(message, recipient.Ip+":"+recipient.Port)
			if err != nil {
				return fmt.Errorf("Failed send reply: %v\n", err)
			}
		}
	}
	return nil
}

// Send reply to the sender of the message
func sendMessage(message *chatgroup.Message, recipient string) error {

	// Connect to the recipient
	conn, err := net.Dial("udp", recipient)
	if err != nil {
		return fmt.Errorf("could not connect to recipient %q: %v", recipient, err)
	}

	// Marshal into binary format
	byteArray, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not encode message: %v", err)
	}

	// Write the bytes to the connection
	n, err := conn.Write(byteArray)
	if err != nil {
		return fmt.Errorf("could not write message to the connection: %v", err)
	}
	log.Printf("Message (%v byte) sent (%v byte): %v\n", len(byteArray), n, message)

	// Close connection
	return conn.Close()
}
