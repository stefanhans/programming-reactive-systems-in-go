package bootstrapApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

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

// Config has the configuration of the bootstrap service.
type Config struct {
	MaxPeers int `json:"maxpeers,omitempty"` // Max number of bootstrap peers to be saved
	MinPeers int `json:"minpeers,omitempty"` // Min number of bootstrap peers, i.e. triggers refill
	NumPeers int `json:"numpeers,omitempty"` // Number of bootstrap peers
}

// BootstrapData is the complete data structure
type BootstrapData struct {
	Config Config
	Peers  map[string]*Peer
}

// BootstrapDataAPI contains all data needed
type BootstrapDataAPI struct {
	ServiceUrl    string
	Self          *Peer
	BootstrapData *BootstrapData
}

// Create returns the individual BootstrapDataAPI initially
func Create(peer *Peer) (*BootstrapDataAPI, error) {

	// https://europe-west1-bootstrap-peers.cloudfunctions.net
	// http://localhost:8080
	serviceUrl := os.Getenv("BOOTSTRAP_DATA_SERVER")
	if serviceUrl == "" {
		return nil, fmt.Errorf("BOOTSTRAP_DATA_SERVER environment variable unset or missing")
	}

	// Send ping request to service
	res, err := http.Post(serviceUrl+"/ping",
		"application/x-www-form-urlencoded",
		strings.NewReader(""))
	if err != nil {
		return nil, fmt.Errorf("failed to ping BOOTSTRAP_DATA_SERVER: %v\n", err)
	}
	fmt.Printf("Received reply from Ping: %v\n", res.Status)
	err = res.Body.Close()
	if err != nil {
		fmt.Printf("cannot close response body\n")
	}

	peer.ID = uuid.NewV4().String()
	peer.Timestamp = fmt.Sprintf("%d", time.Now().Unix())

	return &BootstrapDataAPI{
		ServiceUrl:    serviceUrl,
		Self:          peer,
		BootstrapData: &BootstrapData{},
	}, nil
}

// Join sends the data of a peer to the bootstrap service to join
// and returns the current bootstrap data.
func (api *BootstrapDataAPI) Join() *BootstrapData {

	// Send request to service
	res, err := http.Post(api.ServiceUrl+"/join",
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("%s %s %s %s %s %s %s",
			api.Self.ID, api.Self.Name, api.Self.Ip, api.Self.Port, api.Self.Protocol, "0", api.Self.Timestamp)))
	if err != nil {
		fmt.Printf("failed to join bootstrap-peers (%v): %v\n", api.Self, err)
		return nil
	}
	if res.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("Received no \"200 OK\" from Join: %q\n", strings.TrimSuffix(string(b), "\n"))
		return nil
	}
	fmt.Printf("Received reply from Join: %v\n", res.Status)

	// Read response body in JSON
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Printf("failed to read response of bootstrap-peers join (%v): %v\n", api.Self, err)
		return nil
	}

	// Unmarshall bootstrap data
	var bootstrapData BootstrapData
	err = json.Unmarshal(body, &bootstrapData)
	if err != nil {
		fmt.Printf("failed to unmarshall response of bootstrap-peers join (%v): %v\n", res.Proto, err)
		return nil
	}

	return &bootstrapData
}

// Leave sends the id of a peer which wants to leave the list
// and returns the current bootstrap data.
func (api *BootstrapDataAPI) Leave(id string) *BootstrapData {

	// Send request to service
	res, err := http.Post(api.ServiceUrl+"/leave",
		"application/x-www-form-urlencoded",
		strings.NewReader(id))
	if err != nil {
		fmt.Printf("failed to leave bootstrap-peers (%s): %v\n", id, err)
		return nil
	}
	fmt.Printf("Received reply from Leave: %v\n", res.Status)

	// Read response body in JSON
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Printf("failed to read response of bootstrap-peers leave (%s): %v\n", id, err)
		return nil
	}

	// Unmarshall bootstrap data
	var bootstrapData BootstrapData
	err = json.Unmarshal(body, &bootstrapData)
	if err != nil {
		fmt.Printf("failed to unmarshall response of bootstrap-peers leave (%v): %v\n", res.Proto, err)
		return nil
	}

	return &bootstrapData
}

