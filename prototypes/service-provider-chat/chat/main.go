package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/service-provider-chat/chat/chat-group"
	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/service-provider-chat/serviceprovider"
)

func main() {

	// Check command args and set own chatgroup.Member
	flag.Parse()
	if flag.NArg() < 5 {
		fmt.Fprintln(os.Stderr, "missing parameter: <name> <sp ip> <sp port> <chat ip> <chat port>")
		os.Exit(1)
	}

	config = &Config{
		Name: flag.Arg(0),
		ServiceProvider: &ServiceConfig{
			Name:     "root",
			Protocol: "tcp",
			Ip:       "127.0.0.1",
			Port:     "22365",
		},
		ChatService: &ServiceConfig{
			Name:     flag.Arg(0),
			Protocol: "tcp",
			Ip:       flag.Arg(1), // <sp ip>
			Port:     flag.Arg(2), // <sp port>
		},
		ChatListener: &ServiceConfig{
			Name:     flag.Arg(0),
			Protocol: "tcp",
			Ip:       flag.Arg(3), // <chat ip>
			Port:     flag.Arg(4), // <chat port>
		},
	}

	// Prepare logfile for logging
	year, month, day := time.Now().Date()
	hour, minute, second := time.Now().Clock()
	logfilename = fmt.Sprintf("rudimentary-chat-tcp-%s-%v%02d%02d%02d%02d%02d.log", config.Name,
		year, int(month), int(day), int(hour), int(minute), int(second))

	//var err error
	logfile, err := os.OpenFile(logfilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening logfile %v: %v", logfilename, err)
	}
	defer logfile.Close()

	// Config logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("DEBUG: ")

	// Switch logging to logfile
	log.SetOutput(logfile)

	log.Printf("config: %v\n", config)

	// Start displaying service
	go func() {

		//displayingService = config.ChatListenerAddress()

		err = startDisplayer()
		if err != nil {
			log.Fatalf("Failed to start displaying service on %q: %v", displayingService, err)
		}
	}()

	selfMember = &chatgroup.Member{Name: config.ChatListener.Name, Ip: config.ChatListener.Ip, Port: config.ChatListener.Port, Leader: false}

	// Initialize 'serviceprovider' candidate for chat service
	clientServiceProvider, err := serviceprovider.NewServiceProvider(
		serviceprovider.SERVICE, config.Name, config.ChatService.Ip, config.ChatService.Port,
		config.ServiceProvider.Name, config.ServiceProvider.Ip, config.ServiceProvider.Port)
	if err != nil {
		log.Fatalf("could not create new service provider: %v", err)
	}
	clientServiceProvider.StartClientService()

	// Send request for chat service to service provider
	clientServiceProvider.RequestServices()

	// Set the version and wait for the reply increased the version
	currentVersion := clientServiceProvider.Version()

	for clientServiceProvider.Version() == currentVersion {
		time.Sleep(time.Millisecond * 10)
	}

	// Get current chat service and save it in the local configuration
	service, err := clientServiceProvider.GetService()
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
	config.SetChatService(service.Name, service.Ip, service.Port)
	log.Printf("ChatServiceAddress: %v\n", config.ChatServiceAddress())

	//time.Sleep(time.Second * 10)

	// Initialize chat command usage
	commandUsageInit()

	// Create the TUI
	clientGui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Printf("could not create tui: %v\n", err)
	}
	defer clientGui.Close()

	// Try to start publishing service and subscribe accordingly
	go func() {

		err := startPublisher()

		// Check if Publisher is "already in use"
		if err != nil && strings.Contains(err.Error(), syscall.EADDRINUSE.Error()) {

			// Subscribe to the already running publishing service
			err = Subscribe()
			if err != nil {
				log.Fatalf("Failed to subscribe to running publishing service: %v", err)
			}
			log.Printf("Subscribed to the already running publishing service\n")
		}
	}()

	// Start text-based UI
	err = runTUI()
	if err != nil {
		log.Fatalf("runTUI: %v", err)
	}
}
