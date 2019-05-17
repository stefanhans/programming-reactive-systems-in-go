package bootstrap_data_api

import (
	"net/http"
	"os"
	"strings"
	"testing"
)

func executeApiFunc(api *BootstrapDataAPI, cmd string, arguments ...string) *BootstrapData {

	switch cmd {
	case "join":
		if len(arguments) != 0 {
			api.Self.ID = arguments[0]
		}
		return api.Join()
	case "leave":
		return api.Leave(arguments[0])
	case "refill":
		return api.Refill()
	case "reset":
		return api.Reset()
	}

	return nil
}

func TestLocalhostPing(t *testing.T) {
	serviceUrl := os.Getenv("BOOTSTRAP_DATA_SERVER")
	if serviceUrl != "http://localhost:8080" {
		t.Errorf("BOOTSTRAP_DATA_SERVER environment variable expected %q set %q",
			"http://localhost:8080", serviceUrl)
	}

	// Send ping request to service
	res, err := http.Post(serviceUrl+"/ping",
		"application/x-www-form-urlencoded",
		strings.NewReader(""))
	if err != nil {
		t.Errorf("failed to ping BOOTSTRAP_DATA_SERVER: %v\n", err)
		return
	}
	err = res.Body.Close()
	if err != nil {
		t.Errorf("cannot close response body\n")
	}

	//Send ping request to wrong service
	serviceUrl = "http://unknown"
	res, err = http.Post(serviceUrl+"/ping",
		"application/x-www-form-urlencoded",
		strings.NewReader(""))
	if err == nil {
		t.Errorf("unexpected success to ping unknown BOOTSTRAP_DATA_SERVER\n")
		err = res.Body.Close()
		if err != nil {
			t.Errorf("cannot close response body\n")
		}
		return
	}

}

func TestLocalhostReset(t *testing.T) {
	serviceUrl := os.Getenv("BOOTSTRAP_DATA_SERVER")
	if serviceUrl != "http://localhost:8080" {
		t.Errorf("BOOTSTRAP_DATA_SERVER environment variable expected %q set %q",
			"http://localhost:8080", serviceUrl)
	}

	dummy, err := Create(&Peer{})
	if err != nil {
		t.Errorf("cannot create dummy API: %v\n", err)
		return
	}

	var testCases map[string]struct {
		bootstrapDataAPI BootstrapDataAPI
		peers            []struct {
			peer *Peer
			cmds []string
		}
		expectedNumberOfPeers int
	}
	testCases = map[string]struct {
		bootstrapDataAPI BootstrapDataAPI
		peers            []struct {
			peer *Peer
			cmds []string
		}
		expectedNumberOfPeers int
	}{
		"Empty Reset": {
			bootstrapDataAPI: BootstrapDataAPI{},
			peers: []struct {
				peer *Peer
				cmds []string
			}{
				{
					peer: &Peer{
						Name:     "alice",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"reset"},
				},
			},
			expectedNumberOfPeers: 0,
		},
		"Reset after one": {
			bootstrapDataAPI: BootstrapDataAPI{},
			peers: []struct {
				peer *Peer
				cmds []string
			}{
				{
					peer: &Peer{
						Name:     "alice",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"join", "reset"},
				},
			},
			expectedNumberOfPeers: 0,
		},
		"Reset after three": {
			bootstrapDataAPI: BootstrapDataAPI{},
			peers: []struct {
				peer *Peer
				cmds []string
			}{
				{
					peer: &Peer{
						Name:     "alice",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"join"},
				},
				{
					peer: &Peer{
						Name:     "bob",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"join"},
				},
				{
					peer: &Peer{
						Name:     "charly",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"join", "reset"},
				},
			},
			expectedNumberOfPeers: 0,
		},
	}

	var bootstrapData *BootstrapData

	for n, tc := range testCases {

		dummy.Reset()

		t.Run(n, func(t *testing.T) {
			for _, peer := range tc.peers {
				bootstrapDataAPI, err := Create(peer.peer)
				if err != nil {
					t.Errorf("cannot create API: %v\n", err)
					return
				}
				for _, cmd := range peer.cmds {
					cmdArgs := strings.Split(cmd, " ")
					//fmt.Printf("%q.%q: %v\n", peer.peer.Name, cmd, cmdArgs)
					if len(cmdArgs) > 1 {
						bootstrapData = executeApiFunc(bootstrapDataAPI, cmdArgs[0],
							strings.Join(cmdArgs[1:], " "))
					} else {
						bootstrapData = executeApiFunc(bootstrapDataAPI, cmdArgs[0])
					}
				}
			}
			if len(bootstrapData.Peers) != tc.expectedNumberOfPeers {
				t.Errorf("number of listed peers %d and not %d",
					len(bootstrapData.Peers), tc.expectedNumberOfPeers)
				return
			}
			if bootstrapData.Config.NumPeers != tc.expectedNumberOfPeers {
				t.Errorf("NumPeers in config %d and not %d",
					bootstrapData.Config.NumPeers, tc.expectedNumberOfPeers)
				return
			}
		})
	}
}

