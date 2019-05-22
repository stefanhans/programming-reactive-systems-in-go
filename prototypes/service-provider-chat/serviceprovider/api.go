package serviceprovider

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/service-provider-chat/serviceprovider/serviceproviders"
)

// NewServiceProvider is the core constructor
func NewServiceProvider(
	serviceProviderType ServiceProviderType,
	candidateName string,
	candidateIp string,
	candidatePort string,
	providerName string,
	providerIp string,
	providerPort string) (*ServiceProvider, error) {

	// Resolve Ip string of candidate and update accordingly
	addr, err := net.ResolveIPAddr("ip", candidateIp)
	if err != nil {
		return nil, fmt.Errorf("no valid ip address %q for candidate: %v\n", candidateIp, err.Error())
	}
	candidateIp = addr.String()

	// Resolve Ip string of provider and update accordingly
	addr, err = net.ResolveIPAddr("ip", providerIp)
	if err != nil {
		return nil, fmt.Errorf("no valid ip address %q for provider: %v\n", providerIp, err.Error())
	}
	providerIp = addr.String()

	candidate := &serviceproviders.Service{
		Name:   candidateName,
		Ip:     candidateIp,
		Port:   candidatePort,
		Status: serviceproviders.Service_CANDIDATE,
	}

	provider := &serviceproviders.Service{
		Name:   providerName,
		Ip:     providerIp,
		Port:   providerPort,
		Status: serviceproviders.Service_PROVIDER,
	}

	var services []*serviceproviders.Service
	if serviceProviderType == SERVICE {
		services = append(services, &serviceproviders.Service{
			Name:   candidateName,
			Ip:     candidateIp,
			Port:   candidatePort,
			Status: serviceproviders.Service_CANDIDATE,
		})
	}

	internalServiceProvider := &serviceproviders.ServiceProvider{
		Provider: provider,
		Services: services,
	}

	messageHandlerFunctions := make(map[serviceproviders.Message_MessageType]func(*serviceproviders.Message, net.Addr) error)
	replyChannels := make(map[serviceproviders.Message_MessageType]chan *serviceproviders.Message)

	newServiceProvider := &ServiceProvider{
		serviceProviderType: serviceProviderType,
		version:             0,
		candidate:           candidate,
		provider:            provider,
		serviceProvider:     internalServiceProvider,
		message: &serviceproviders.Message{
			MsgType:         serviceproviders.Message_SERVICE_REQUEST,
			ServiceProvider: internalServiceProvider,
			Sender:          candidate,
		},
		messageHandlerFunctions: messageHandlerFunctions,
		replyChannels:           replyChannels,
	}

	messageHandlerFunctions[serviceproviders.Message_SERVICE_REQUEST] = newServiceProvider.handleServiceRequest
	messageHandlerFunctions[serviceproviders.Message_SERVICE_REPLY] = newServiceProvider.handleServiceReply

	replyChannels[serviceproviders.Message_SERVICE_REPLY] = make(chan *serviceproviders.Message)

	return newServiceProvider, nil
}

// GetReplyChannel returns a channel for messages of a certain type
func (serviceProvider *ServiceProvider) GetReplyChannel(replyType serviceproviders.Message_MessageType) chan *serviceproviders.Message {
	return serviceProvider.replyChannels[replyType]
}

// String shows a textual representation of a service provider
func (serviceProvider *ServiceProvider) String() string {
	out := "serviceprovider.ServiceProvider:\n"
	out += fmt.Sprintf("\tServiceProviderType: %v\n", serviceProvider.serviceProviderType)
	out += fmt.Sprintf("\tVersion: %v\n", serviceProvider.version)
	out += fmt.Sprintf("\tCandidate: %v\n", serviceProvider.candidate)
	out += fmt.Sprintf("\tProvider: %v\n", serviceProvider.provider)
	out += fmt.Sprintf("\tServiceProvider: \n")
	out += fmt.Sprintf("\t\tProvider: %v\n", serviceProvider.serviceProvider.Provider)
	for i, service := range serviceProvider.serviceProvider.Services {
		out += fmt.Sprintf("\t\tService[%d]: name:%q ip:%q port:%q status:%v\n",
			i, service.Name, service.Ip, service.Port, service.Status)
	}
	out += fmt.Sprintf("\tMessage:\n")
	out += fmt.Sprintf("\t\tMsgType: %v\n", serviceProvider.message.MsgType)
	out += fmt.Sprintf("\t\tSender: %v\n", serviceProvider.message.Sender)
	out += fmt.Sprintf("\t\tServiceProvider:\n")
	out += fmt.Sprintf("\t\t\tProvider: %v\n", serviceProvider.message.ServiceProvider.Provider)
	for i, service := range serviceProvider.message.ServiceProvider.Services {
		out += fmt.Sprintf("\t\t\tService[%d]: name:%q ip:%q port:%q status:%v\n",
			i, service.Name, service.Ip, service.Port, service.Status)
	}
	out += fmt.Sprintf("\tMessageHandlerFunctions: %v\n", serviceProvider.messageHandlerFunctions)
	out += fmt.Sprintf("\tReplyChannels: %v\n", serviceProvider.replyChannels)

	return out
}

// ServiceProviderType returns the type of the service provider, i.e. PROVIDER or SERVICE
func (serviceProvider *ServiceProvider) ServiceProviderType() ServiceProviderType {
	return serviceProvider.serviceProviderType
}

// Version returns the version of the service list
func (serviceProvider *ServiceProvider) Version() int {
	return serviceProvider.version
}

// CandidateName returns the name of the service candidate
func (serviceProvider *ServiceProvider) CandidateName() string {
	return serviceProvider.candidate.Name
}

