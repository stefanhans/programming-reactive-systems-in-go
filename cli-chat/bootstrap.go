package main

import (
	"encoding/json"
	"fmt"

	bootstrap_data_api "github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-api"
)

var (
	bootstrapApi  *bootstrap_data_api.BootstrapDataAPI
	bootstrapData *bootstrap_data_api.BootstrapData
)

func displayBootstrapData(bootstrapData *bootstrap_data_api.BootstrapData) {

	//
	bootstrapDataJson, err := json.MarshalIndent(bootstrapData, "", "  ")
	if err != nil {
		displayError("failed to marshall bootstrap peers", err)
		return
	}

	displayJson(bootstrapDataJson)
}

func displayBootstrapPeers(bootstrapData *bootstrap_data_api.BootstrapData) {

	bootstrapPeersJson, err := json.MarshalIndent(bootstrapData.Peers, "", "  ")
	if err != nil {
		displayError("failed to marshall bootstrap peers", err)
		return
	}

	displayJson(bootstrapPeersJson)
}

func displayBootstrapConfig(bootstrapData *bootstrap_data_api.BootstrapData) {

	bootstrapConfigJson, err := json.MarshalIndent(bootstrapData.Config, "", "  ")
	if err != nil {
		displayError("failed to marshall bootstrap peers", err)
		return
	}

	displayJson(bootstrapConfigJson)
}

func initializeBootstrapApi() {

	if mlist == nil {
		displayError("cannot create bootstrap API without local node data of memberlist")
		return
	}

	bootstrapApi, err = bootstrap_data_api.Create(&bootstrap_data_api.Peer{
		Name:     mlist.LocalNode().Name,
		Ip:       mlist.LocalNode().Addr.String(),
		Port:     fmt.Sprintf("%d", mlist.LocalNode().Port),
		Protocol: "tcp",
	})
	if err != nil {
		displayError("failed to create bootstrap peer API", err)
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

	displayBootstrapData(bootstrapData)
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

	displayBootstrapData(bootstrapData)
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

	bootstrapData = bootstrapApi.Refill()

	displayBootstrapData(bootstrapData)
}

func showBootstrapData(arguments []string) {

	// Get rid off warning
	_ = arguments

	if bootstrapApi == nil {
		initializeBootstrapApi()
		if bootstrapApi == nil {
			return
		}
	}

	bootstrapData = bootstrapApi.List()

	displayBootstrapData(bootstrapData)
}

func listBootstrapConfig(arguments []string) {

	// Get rid off warning
	_ = arguments

	if bootstrapApi == nil {
		initializeBootstrapApi()
		if bootstrapApi == nil {
			return
		}
	}

	bootstrapData = bootstrapApi.List()

	displayBootstrapConfig(bootstrapData)
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

	displayBootstrapPeers(bootstrapData)
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

	bootstrapData = bootstrapApi.Reset()

	displayBootstrapData(bootstrapData)
}

func listLocalBootstrapPeers(arguments []string) {

	// Get rid off warning
	_ = arguments

	displayBootstrapPeers(bootstrapData)
}