func TestLocalhostJoin(t *testing.T) {
	serviceUrl := os.Getenv("BOOTSTRAP_DATA_SERVER")
	if serviceUrl != "http://localhost:8080" {
		t.Errorf("BOOTSTRAP_DATA_SERVER environment variable expected %q set %q",
			"http://localhost:8080", serviceUrl)
	}

	dummy, err := Create(&Peer{})
	if err != nil {
		t.Errorf("cannot create dummy API: %v\n", err)
		return
	}
	dummy.Reset()

	var testCases map[string]struct {
		bootstrapDataAPI BootstrapDataAPI
		peers            []struct {
			peer *Peer
			cmds []string
		}
		expectedNumberOfPeers int
		expectedName          string
		expectedIp            string
		expectedPort          string
		expectedProtocol      string
	}
	testCases = map[string]struct {
		bootstrapDataAPI BootstrapDataAPI
		peers            []struct {
			peer *Peer
			cmds []string
		}
		expectedNumberOfPeers int
		expectedName          string
		expectedIp            string
		expectedPort          string
		expectedProtocol      string
	}{
		"Join empty": {
			bootstrapDataAPI: BootstrapDataAPI{},
			peers: []struct {
				peer *Peer
				cmds []string
			}{
				{
					peer: &Peer{
						Name:     "alice",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"join"},
				},
			},
			expectedNumberOfPeers: 1,
			expectedName:          "alice",
			expectedIp:            "127.0.0.1",
			expectedPort:          "22365",
			expectedProtocol:      "tcp",
		},
		"Join more than MaxPeers 2": {
			bootstrapDataAPI: BootstrapDataAPI{},
			peers: []struct {
				peer *Peer
				cmds []string
			}{
				{
					peer: &Peer{
						Name:     "alice",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"join"},
				},
				{
					peer: &Peer{
						Name:     "bob",
						Ip:       "127.0.0.2",
						Port:     "22366",
						Protocol: "tcp",
					},
					cmds: []string{"join"},
				},
				{
					peer: &Peer{
						Name:     "charly",
						Ip:       "127.0.0.3",
						Port:     "22367",
						Protocol: "tcp",
					},
					cmds: []string{"join"},
				},
			},
			expectedNumberOfPeers: 2,
			expectedName:          "bob",
			expectedIp:            "127.0.0.2",
			expectedPort:          "22366",
			expectedProtocol:      "tcp",
		},
		"Join less than MaxPeers 2": {
			bootstrapDataAPI: BootstrapDataAPI{},
			peers: []struct {
				peer *Peer
				cmds []string
			}{
				{
					peer: &Peer{
						Name:     "alice",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"join"},
				},
				{
					peer: &Peer{
						Name:     "bob",
						Ip:       "127.0.0.2",
						Port:     "22366",
						Protocol: "tcp",
					},
					cmds: []string{"join"},
				},
			},
			expectedNumberOfPeers: 2,
			expectedName:          "bob",
			expectedIp:            "127.0.0.2",
			expectedPort:          "22366",
			expectedProtocol:      "tcp",
		},
	}

	var bootstrapData *BootstrapData

	for n, tc := range testCases {

		dummy.Reset()

		t.Run(n, func(t *testing.T) {
			for _, peer := range tc.peers {
				bootstrapDataAPI, err := Create(peer.peer)
				if err != nil {
					t.Errorf("cannot create API: %v\n", err)
					return
				}
				for _, cmd := range peer.cmds {
					cmdArgs := strings.Split(cmd, " ")
					//fmt.Printf("%q.%q: %v\n", peer.peer.Name, cmd, cmdArgs)
					if len(cmdArgs) > 1 {
						bootstrapData = executeApiFunc(bootstrapDataAPI, cmdArgs[0],
							strings.Join(cmdArgs[1:], " "))
					} else {
						bootstrapData = executeApiFunc(bootstrapDataAPI, cmdArgs[0])
					}
				}
			}
			if len(bootstrapData.Peers) != tc.expectedNumberOfPeers {
				t.Errorf("number of listed peers %d and not %d",
					len(bootstrapData.Peers), tc.expectedNumberOfPeers)
				return
			}
			if bootstrapData.Config.NumPeers != tc.expectedNumberOfPeers {
				t.Errorf("NumPeers in config %d and not %d",
					bootstrapData.Config.NumPeers, tc.expectedNumberOfPeers)
				return
			}
			found := false
			for key, peer := range bootstrapData.Peers {
				if peer.Name == tc.expectedName {
					found = true
					if bootstrapData.Peers[key].Ip != tc.expectedIp {
						t.Errorf("Ip %q found and not expected %q",
							bootstrapData.Peers[key].Ip, tc.expectedIp)
						return
					}
					if bootstrapData.Peers[key].Port != tc.expectedPort {
						t.Errorf("Port %q found and not expected %q",
							bootstrapData.Peers[key].Port, tc.expectedPort)
						return
					}
					if bootstrapData.Peers[key].Protocol != tc.expectedProtocol {
						t.Errorf("Protocol %q found and not expected %q",
							bootstrapData.Peers[key].Protocol, tc.expectedProtocol)
						return
					}
				}
			}
			if !found {
				t.Errorf("expected peer %q not found", tc.expectedName)
				return
			}
		})
	}
}

