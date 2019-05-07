package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Reset(w http.ResponseWriter, r *http.Request) {

	// Get rid off warning
	_ = r

	// Reset the bootstrap data
	for k := range bootstrapData.Peers {
		delete(bootstrapData.Peers, k)
	}
	bootstrapData.Config.NumPeers = 0

	// Marshall all bootstrap peers
	bootstrapPeersJson, err := json.MarshalIndent(bootstrapData, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshall bootstrap peers for %q: %v", collectionFileName, err),
			http.StatusInternalServerError)
		return
	}

	// Save JSON to file
	mtx.Lock()
	ioutil.WriteFile(collectionFileName, append(bootstrapPeersJson, byte('\n')), 0600)
	mtx.Unlock()

	// Send JSON as response
	fmt.Fprintf(w, "%v", string(bootstrapPeersJson))

}

/*

curl http://localhost:8080/reset

*/