// Refill sends the data of a peer to the bootstrap service to fill a possible gap
// and returns the current bootstrap data.
func (api *BootstrapDataAPI) Refill() *BootstrapData {

	// Send request to service
	res, err := http.Post(api.ServiceUrl+"/refill",
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("%s %s %s %s %s %s %s",
			api.Self.ID, api.Self.Name, api.Self.Ip, api.Self.Port, api.Self.Protocol, api.Self.Status, api.Self.Timestamp)))
	if err != nil {
		fmt.Printf("failed to join bootstrap-peers (%v): %v\n", api.Self, err)
		return nil
	}
	fmt.Printf("Received reply from Refill: %v\n", res.Status)

	// Read response body in JSON
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Printf("failed to read response of bootstrap-peers refill (%v): %v\n", api.Self, err)
		return nil
	}

	// Unmarshall bootstrap data
	var bootstrapData BootstrapData
	err = json.Unmarshal(body, &bootstrapData)
	if err != nil {
		fmt.Printf("failed to unmarshall response of bootstrap-peers leave (%v): %v\n", res.Proto, err)
		return nil
	}
	return &bootstrapData
}

// List returns the current bootstrap data.
func (api *BootstrapDataAPI) List() *BootstrapData {

	// Send request to service
	res, err := http.Post(api.ServiceUrl+"/list",
		"application/x-www-form-urlencoded",
		strings.NewReader(""))
	if err != nil {
		fmt.Printf("failed to list bootstrap-peers: %v\n", err)
		return nil
	}
	fmt.Printf("Received reply from List: %v\n", res.Status)

	// Read response body in JSON
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Printf("failed to read response of bootstrap-peers list: %v\n", err)
		return nil
	}

	// Unmarshall bootstrap data
	var bootstrapData BootstrapData
	err = json.Unmarshal(body, &bootstrapData)
	if err != nil {
		fmt.Printf("failed to unmarshall response of bootstrap-peers join (%v): %v\n", res.Proto, err)
		return nil
	}

	return &bootstrapData
}

// Reset resets the bootstrap service
// and returns the bootstrap data of the empty list.
func (api *BootstrapDataAPI) Reset() *BootstrapData {

	// Send request to service
	res, err := http.Post(api.ServiceUrl+"/reset",
		"application/x-www-form-urlencoded",
		strings.NewReader(""))
	if err != nil {
		fmt.Printf("failed to reset bootstrap-peers: %v\n", err)
		return nil
	}
	fmt.Printf("Received reply from Reset: %v\n", res.Status)

	// Read response body in JSON
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Printf("failed to read response of bootstrap-peers list: %v\n", err)
		return nil
	}

	// Unmarshall bootstrap data
	var bootstrapData BootstrapData
	err = json.Unmarshal(body, &bootstrapData)
	if err != nil {
		fmt.Printf("failed to unmarshall response of bootstrap-peers Reset (%v): %v\n", res.Proto, err)
		return nil
	}

	return &bootstrapData
}

// UpdateConfig updates the configuration of the service.
func (api *BootstrapDataAPI) UpdateConfig() (*BootstrapData, error) {

	// Send request to service
	res, err := http.Post(api.ServiceUrl+"/config",
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("%d %d",
			api.BootstrapData.Config.MaxPeers, api.BootstrapData.Config.MinPeers)))
	if err != nil {
		return nil, fmt.Errorf("failed to config bootstrap data: %v", err)
	}
	fmt.Printf("Received reply from UpdateConfig: %v\n", res.Status)

	// Read response body in JSON
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read response of bootstrap data update config: %v", err)
	}

	// Unmarshall bootstrap data
	var bootstrapData BootstrapData
	err = json.Unmarshal(body, &bootstrapData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall response of bootstrap data UpdateConfig (%v): %v\n", res.Proto, err)
	}

	return &bootstrapData, nil
}
