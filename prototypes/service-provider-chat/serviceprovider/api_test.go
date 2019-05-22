package serviceprovider

import (
	"fmt"
	"log"
	"net"
	"strings"
	"testing"

	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/service-provider-chat/serviceprovider/serviceproviders"
	"time"
)

func TestNewServiceProvider(t *testing.T) {

	// TestCases
	var testcases = []struct {
		serviceProviderType ServiceProviderType
		candidateName       string
		candidateIp         string
		candidatePort       string
		providerName        string
		providerIp          string
		providerPort        string
	}{
		{serviceProviderType: PROVIDER,
			candidateName: "candidate",
			candidateIp:   "127.0.0.1",
			candidatePort: "12345",
			providerName:  "provider",
			providerIp:    "127.0.0.2",
			providerPort:  "22365"},
		{serviceProviderType: SERVICE,
			candidateName: "candidate",
			candidateIp:   "127.0.0.1",
			candidatePort: "12345",
			providerName:  "provider",
			providerIp:    "127.0.0.2",
			providerPort:  "22365"},
	}
	for _, tc := range testcases {
		serviceProvider, err := NewServiceProvider(
			tc.serviceProviderType,
			tc.candidateName,
			tc.candidateIp,
			tc.candidatePort,
			tc.providerName,
			tc.providerIp,
			tc.providerPort)
		if err != nil {
			log.Fatalf("could not create new serviceProvider: %v", err)
		}

		// serviceProviderType
		t.Run("serviceProviderType", func(t *testing.T) {
			if serviceProvider.serviceProviderType != tc.serviceProviderType {
				t.Errorf("Unexpected serviceProviderType: %v", serviceProvider.serviceProviderType)
			}
		})

		// version
		t.Run("version", func(t *testing.T) {
			if serviceProvider.version != 0 {
				t.Errorf("Unexpected version: %v", serviceProvider.version)
			}
		})

		// candidate.Name
		t.Run("candidate.Name", func(t *testing.T) {
			if serviceProvider.candidate.Name != tc.candidateName {
				t.Errorf("Unexpected candidate.Name: %v", serviceProvider.candidate.Name)
			}
		})

		// candidate.Ip
		t.Run("candidate.Ip", func(t *testing.T) {
			if serviceProvider.candidate.Ip != tc.candidateIp {
				t.Errorf("Unexpected candidate.Ip: %v", serviceProvider.candidate.Ip)
			}
		})

		// candidate.Port
		t.Run("candidate.Port", func(t *testing.T) {
			if serviceProvider.candidate.Port != tc.candidatePort {
				t.Errorf("Unexpected candidate.Port: %v", serviceProvider.candidate.Port)
			}
		})

		// provider.Name
		t.Run("provider.Name", func(t *testing.T) {
			if serviceProvider.provider.Name != tc.providerName {
				t.Errorf("Unexpected provider.Name: %v", serviceProvider.provider.Name)
			}
		})

		// provider.Ip
		t.Run("provider.Ip", func(t *testing.T) {
			if serviceProvider.provider.Ip != tc.providerIp {
				t.Errorf("Unexpected provider.Ip: %v", serviceProvider.provider.Ip)
			}
		})

		// provider.Port
		t.Run("provider.Port", func(t *testing.T) {
			if serviceProvider.provider.Port != tc.providerPort {
				t.Errorf("Unexpected provider.Port: %v", serviceProvider.provider.Port)
			}
		})

		// ServiceProvider.Provider
		t.Run("ServiceProvider.Provider", func(t *testing.T) {
			if serviceProvider.serviceProvider.Provider != serviceProvider.provider {
				t.Errorf("Unexpected ServiceProvider.Provider: %v", serviceProvider.serviceProvider.Provider)
			}
		})

		if serviceProvider.serviceProviderType == SERVICE {
			// ServiceProvider.Services[0]
			t.Run("ServiceProvider.Services[0]", func(t *testing.T) {
				if serviceProvider.serviceProvider.Services[0].Name != serviceProvider.candidate.Name {
					t.Errorf("Unexpected ServiceProvider.Services[0].Name: %v", serviceProvider.serviceProvider.Services[0].Name)
				}
				if serviceProvider.serviceProvider.Services[0].Ip != serviceProvider.candidate.Ip {
					t.Errorf("Unexpected ServiceProvider.Services[0].Ip: %v", serviceProvider.serviceProvider.Services[0].Ip)
				}
				if serviceProvider.serviceProvider.Services[0].Port != serviceProvider.candidate.Port {
					t.Errorf("Unexpected ServiceProvider.Services[0].Port: %v", serviceProvider.serviceProvider.Services[0].Port)
				}
				if serviceProvider.serviceProvider.Services[0].Status != serviceProvider.candidate.Status {
					t.Errorf("Unexpected ServiceProvider.Services[0].Status: %v", serviceProvider.serviceProvider.Services[0].Status)
				}
			})
		}

		// message.MsgType
		t.Run("message.MsgType", func(t *testing.T) {
			if serviceProvider.message.MsgType != serviceproviders.Message_SERVICE_REQUEST {
				t.Errorf("Unexpected message.MsgType: %v", serviceProvider.message.MsgType)
			}
		})

		// message.ServiceProvider
		t.Run("message.ServiceProvider", func(t *testing.T) {
			if serviceProvider.message.ServiceProvider != serviceProvider.serviceProvider {
				t.Errorf("Unexpected message.ServiceProvider: %v", serviceProvider.message.ServiceProvider)
			}
		})

		// message.Sender
		t.Run("message.Sender", func(t *testing.T) {
			if serviceProvider.message.Sender != serviceProvider.candidate {
				t.Errorf("Unexpected message.Sender: %v", serviceProvider.message.Sender)
			}
		})
	}
}

