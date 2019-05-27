package main

import (
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/service-provider-chat/chat/chat-group"
)

// Subscribe sends a subscribe request to the publisher
func Subscribe() error {

	newMember := &chatgroup.Message{
		MsgType: chatgroup.Message_SUBSCRIBE_REQUEST,
		Sender:  selfMember}

	return sendPublisherRequest(newMember)
}

// Unsubscribe sends a unsubscribe request to the publisher
func Unsubscribe(memberName string) error {

	leavingMember := &chatgroup.Message{
		MsgType: chatgroup.Message_UNSUBSCRIBE_REQUEST,
		Sender: &chatgroup.Member{
			Name: memberName}}

	return sendPublisherRequest(leavingMember)
}

// Publish send a publish request to the publisher
func Publish(text string) error {

	message := &chatgroup.Message{
		MsgType: chatgroup.Message_PUBLISH_REQUEST,
		Sender:  selfMember,
		Text:    text}

	// Append text message in "messages" view
	displayText(fmt.Sprintf("%s: %s", selfMember.Name, message.Text))

	return sendPublisherRequest(message)
}

// Dial publisher and return connection
func sendPublisherRequest(message *chatgroup.Message) error {

	// Connect to publishing service
	conn, err := net.Dial("tcp", config.ChatServiceAddress())
	if err != nil {
		return fmt.Errorf("could not connect to publishing service: %v", err)
	}

	// Marshal into binary format
	byteArray, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not encode message: %v", err)
	}

	// Write message into connection
	n, err := conn.Write(byteArray)
	if err != nil {
		return fmt.Errorf("could not write message: %v", err)
	}
	log.Printf("Message (%v byte) sent (%v byte): %v\n", len(byteArray), n, message)

	// Close connection
	return conn.Close()
}
