package main

// File for the collection
var collectionFileName string = "peers.json"

// Peer is the struct for the collection
type Peer struct {
	ID       string `json:"id,omitempty"`   // UUID
	Name     string `json:"name,omitempty"` // chat name
	Ip       string `json:"ip,omitempty"`
	Port     string `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"` // "tcp" or "udp"
	// todo get rid of unused field status
	Status    string `json:"status,omitempty"`
	Timestamp string `json:"timestamp,omitempty"` // Unix time in seconds
}

type Config struct {
	MaxPeers            int `json:"maxpeers,omitempty"`            // Max number of bootstrap peers to be saved
	MinRefillCandidates int `json:"minrefillcandidates,omitempty"` // Number used to decide peer send refill request
	NumPeers            int `json:"numpeers,omitempty"`            // Number of bootstrap peers
}

// BootstrapPeers is a complete data structure
type BootstrapPeers struct {
	Config Config
	Peers  map[string]*Peer
}

// Default config
var bootstrapData = BootstrapPeers{
	Config: Config{
		MaxPeers:            2,
		MinRefillCandidates: 2,
		NumPeers:            0,
	},
	Peers: map[string]*Peer{},
}