// GetReplyChannel(serviceproviders.Message_SERVICE_REPLY) returns (chan *serviceproviders.Message)
func TestGetReplyChannel(t *testing.T) {

	serviceProvider, err := NewServiceProvider(PROVIDER,
		"", "", "",
		"", "", "")
	if err != nil {
		log.Fatalf("could not create new serviceProvider: %v", err)
	}

	t.Run("Message_SERVICE_REPLY", func(t *testing.T) {
		retType := serviceProvider.GetReplyChannel(serviceproviders.Message_SERVICE_REPLY)

		if fmt.Sprintf("%T", retType) != "chan *serviceproviders.Message" {
			t.Errorf("Unexpected channel type: %T", retType)
		}
	})
}

func TestString(t *testing.T) {

	// TestCases
	var testcases = []struct {
		serviceProviderType    ServiceProviderType
		expServiceProviderType string
		candidateName          string
		candidateIp            string
		candidatePort          string
		providerName           string
		providerIp             string
		providerPort           string
	}{
		{
			serviceProviderType:    PROVIDER,
			expServiceProviderType: "PROVIDER",
			candidateName:          "candidate",
			candidateIp:            "127.0.0.1",
			candidatePort:          "12345",
			providerName:           "provider",
			providerIp:             "127.0.0.2",
			providerPort:           "22365"},
		{
			serviceProviderType:    SERVICE,
			expServiceProviderType: "SERVICE",
			candidateName:          "candidate",
			candidateIp:            "127.0.0.1",
			candidatePort:          "12345",
			providerName:           "provider",
			providerIp:             "127.0.0.2",
			providerPort:           "22365"},
	}
	for _, tc := range testcases {
		serviceProvider, err := NewServiceProvider(
			tc.serviceProviderType,
			tc.candidateName,
			tc.candidateIp,
			tc.candidatePort,
			tc.providerName,
			tc.providerIp,
			tc.providerPort)
		if err != nil {
			log.Fatalf("could not create new serviceProvider: %v", err)
		}

		// String
		t.Run("serviceProviderType", func(t *testing.T) {
			if !strings.Contains(serviceProvider.String(), tc.expServiceProviderType) {
				t.Errorf("no %q in String(): %v", tc.expServiceProviderType, serviceProvider.String())
			}
		})
		t.Run("candidate", func(t *testing.T) {
			if !strings.Contains(serviceProvider.String(), tc.candidateName) {
				t.Errorf("no %q in String(): %v", tc.candidateName, serviceProvider.String())
			}
			if !strings.Contains(serviceProvider.String(), tc.candidateIp) {
				t.Errorf("no %q in String(): %v", tc.candidateIp, serviceProvider.String())
			}
			if !strings.Contains(serviceProvider.String(), tc.candidatePort) {
				t.Errorf("no %q in String(): %v", tc.candidatePort, serviceProvider.String())
			}
			if serviceProvider.serviceProviderType == SERVICE &&
				!strings.Contains(serviceProvider.String(), "status:CANDIDATE") {
				t.Errorf("no \"status:CANDIDATE\" in String(): %v", serviceProvider.String())
			}
		})
		t.Run("provider", func(t *testing.T) {
			if !strings.Contains(serviceProvider.String(), tc.providerName) {
				t.Errorf("no %q in String(): %v", tc.providerName, serviceProvider.String())
			}
			if !strings.Contains(serviceProvider.String(), tc.providerIp) {
				t.Errorf("no %q in String(): %v", tc.providerIp, serviceProvider.String())
			}
			if !strings.Contains(serviceProvider.String(), tc.providerPort) {
				t.Errorf("no %q in String(): %v", tc.providerPort, serviceProvider.String())
			}
			if !strings.Contains(serviceProvider.String(), "status:PROVIDER") {
				t.Errorf("no \"status:PROVIDER\" in String(): %v", serviceProvider.String())
			}
		})
		t.Run("SERVICE_REQUEST", func(t *testing.T) {
			if !strings.Contains(serviceProvider.String(), "SERVICE_REQUEST") {
				t.Errorf("no \"SERVICE_REQUEST\" in String(): %v", serviceProvider.String())
			}
		})
		t.Run("MessageHandlerFunctions", func(t *testing.T) {
			if !strings.Contains(serviceProvider.String(), "MessageHandlerFunctions") {
				t.Errorf("no \"MessageHandlerFunctions\" in String(): %v", serviceProvider.String())
			}
		})
		t.Run("ReplyChannels", func(t *testing.T) {
			if !strings.Contains(serviceProvider.String(), "ReplyChannels") {
				t.Errorf("no \"ReplyChannels\" in String(): %v", serviceProvider.String())
			}
		})
	}
}