func TestLocalhostLeave(t *testing.T) {
	serviceUrl := os.Getenv("BOOTSTRAP_DATA_SERVER")
	if serviceUrl != "http://localhost:8080" {
		t.Errorf("BOOTSTRAP_DATA_SERVER environment variable expected %q set %q",
			"http://localhost:8080", serviceUrl)
	}

	dummy, err := Create(&Peer{})
	if err != nil {
		t.Errorf("cannot create dummy API: %v\n", err)
		return
	}
	dummy.Reset()

	var testCases map[string]struct {
		bootstrapDataAPI BootstrapDataAPI
		peers            []struct {
			peer *Peer
			cmds []string
		}
		expectedNumberOfPeers int
		removedName           string
	}
	testCases = map[string]struct {
		bootstrapDataAPI BootstrapDataAPI
		peers            []struct {
			peer *Peer
			cmds []string
		}
		expectedNumberOfPeers int
		removedName           string
	}{
		"Leave single peer": {
			bootstrapDataAPI: BootstrapDataAPI{},
			peers: []struct {
				peer *Peer
				cmds []string
			}{
				{
					peer: &Peer{
						Name:     "alice",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"join alice", "leave alice"},
				},
			},
			expectedNumberOfPeers: 0,
			removedName:           "alice",
		},
		"Leave list peer": {
			bootstrapDataAPI: BootstrapDataAPI{},
			peers: []struct {
				peer *Peer
				cmds []string
			}{
				{
					peer: &Peer{
						Name:     "alice",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"join"},
				},
				{
					peer: &Peer{
						Name:     "bob",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"join bob", "leave bob"},
				},
			},
			expectedNumberOfPeers: 1,
			removedName:           "bob",
		},
		"Leave not existing gt 2": {
			bootstrapDataAPI: BootstrapDataAPI{},
			peers: []struct {
				peer *Peer
				cmds []string
			}{
				{
					peer: &Peer{
						Name:     "alice",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"join"},
				},
				{
					peer: &Peer{
						Name:     "bob",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"join"},
				},
				{
					peer: &Peer{
						Name:     "charly",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"join charly", "leave charly"},
				},
			},
			expectedNumberOfPeers: 2,
			removedName:           "charly",
		},
	}

	var bootstrapData *BootstrapData

	for n, tc := range testCases {

		dummy.Reset()

		t.Run(n, func(t *testing.T) {
			for _, peer := range tc.peers {
				bootstrapDataAPI, err := Create(peer.peer)
				if err != nil {
					t.Errorf("cannot create API: %v\n", err)
					return
				}
				for _, cmd := range peer.cmds {
					cmdArgs := strings.Split(cmd, " ")
					//fmt.Printf("%q.%q: %v\n", peer.peer.Name, cmd, cmdArgs)
					if len(cmdArgs) > 1 {
						bootstrapData = executeApiFunc(bootstrapDataAPI, cmdArgs[0],
							strings.Join(cmdArgs[1:], " "))
					} else {
						bootstrapData = executeApiFunc(bootstrapDataAPI, cmdArgs[0])
					}
				}
			}
			if len(bootstrapData.Peers) != tc.expectedNumberOfPeers {
				t.Errorf("number of listed peers %d and not %d",
					len(bootstrapData.Peers), tc.expectedNumberOfPeers)
				return
			}
			if bootstrapData.Config.NumPeers != tc.expectedNumberOfPeers {
				t.Errorf("NumPeers in config %d and not %d",
					bootstrapData.Config.NumPeers, tc.expectedNumberOfPeers)
				return
			}
			found := false
			for _, peer := range bootstrapData.Peers {
				if peer.Name == tc.removedName {
					found = true
				}
			}
			if found {
				t.Errorf("removed peer %q found", tc.removedName)
				return
			}
		})
	}
}
func TestLocalhostRefill(t *testing.T) {
	serviceUrl := os.Getenv("BOOTSTRAP_DATA_SERVER")
	if serviceUrl != "http://localhost:8080" {
		t.Errorf("BOOTSTRAP_DATA_SERVER environment variable expected %q set %q",
			"http://localhost:8080", serviceUrl)
	}

	dummy, err := Create(&Peer{})
	if err != nil {
		t.Errorf("cannot create dummy API: %v\n", err)
		return
	}
	dummy.Reset()

	dummy.BootstrapData.Config.MaxPeers = 2
	dummy.BootstrapData.Config.MinPeers = 1

	_, err = dummy.UpdateConfig()
	if err != nil {
		t.Errorf("cannot set dummy API config to default: %v\n", err)
		return
	}

	var testCases map[string]struct {
		bootstrapDataAPI BootstrapDataAPI
		peers            []struct {
			peer *Peer
			cmds []string
		}
		expectedNumberOfPeers int
		expectedName          string
		expectedIp            string
		expectedPort          string
		expectedProtocol      string
	}
	testCases = map[string]struct {
		bootstrapDataAPI BootstrapDataAPI
		peers            []struct {
			peer *Peer
			cmds []string
		}
		expectedNumberOfPeers int
		expectedName          string
		expectedIp            string
		expectedPort          string
		expectedProtocol      string
	}{
		"Refill empty": {
			bootstrapDataAPI: BootstrapDataAPI{},
			peers: []struct {
				peer *Peer
				cmds []string
			}{
				{
					peer: &Peer{
						Name:     "alice",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"refill"},
				},
			},
			expectedNumberOfPeers: 1,
			expectedName:          "alice",
			expectedIp:            "127.0.0.1",
			expectedPort:          "22365",
			expectedProtocol:      "tcp",
		},
		"Refill more than MaxPeers 2 and older": {
			bootstrapDataAPI: BootstrapDataAPI{},
			peers: []struct {
				peer *Peer
				cmds []string
			}{
				{
					peer: &Peer{
						Name:      "alice",
						Ip:        "127.0.0.1",
						Port:      "22365",
						Protocol:  "tcp",
						Timestamp: "10",
					},
					cmds: []string{"refill"},
				},
				{
					peer: &Peer{
						Name:      "bob",
						Ip:        "127.0.0.2",
						Port:      "22366",
						Protocol:  "tcp",
						Timestamp: "20",
					},
					cmds: []string{"refill"},
				},
				{
					peer: &Peer{
						Name:      "charly",
						Ip:        "127.0.0.3",
						Port:      "22367",
						Protocol:  "tcp",
						Timestamp: "1",
					},
					cmds: []string{"refill"},
				},
			},
			expectedNumberOfPeers: 2,
			expectedName:          "charly",
			expectedIp:            "127.0.0.3",
			expectedPort:          "22367",
			expectedProtocol:      "tcp",
		},
		"Refill more than MaxPeers 2 and younger": {
			bootstrapDataAPI: BootstrapDataAPI{},
			peers: []struct {
				peer *Peer
				cmds []string
			}{
				{
					peer: &Peer{
						Name:      "alice",
						Ip:        "127.0.0.1",
						Port:      "22365",
						Protocol:  "tcp",
						Timestamp: "10",
					},
					cmds: []string{"join"},
				},
				{
					peer: &Peer{
						Name:      "bob",
						Ip:        "127.0.0.2",
						Port:      "22366",
						Protocol:  "tcp",
						Timestamp: "20",
					},
					cmds: []string{"join"},
				},
				{
					peer: &Peer{
						Name:      "charly",
						Ip:        "127.0.0.3",
						Port:      "22367",
						Protocol:  "tcp",
						Timestamp: "30",
					},
					cmds: []string{"refill"},
				},
			},
			expectedNumberOfPeers: 2,
			expectedName:          "bob",
			expectedIp:            "127.0.0.2",
			expectedPort:          "22366",
			expectedProtocol:      "tcp",
		},
		"Refill less than MaxPeers 2": {
			bootstrapDataAPI: BootstrapDataAPI{},
			peers: []struct {
				peer *Peer
				cmds []string
			}{
				{
					peer: &Peer{
						Name:     "alice",
						Ip:       "127.0.0.1",
						Port:     "22365",
						Protocol: "tcp",
					},
					cmds: []string{"join"},
				},
				{
					peer: &Peer{
						Name:     "bob",
						Ip:       "127.0.0.2",
						Port:     "22366",
						Protocol: "tcp",
					},
					cmds: []string{"refill"},
				},
			},
			expectedNumberOfPeers: 2,
			expectedName:          "bob",
			expectedIp:            "127.0.0.2",
			expectedPort:          "22366",
			expectedProtocol:      "tcp",
		},
	}

	var bootstrapData *BootstrapData

	for n, tc := range testCases {

		dummy.Reset()

		t.Run(n, func(t *testing.T) {
			for _, peer := range tc.peers {
				timestamp := peer.peer.Timestamp
				bootstrapDataAPI, err := Create(peer.peer)
				if err != nil {
					t.Errorf("cannot create API: %v\n", err)
					return
				}

				bootstrapDataAPI.Self.Timestamp = timestamp

				for _, cmd := range peer.cmds {
					cmdArgs := strings.Split(cmd, " ")
					//fmt.Printf("%q.%q: %v\n", peer.peer.Name, cmd, cmdArgs)
					if len(cmdArgs) > 1 {
						bootstrapData = executeApiFunc(bootstrapDataAPI, cmdArgs[0],
							strings.Join(cmdArgs[1:], " "))
					} else {
						bootstrapData = executeApiFunc(bootstrapDataAPI, cmdArgs[0])
					}
				}
			}
			if len(bootstrapData.Peers) != tc.expectedNumberOfPeers {
				t.Errorf("number of listed peers %d and not %d",
					len(bootstrapData.Peers), tc.expectedNumberOfPeers)
				return
			}
			if bootstrapData.Config.NumPeers != tc.expectedNumberOfPeers {
				t.Errorf("NumPeers in config %d and not %d",
					bootstrapData.Config.NumPeers, tc.expectedNumberOfPeers)
				return
			}
			found := false
			for key, peer := range bootstrapData.Peers {
				if peer.Name == tc.expectedName {
					found = true
					if bootstrapData.Peers[key].Ip != tc.expectedIp {
						t.Errorf("Ip %q found and not expected %q",
							bootstrapData.Peers[key].Ip, tc.expectedIp)
						return
					}
					if bootstrapData.Peers[key].Port != tc.expectedPort {
						t.Errorf("Port %q found and not expected %q",
							bootstrapData.Peers[key].Port, tc.expectedPort)
						return
					}
					if bootstrapData.Peers[key].Protocol != tc.expectedProtocol {
						t.Errorf("Protocol %q found and not expected %q",
							bootstrapData.Peers[key].Protocol, tc.expectedProtocol)
						return
					}
				}
			}
			if !found {
				t.Errorf("expected peer %q not found", tc.expectedName)
				return
			}
		})
	}
}

