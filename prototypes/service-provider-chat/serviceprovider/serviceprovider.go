package serviceprovider

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/service-provider-chat/serviceprovider/serviceproviders"
)

type ServiceProviderType int32

const (
	PROVIDER ServiceProviderType = 0 // PROVIDER is a root service and has to be known explicitly
	SERVICE  ServiceProviderType = 1 // SERVICE can be requested from a PROVIDER
)

type ServiceProvider struct {

	// serviceProviderType can be PROVIDER or SERVICE
	serviceProviderType ServiceProviderType

	// version shows the current version
	version int

	// candidate is oneself (as possible working service)
	candidate *serviceproviders.Service

	// provider replies the request for a working service or provider, respectively
	provider *serviceproviders.Service

	// ServiceProvider stores the provider and its list of services (candidates and one working)
	serviceProvider *serviceproviders.ServiceProvider

	// message to send and functions to handle replies
	message                 *serviceproviders.Message
	messageHandlerFunctions map[serviceproviders.Message_MessageType]func(*serviceproviders.Message, net.Addr) error

	// replyChannels store the channels to receive the message replies
	replyChannels map[serviceproviders.Message_MessageType]chan *serviceproviders.Message
}

func (serviceProviderType ServiceProviderType) String() string {
	switch serviceProviderType {
	case PROVIDER:
		return "PROVIDER"
	case SERVICE:
		return "SERVICE"
	}
	return ""
}

// handleRequest reads all incoming data, take the leading byte as message type,
// and use the appropriate handler for the rest
func (serviceProvider *ServiceProvider) handleRequest(conn net.Conn) {

	defer conn.Close()

	// Read all data from the connection
	var buf [512]byte
	var data []byte
	addr := conn.RemoteAddr()

	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			break
		}
		data = append(data, buf[0:n]...)
	}

	// Unmarshall message
	var msg serviceproviders.Message
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		fmt.Printf("could not unmarshall leadergroup.Message: %v\n", err)
		return
	}

	fmt.Printf("Message: %v\n", msg)

	// Fetch the handler from a map by the message type and call it accordingly
	if replyAction, ok := serviceProvider.messageHandlerFunctions[msg.MsgType]; ok {
		err := replyAction(&msg, conn.RemoteAddr())
		if err != nil {
			fmt.Printf("could not handle request %v from %v: %v", msg.MsgType, addr, err)
		}
	} else {
		log.Printf("server: unknown message type %v\n", msg.MsgType)
	}
}

// handleServiceRequest
func (serviceProvider *ServiceProvider) handleServiceRequest(message *serviceproviders.Message, addr net.Addr) error {

	// Update sender Ip
	message.Sender.Ip = strings.Split(addr.String(), ":")[0]

	fmt.Printf("handleServiceRequest\n")
	// Insert or update sender in services and send reply, if so
	if serviceProvider.UpdateServices(message) {

		message.ServiceProvider = serviceProvider.serviceProvider
		message.MsgType = serviceproviders.Message_SERVICE_REPLY

		serviceProvider.version++

		fmt.Printf("UpdateServices: version %d\n", serviceProvider.version)

		return tcpSend(message, net.JoinHostPort(message.Sender.Ip, message.Sender.Port))
	}

	return tcpSend(message, net.JoinHostPort(message.Sender.Ip, message.Sender.Port))
}

// handleServiceReply
func (serviceProvider *ServiceProvider) handleServiceReply(message *serviceproviders.Message, addr net.Addr) error {

	serviceProvider.serviceProvider.Services = message.ServiceProvider.Services
	serviceProvider.version++

	return nil
}

// tcpSend
func tcpSend(message *serviceproviders.Message, recipient string) error {

	fmt.Printf("Message: %v\n", message)

	// Connect to the recipient
	conn, err := net.Dial("tcp", recipient)
	if err != nil {
		return fmt.Errorf("could not connect to recipient %q: %v", recipient, err)
	}

	// Marshal into binary format
	byteArray, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not encode message: %v", err)
	}

	// Write the bytes to the connection
	_, err = conn.Write(byteArray)
	if err != nil {
		return fmt.Errorf("could not write message to the connection: %v", err)
	}

	// Close connection
	return conn.Close()
}