func TestMessage(t *testing.T) {

	// TestCases
	var testcases = []struct {
		serviceProviderType ServiceProviderType
		candidateName       string
		candidateIp         string
		candidatePort       string
		providerName        string
		providerIp          string
		providerPort        string
	}{
		{serviceProviderType: PROVIDER,
			candidateName: "candidate",
			candidateIp:   "127.0.0.1",
			candidatePort: "12345",
			providerName:  "provider",
			providerIp:    "127.0.0.2",
			providerPort:  "22365"},
		{serviceProviderType: SERVICE,
			candidateName: "candidate",
			candidateIp:   "127.0.0.1",
			candidatePort: "12345",
			providerName:  "provider",
			providerIp:    "127.0.0.2",
			providerPort:  "22365"},
	}
	for _, tc := range testcases {
		serviceProvider, err := NewServiceProvider(
			tc.serviceProviderType,
			tc.candidateName,
			tc.candidateIp,
			tc.candidatePort,
			tc.providerName,
			tc.providerIp,
			tc.providerPort)
		if err != nil {
			log.Fatalf("could not create new serviceProvider: %v", err)
		}

		// message.MsgType
		t.Run("message.MsgType", func(t *testing.T) {
			if serviceProvider.Message().MsgType != serviceproviders.Message_SERVICE_REQUEST {
				t.Errorf("Unexpected message.MsgType: %v", serviceProvider.message.MsgType)
			}
		})

		// message.ServiceProvider
		t.Run("message.ServiceProvider", func(t *testing.T) {
			if serviceProvider.Message().ServiceProvider != serviceProvider.serviceProvider {
				t.Errorf("Unexpected message.ServiceProvider: %v", serviceProvider.message.ServiceProvider)
			}
		})

		// message.Sender
		t.Run("message.Sender", func(t *testing.T) {
			if serviceProvider.Message().Sender != serviceProvider.candidate {
				t.Errorf("Unexpected message.Sender: %v", serviceProvider.message.Sender)
			}
		})
	}
}