func TestLocalhostConfig(t *testing.T) {

	serviceUrl := os.Getenv("BOOTSTRAP_DATA_SERVER")
	if serviceUrl != "http://localhost:8080" {
		t.Errorf("BOOTSTRAP_DATA_SERVER environment variable expected %q set %q",
			"http://localhost:8080", serviceUrl)
	}

	dummy, err := Create(&Peer{})
	if err != nil {
		t.Errorf("cannot create dummy API: %v\n", err)
		return
	}

	var testCases map[string]struct {
		bootstrapDataAPI BootstrapDataAPI
		maxPeers         int
		minPeers         int
		expectedError    bool
	}

	testCases = map[string]struct {
		bootstrapDataAPI BootstrapDataAPI
		maxPeers         int
		minPeers         int
		expectedError    bool
	}{
		"Reset to zero": {
			bootstrapDataAPI: BootstrapDataAPI{},
			maxPeers:         0,
			minPeers:         0,
			expectedError:    false,
		},
		"Set to default": {
			bootstrapDataAPI: BootstrapDataAPI{},
			maxPeers:         2,
			minPeers:         1,
			expectedError:    false,
		},
		"Change both": {
			bootstrapDataAPI: BootstrapDataAPI{},
			maxPeers:         8,
			minPeers:         5,
			expectedError:    false,
		},
		"max lt min": {
			bootstrapDataAPI: BootstrapDataAPI{},
			maxPeers:         1,
			minPeers:         2,
			expectedError:    true,
		},
	}

	var bootstrapData *BootstrapData

	for n, tc := range testCases {

		t.Run(n, func(t *testing.T) {

			dummy.BootstrapData.Config.MaxPeers = tc.maxPeers
			dummy.BootstrapData.Config.MinPeers = tc.minPeers

			bootstrapData, err = dummy.UpdateConfig()
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error %v, received %v", tc.expectedError, err)
			}

			if !tc.expectedError {
				if bootstrapData.Config.MaxPeers != tc.maxPeers ||
					bootstrapData.Config.MinPeers != tc.minPeers {
					t.Errorf("expected MaxPeers/MinPeers %d/%d, received  %d/%d",
						tc.maxPeers, tc.minPeers,
						bootstrapData.Config.MaxPeers, bootstrapData.Config.MinPeers)
				}
			}
		})
	}
}
