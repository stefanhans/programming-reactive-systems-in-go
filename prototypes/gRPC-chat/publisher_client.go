package main

import (
	"fmt"
	"strings"

	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/gRPC-chat/chat-group"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Subscribe sends a subscribe request to the publisher
func Subscribe(name string, ip string, port string) error {

	// Create gRPC client connected with gRPC publisher
	client, err := dialPublisher()
	if err != nil {
		return err
	}

	// Subscribe via gRPC client
	_, err = client.Subscribe(context.Background(), &chatgroup.Member{Name: name, Ip: ip, Port: port})
	if err != nil {
		return fmt.Errorf("could not subscribe to the chatgroup: %v", err)
	}
	return nil
}

// Unsubscribe sends a unsubscribe request to the publisher
func Unsubscribe(name string) error {

	// Create gRPC client connected with gRPC publisher
	client, err := dialPublisher()
	if err != nil {
		return err
	}

	// Unsubscribe via gRPC client
	_, err = client.Unsubscribe(context.Background(), &chatgroup.Member{Name: name})
	if err != nil {
		return fmt.Errorf("could not unsubscribe from the chatgroup: %v", err)
	}
	return nil
}

// Publish send a publish request to the publisher
func Publish(name string, text ...string) error {

	// Create gRPC client connected with gRPC publisher
	client, err := dialPublisher()
	if err != nil {
		return err
	}

	// Prepare message
	msg := chatgroup.Message{Sender: &chatgroup.Member{Name: name}, Text: strings.Join(text[:], " ")}

	// Publish via gRPC client
	_, err = client.Publish(context.Background(), &msg)
	if err != nil {
		return fmt.Errorf("could not publish to the chatgroup: %v", err)
	}
	return nil
}

// Dial gRPC publisher and return gRPC client
func dialPublisher() (chatgroup.PublisherClient, error) {

	conn, err := grpc.Dial(":"+serverPort, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect to publisher: %v", err)
	}
	return chatgroup.NewPublisherClient(conn), nil
}