func TestUpdateServices(t *testing.T) {

	// TestCases
	var testcases = map[string]struct {
		// Slice of services to add
		testservices []struct {
			newServiceName   string
			newServiceIp     string
			newServicePort   string
			newServiceStatus serviceproviders.Service_Status
		}
		// Expected final results
		expectedUpdated     bool
		expectedLen         int
		expectedServiceName string
	}{
		"empty list": {
			testservices: []struct {
				newServiceName   string
				newServiceIp     string
				newServicePort   string
				newServiceStatus serviceproviders.Service_Status
			}{
				{
					// #1 Service
					newServiceName:   "srv1",
					newServiceIp:     "127.0.1.1",
					newServicePort:   "1111",
					newServiceStatus: serviceproviders.Service_CANDIDATE,
				},
				{
					newServiceName:   "srv2",
					newServiceIp:     "127.0.1.2",
					newServicePort:   "2222",
					newServiceStatus: serviceproviders.Service_CANDIDATE,
				},
			},
			expectedUpdated:     true,
			expectedLen:         2,
			expectedServiceName: "srv1",
		},
		"service already exist": {
			testservices: []struct {
				newServiceName   string
				newServiceIp     string
				newServicePort   string
				newServiceStatus serviceproviders.Service_Status
			}{
				{
					// #1 Service
					newServiceName:   "srv1",
					newServiceIp:     "127.0.1.1",
					newServicePort:   "1111",
					newServiceStatus: serviceproviders.Service_SERVICE,
				},
				{
					newServiceName:   "srv1",
					newServiceIp:     "127.0.1.2",
					newServicePort:   "2222",
					newServiceStatus: serviceproviders.Service_CANDIDATE,
				},
			},
			expectedUpdated:     true,
			expectedLen:         1,
			expectedServiceName: "srv1",
		},
	}

	for _, tc := range testcases {
		serviceProvider, err := NewServiceProvider(PROVIDER,
			"", "", "",
			"", "", "")
		if err != nil {
			log.Fatalf("could not create new serviceProvider: %v", err)
		}
		for _, ts := range tc.testservices {
			service := &serviceproviders.Service{
				Name:   ts.newServiceName,
				Ip:     ts.newServiceIp,
				Port:   ts.newServicePort,
				Status: ts.newServiceStatus}

			var services []*serviceproviders.Service
			services = append(services, service)

			message := &serviceproviders.Message{
				MsgType: serviceproviders.Message_SERVICE_REQUEST,
				ServiceProvider: &serviceproviders.ServiceProvider{
					Provider: &serviceproviders.Service{},
					Services: services,
				},
			}
			serviceProvider.UpdateServices(message)

		}

		// len(services)
		t.Run("len(services)", func(t *testing.T) {
			if len(serviceProvider.serviceProvider.Services) != tc.expectedLen {
				t.Errorf("Length of service list is not %d: %d", tc.expectedLen, len(serviceProvider.serviceProvider.Services))
			}
		})

		// prepare expectedServiceCount and expectedServiceName
		srvName := ""
		cnt := 0
		for _, srv := range serviceProvider.serviceProvider.Services {
			if srv.Status == serviceproviders.Service_SERVICE {
				srvName = srv.Name
				cnt++
			}
		}

		// expectedServiceCount
		t.Run("number of services", func(t *testing.T) {
			if cnt != 1 {
				t.Errorf("Unexpected number of working services: %v", cnt)
			}
		})

		// expectedServiceName
		t.Run("name of service", func(t *testing.T) {
			if srvName != tc.expectedServiceName {
				t.Errorf("Unexpected name of working service (expected %q): %q", tc.expectedServiceName, srvName)
			}
		})
	}
}