// SetCandidateIp sets the Ip address of the service candidate
func (serviceProvider *ServiceProvider) SetCandidateIp(ip string) {
	serviceProvider.candidate.Ip = ip
}

// CandidateIp returns the Ip address of the service candidate
func (serviceProvider *ServiceProvider) CandidateIp() string {
	return serviceProvider.candidate.Ip
}

// SetCandidatePort sets the port number of the service candidate
func (serviceProvider *ServiceProvider) SetCandidatePort(port string) error {

	// Port number is an integer
	p, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	// Within free port number range without root access, i.e. [1024, 65535]
	if p < 1024 || p > 65535 {
		return fmt.Errorf("portnumber %d not between 1024 and 65535", p)
	}

	serviceProvider.candidate.Port = port

	return nil
}

// CandidatePort returns the port number of the service candidate
func (serviceProvider *ServiceProvider) CandidatePort() string {
	return serviceProvider.candidate.Port
}

// SetCandidateStatus sets the status of the service candidate
func (serviceProvider *ServiceProvider) SetCandidateStatus(status serviceproviders.Service_Status) {
	serviceProvider.candidate.Status = status
}

// CandidateStatus returns the status of the service candidate
func (serviceProvider *ServiceProvider) CandidateStatus() serviceproviders.Service_Status {
	return serviceProvider.candidate.Status
}

// ProviderName returns the name of the service provider
func (serviceProvider *ServiceProvider) ProviderName() string {
	return serviceProvider.provider.Name
}

// SetProviderIp sets the Ip address of the provider
func (serviceProvider *ServiceProvider) SetProviderIp(ip string) {
	serviceProvider.provider.Ip = ip
}

// ProviderIp returns the Ip address of the provider
func (serviceProvider *ServiceProvider) ProviderIp() string {
	return serviceProvider.provider.Ip
}

// SetProviderPort sets the port number of the provider
func (serviceProvider *ServiceProvider) SetProviderPort(port string) error {

	// Port number is an integer
	p, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	// Within free port number range without root access, i.e. [1024, 65535]
	if p < 1024 || p > 65535 {
		return fmt.Errorf("portnumber %d not between 1024 and 65535", p)
	}

	serviceProvider.provider.Port = port

	return nil
}

// ProviderPort returns the port number of the provider
func (serviceProvider *ServiceProvider) ProviderPort() string {
	return serviceProvider.provider.Port
}

// Message returns the current message of the provider
func (serviceProvider *ServiceProvider) Message() *serviceproviders.Message {
	return serviceProvider.message
}

// UpdateServices updates the list of services of the service provider
func (serviceProvider *ServiceProvider) UpdateServices(message *serviceproviders.Message) bool {

	// Insert service
	err := serviceProvider.insertService(message)
	if err != nil {
		return false
	}

	// Return, if one service is working
	for _, s := range serviceProvider.serviceProvider.Services {
		if s.Status == serviceproviders.Service_SERVICE {
			return true
		}
	}

	// Elect a service
	serviceProvider.serviceElection()
	return true
}

// StartRootService starts the server providing information and administration of the service
func (serviceProvider *ServiceProvider) StartRootService() {

	// Try to create listener
	listener, err := net.Listen("tcp", net.JoinHostPort(serviceProvider.ProviderIp(), serviceProvider.ProviderPort()))
	if err != nil {
		log.Fatalf("could not listenen on %q: %v", net.JoinHostPort(serviceProvider.ProviderIp(), serviceProvider.ProviderPort()), err)
	}

	// Goroutine for handling requests
	go func(listener net.Listener) {
		defer listener.Close()
		for {
			// Wait for a connection.
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			// Handle the connection in a new goroutine.
			// The loop then returns to accepting, so that
			// multiple connections may be served concurrently.
			go serviceProvider.handleRequest(conn)
		}
	}(listener)
}

// StartClientService starts the client providing the potential service.
func (serviceProvider *ServiceProvider) StartClientService() {

	// Try to create listener
	listener, err := net.Listen("tcp", net.JoinHostPort(serviceProvider.CandidateIp(), ""))
	if err != nil {
		log.Fatalf("could not listenen on %q: %v", net.JoinHostPort(serviceProvider.CandidateIp(), ""), err)
	}

	// Set the port number of the listener as member port
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		fmt.Printf("cannot split host from (new) port %q: %v", listener.Addr().String(), err)
	}
	serviceProvider.message.Sender.Port = port

	// Goroutine for handling requests
	go func(listener net.Listener) {
		defer listener.Close()
		for {
			// Wait for a connection.
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			// Handle the connection in a new goroutine.
			// The loop then returns to accepting, so that
			// multiple connections may be served concurrently.
			go serviceProvider.handleRequest(conn)
		}
	}(listener)
}

func (serviceProvider *ServiceProvider) RequestServices() error {

	serviceProvider.message.MsgType = serviceproviders.Message_SERVICE_REQUEST
	serviceProvider.message.Sender = serviceProvider.candidate

	return tcpSend(serviceProvider.message, net.JoinHostPort(serviceProvider.provider.Ip, serviceProvider.provider.Port))
}

func (serviceProvider *ServiceProvider) GetService() (*serviceproviders.Service, error) {

	for _, s := range serviceProvider.serviceProvider.Services {
		if s.Status == serviceproviders.Service_SERVICE {
			return s, nil
		}
	}
	return nil, fmt.Errorf("no active service found")
}

func (serviceProvider *ServiceProvider) GetServiceAddress() (string, error) {

	for _, s := range serviceProvider.serviceProvider.Services {
		if s.Status == serviceproviders.Service_SERVICE {
			return net.JoinHostPort(s.Ip, s.Port), nil
		}
	}
	return "", fmt.Errorf("no active service found")
}
