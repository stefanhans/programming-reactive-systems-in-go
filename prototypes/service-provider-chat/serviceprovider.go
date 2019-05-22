package main

import (
	"fmt"
	"os"

	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/service-provider-chat/serviceprovider"
)

func main() {

	// Create pure service provider without being a client as well
	rootServiceProvider, err := serviceprovider.NewServiceProvider(
		serviceprovider.PROVIDER, "", "", "",
		"rootService", "127.0.0.1", "22365")

	// Exit on error
	if err != nil {
		fmt.Printf("could not create new service provider: %v", err)
		os.Exit(2)
	}

	// Start only server
	rootServiceProvider.StartRootService()

	// Keep the service running - stop via Ctrl-C
	select {}

}