func TestResolveIpNewServiceProvider(t *testing.T) {

	testcases := map[string]struct {
		expectedNewServiceProviderError bool
		candidateIp                     string
		expectedCandidateIp             string
		expectedCandidateIpResolveError bool
		providerIp                      string
		expectedProviderIp              string
		expectedProviderIpResolveError  bool
	}{
		"127.0.0.1": {
			expectedNewServiceProviderError: false,
			candidateIp:                     "127.0.0.1",
			expectedCandidateIp:             "127.0.0.1",
			expectedCandidateIpResolveError: false,
			providerIp:                      "127.0.0.1",
			expectedProviderIp:              "127.0.0.1",
			expectedProviderIpResolveError:  false,
		},
		"localhost": {
			expectedNewServiceProviderError: false,
			candidateIp:                     "localhost",
			expectedCandidateIp:             "127.0.0.1",
			expectedCandidateIpResolveError: false,
			providerIp:                      "localhost",
			expectedProviderIp:              "127.0.0.1",
			expectedProviderIpResolveError:  false,
		},
		"invalid": {
			expectedNewServiceProviderError: true,
			candidateIp:                     "invalid",
			expectedCandidateIp:             "",
			expectedCandidateIpResolveError: true,
			providerIp:                      "invalid",
			expectedProviderIp:              "",
			expectedProviderIpResolveError:  true,
		},
	}
	for n, tc := range testcases {
		serviceProvider, err := NewServiceProvider(
			PROVIDER, "", tc.candidateIp, "",
			"", tc.providerIp, "")
		if err != nil {
			if tc.expectedNewServiceProviderError {
				continue
			}
			log.Fatalf("could not create new serviceProvider: %v", err)
		}

		// String
		t.Run(n, func(t *testing.T) {
			// candidate
			addr, err := net.ResolveIPAddr("ip", serviceProvider.candidate.Ip)
			if (err != nil) != tc.expectedCandidateIpResolveError {
				t.Errorf("candidate: not the expected error: %v", err)
			}
			if tc.expectedCandidateIp != addr.String() {
				t.Errorf("not the expected candidate Ip (%v): %v",
					tc.expectedCandidateIp, addr.String())
			}
			// provider
			addr, err = net.ResolveIPAddr("ip", serviceProvider.provider.Ip)
			if (err != nil) != tc.expectedProviderIpResolveError {
				t.Errorf("provider: not the expected error: %v", err)
			}
			if tc.expectedProviderIp != addr.String() {
				t.Errorf("not the expected provider Ip (%v): %v",
					tc.expectedProviderIp, addr.String())
			}
		})
	}
}

func TestSetterGetter(t *testing.T) {

	serviceProvider, err := NewServiceProvider(
		PROVIDER, "myCandidateName", "", "",
		"myProviderName", "", "")
	if err != nil {
		log.Fatalf("could not create new serviceProvider: %v", err)
	}

	// ServiceProviderType()
	t.Run("ServiceProviderType", func(t *testing.T) {
		if serviceProvider.ServiceProviderType() != PROVIDER {
			t.Errorf("Unexpected result from ServiceProviderType(): %q", serviceProvider.ServiceProviderType())
		}
	})

	// Version()
	t.Run("Version", func(t *testing.T) {
		if serviceProvider.Version() != 0 {
			t.Errorf("Unexpected result from Version(): %q", serviceProvider.Version())
		}
	})

	// CandidateName
	t.Run("CandidateName", func(t *testing.T) {
		if serviceProvider.CandidateName() != "myCandidateName" {
			t.Errorf("Unexpected result from CandidateName(): %v", serviceProvider.CandidateName())
		}
	})

	// CandidateIp
	t.Run("CandidateIp", func(t *testing.T) {
		serviceProvider.SetCandidateIp("127.0.0.1")
		if serviceProvider.CandidateIp() != "127.0.0.1" {
			t.Errorf("Unexpected result from CandidateIp(): %v", serviceProvider.CandidateIp())
		}
	})

	// CandidatePort
	var porttests = []struct {
		port  string
		valid bool
	}{
		{"22365", true},
		{"1024", true},
		{"65535", true},
		{"invalid", false},
		{"-1", false},
		{"0", false},
		{"1023", false},
		{"65536", false},
		{"1234567890", false},
	}
	for _, pt := range porttests {
		t.Run("CandidatePort", func(t *testing.T) {
			serviceProvider.SetCandidatePort(pt.port)
			if (serviceProvider.CandidatePort() == pt.port) != pt.valid {
				t.Errorf("Unexpected validity %v from CandidatePort(%q): %v", pt.valid, pt.port, serviceProvider.CandidatePort())
			}
		})
	}

	// CandidateStatus
	t.Run("CandidateStatus", func(t *testing.T) {
		serviceProvider.SetCandidateStatus(serviceproviders.Service_CANDIDATE)
		if serviceProvider.CandidateStatus() != serviceproviders.Service_CANDIDATE {
			t.Errorf("Unexpected result from CandidateStatus(): %v", serviceProvider.CandidateStatus())
		}
	})

	// ProviderName
	t.Run("ProviderName", func(t *testing.T) {
		if serviceProvider.ProviderName() != "myProviderName" {
			t.Errorf("Unexpected result from ProviderName(): %v", serviceProvider.ProviderName())
		}
	})

	// ProviderIp
	t.Run("ProviderIp", func(t *testing.T) {
		serviceProvider.SetProviderIp("127.0.0.1")
		if serviceProvider.ProviderIp() != "127.0.0.1" {
			t.Errorf("Unexpected result from ProviderIp(): %v", serviceProvider.ProviderIp())
		}
	})

	// ProviderPort
	for _, pt := range porttests {
		t.Run("ProviderPort", func(t *testing.T) {
			serviceProvider.SetProviderPort(pt.port)
			if (serviceProvider.ProviderPort() == pt.port) != pt.valid {
				t.Errorf("Unexpected validity %v from ProviderPort(%q): %v", pt.valid, pt.port, serviceProvider.ProviderPort())
			}
		})
	}
}

