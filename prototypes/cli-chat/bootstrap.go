package main

import (
	"fmt"

	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/cli-chat/bootstrap-data-api"
)

var (
	bootstrapApi  *bootstrap_data_api.BootstrapDataAPI
	bootstrapData *bootstrap_data_api.BootstrapData
)

func initializeBootstrapApi() {

	if mlist == nil {
		fmt.Printf("cannot create bootstrap API without local node data of memberlist\n")
		return
	}

	bootstrapApi, err = bootstrap_data_api.Create(&bootstrap_data_api.Peer{
		Name:     mlist.LocalNode().Name,
		Ip:       mlist.LocalNode().Addr.String(),
		Port:     fmt.Sprintf("%d", mlist.LocalNode().Port),
		Protocol: "tcp",
	})
	if err != nil {
		fmt.Printf("bootstrap peer API creation failed: %v\n", err)
	}
}

func joinBootstrapPeers(arguments []string) {

	// Get rid off warning
	_ = arguments

	if bootstrapApi == nil {
		initializeBootstrapApi()
		if bootstrapApi == nil {
			return
		}
	}

	bootstrapData = bootstrapApi.Join()
}

func leaveBootstrapPeers(arguments []string) {

	// Get rid off warning
	_ = arguments

	if bootstrapApi == nil {
		initializeBootstrapApi()
		if bootstrapApi == nil {
			return
		}
	}

	bootstrapData = bootstrapApi.Leave(bootstrapApi.Self.ID)
}

func refillBootstrapPeers(arguments []string) {

	// Get rid off warning
	_ = arguments

	if bootstrapApi == nil {
		initializeBootstrapApi()
		if bootstrapApi == nil {
			return
		}
	}

	// Todo: Think about to avoid unnecessary refill (via memberlist ping and timestamp comparison)
	bootstrapData = bootstrapApi.Refill()
	fmt.Printf("%v\n", bootstrapData.Config)
}

func listBootstrapPeers(arguments []string) {

	// Get rid off warning
	_ = arguments

	if bootstrapApi == nil {
		initializeBootstrapApi()
		if bootstrapApi == nil {
			return
		}
	}

	bootstrapData = bootstrapApi.List()

	for k, v := range bootstrapData.Peers {
		fmt.Printf("%q: %v\n", k, v)
	}
}

func resetBootstrapPeers(arguments []string) {

	// Get rid off warning
	_ = arguments

	if bootstrapApi == nil {
		initializeBootstrapApi()
		if bootstrapApi == nil {
			return
		}
	}

	bootstrapApi.Reset()
}

func listLocalBootstrapPeers(arguments []string) {

	// Get rid off warning
	_ = arguments

	for k, v := range bootstrapData.Peers {
		fmt.Printf("%q: %v\n", k, v)
	}
}