func TestGetServices(t *testing.T) {
	candidateName := "testA"
	candidateIp := "127.1.2.3"
	candidatePort := "12345"

	serviceProvider, err := NewServiceProvider(
		SERVICE, candidateName, candidateIp, candidatePort,
		"", "", "")
	if err != nil {
		log.Fatalf("could not create new serviceProvider: %v", err)
	}

	t.Run("Candidate A", func(t *testing.T) {
		services := serviceProvider.serviceProvider.GetServices()

		if services[0].Name != candidateName || services[0].Ip != candidateIp || services[0].Port != candidatePort {
			t.Errorf("Not the expected values of inserted service, i.e. %v != %v", &serviceproviders.Service{
				Name: candidateName,
				Ip:   candidateIp,
				Port: candidatePort,
			}, services[0])
		}
	})
}

func TestGetServiceAddress(t *testing.T) {

	serviceProvider, err := NewServiceProvider(
		PROVIDER, "", "", "",
		"", "", "")
	if err != nil {
		log.Fatalf("could not create new serviceProvider: %v", err)
	}

	t.Run("No active service", func(t *testing.T) {
		addr, err := serviceProvider.GetServiceAddress()
		if addr != "" || err == nil {
			t.Errorf("unexpected result from empty list: %v, %v", addr, err)
		}
	})
}

func TestEndToEnd(t *testing.T) {

	rootServiceProvider, err := NewServiceProvider(
		PROVIDER, "", "", "",
		"rootService", "127.0.0.1", "22365")
	if err != nil {
		log.Fatalf("could not create new service provider: %v", err)
	}
	rootServiceProvider.StartRootService()

	clientServiceProvider, err := NewServiceProvider(
		SERVICE, "clientService", "127.0.0.1", "12345",
		"rootService", "127.0.0.1", "22365")
	if err != nil {
		log.Fatalf("could not create new service provider: %v", err)
	}

	clientServiceProvider.StartClientService()

	clientServiceProvider.RequestServices()

	currentVersion := clientServiceProvider.Version()

	for clientServiceProvider.Version() == currentVersion {
		time.Sleep(time.Millisecond * 10)
	}

	t.Run("client get service", func(t *testing.T) {
		service, err := clientServiceProvider.GetService()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !strings.Contains(service.String(), "clientService") {
			t.Errorf("no %q in String(): %v", "clientService", service)
		}
		if !strings.Contains(service.String(), "127.0.0.1") {
			t.Errorf("no %q in String(): %v", "127.0.0.1", service)
		}
		if !strings.Contains(service.String(), "12345") {
			t.Errorf("no %q in String(): %v", "12345", service)
		}
		if !strings.Contains(service.String(), "SERVICE") {
			t.Errorf("no %q in String(): %v", "SERVICE", service)
		}
	})

	t.Run("client get service address", func(t *testing.T) {
		serviceAddr, err := clientServiceProvider.GetServiceAddress()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if serviceAddr != "127.0.0.1:12345" {
			t.Errorf("not expected %q: %v", "127.0.0.1:12345", serviceAddr)
		}
	})
}
