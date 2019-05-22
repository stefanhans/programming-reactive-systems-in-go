package serviceprovider

import (
	"fmt"
	"log"
	"testing"

	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/service-provider-chat/serviceprovider/serviceproviders"
)

func TestInsertService(t *testing.T) {

	testcases := map[string]struct {

		// Slice of services to add
		testservices []struct {
			newServiceName   string
			newServiceStatus serviceproviders.Service_Status
		}

		// Expected final results
		expectedLen int
		expectError bool
	}{
		"empty list": {
			testservices: []struct {
				newServiceName   string
				newServiceStatus serviceproviders.Service_Status
			}{
				{
					// #1 Service
					newServiceName:   "srv1",
					newServiceStatus: serviceproviders.Service_CANDIDATE,
				},
			},
			expectedLen: 1,
			expectError: false,
		},
		"name exists": {
			testservices: []struct {
				newServiceName   string
				newServiceStatus serviceproviders.Service_Status
			}{
				{
					// #1 Service
					newServiceName:   "newservice",
					newServiceStatus: serviceproviders.Service_CANDIDATE,
				},
				{
					// #2 Service
					newServiceName:   "newservice",
					newServiceStatus: serviceproviders.Service_CANDIDATE,
				},
			},
			expectedLen: 1,
			expectError: true,
		},
		"list with another service": {
			testservices: []struct {
				newServiceName   string
				newServiceStatus serviceproviders.Service_Status
			}{
				{
					// #1 Service
					newServiceName:   "srv1",
					newServiceStatus: serviceproviders.Service_SERVICE,
				},
				{
					// #2 Service
					newServiceName:   "srv2",
					newServiceStatus: serviceproviders.Service_CANDIDATE,
				},
				{
					// #3 Service
					newServiceName:   "srv3",
					newServiceStatus: serviceproviders.Service_CANDIDATE,
				},
			},
			expectedLen: 3,
			expectError: false,
		},
	}
	for n, tc := range testcases {

		serviceProvider, err := NewServiceProvider(PROVIDER,
			"", "", "",
			"", "", "")
		if err != nil {
			log.Fatalf("could not create new serviceProvider: %v", err)
		}

		var isError bool
		for _, ts := range tc.testservices {
			service := &serviceproviders.Service{
				Name:   ts.newServiceName,
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

			isError = serviceProvider.insertService(message) != nil

		}

		// len(services)
		t.Run(fmt.Sprintf("%s/len(services)", n), func(t *testing.T) {
			if len(serviceProvider.serviceProvider.Services) != tc.expectedLen {
				t.Errorf("Length of service list is not %d: %d", tc.expectedLen, len(serviceProvider.serviceProvider.Services))
			}
		})

		// (err != nil)
		t.Run(fmt.Sprintf("%s/(err != nil)", n), func(t *testing.T) {
			if isError != tc.expectError {
				t.Errorf("Error is not expected: %v", err)
			}
		})
	}
}

func TestServiceElection(t *testing.T) {

	testcases := map[string]struct {

		// Slice of services to add
		testservices []struct {
			newServiceStatus serviceproviders.Service_Status
		}

		// Expected final results
		expectedLen           int
		expectError           bool
		expectGetServiceError bool
		expectedServices      int
	}{
		"empty list": {
			testservices: []struct {
				newServiceStatus serviceproviders.Service_Status
			}{},
			expectedLen:           0,
			expectError:           true,
			expectGetServiceError: true,
			expectedServices:      -1,
		},
		"one service": {
			testservices: []struct {
				newServiceStatus serviceproviders.Service_Status
			}{
				{
					// #1 Service
					newServiceStatus: serviceproviders.Service_SERVICE,
				},
			},
			expectedLen:           1,
			expectError:           false,
			expectGetServiceError: false,
			expectedServices:      0,
		},
		"one candidate": {
			testservices: []struct {
				newServiceStatus serviceproviders.Service_Status
			}{
				{
					// #1 Service
					newServiceStatus: serviceproviders.Service_CANDIDATE,
				},
			},
			expectedLen:           1,
			expectError:           false,
			expectGetServiceError: false,
			expectedServices:      0,
		},
		"no candidate": {
			testservices: []struct {
				newServiceStatus serviceproviders.Service_Status
			}{
				{
					// #1 Service
					newServiceStatus: serviceproviders.Service_NOTFOUND,
				},
			},
			expectedLen:           1,
			expectError:           true,
			expectGetServiceError: true,
			expectedServices:      0,
		},
		"second service of two": {
			testservices: []struct {
				newServiceStatus serviceproviders.Service_Status
			}{
				{
					// #1 Service
					newServiceStatus: serviceproviders.Service_CANDIDATE,
				},
				{
					// #2 Service
					newServiceStatus: serviceproviders.Service_SERVICE,
				},
			},
			expectedLen:           2,
			expectError:           false,
			expectGetServiceError: false,
			expectedServices:      1,
		},
		"two candidates": {
			testservices: []struct {
				newServiceStatus serviceproviders.Service_Status
			}{
				{
					// #1 Service
					newServiceStatus: serviceproviders.Service_CANDIDATE,
				},
				{
					// #2 Service
					newServiceStatus: serviceproviders.Service_CANDIDATE,
				},
			},
			expectedLen:           2,
			expectError:           false,
			expectGetServiceError: false,
			expectedServices:      0,
		},
	}
	for n, tc := range testcases {

		serviceProvider, err := NewServiceProvider(PROVIDER,
			"", "", "",
			"", "", "")
		if err != nil {
			log.Fatalf("could not create new serviceProvider: %v", err)
		}

		var isError bool
		if len(tc.testservices) != 0 {

			for i, ts := range tc.testservices {
				service := &serviceproviders.Service{
					Name:   fmt.Sprintf("#%d", i),
					Status: ts.newServiceStatus}
				//fmt.Printf("ServiceProvider before: %v", serviceProvider)
				//_ = service

				var services []*serviceproviders.Service
				services = append(services, service)

				message := &serviceproviders.Message{
					MsgType: serviceproviders.Message_SERVICE_REQUEST,
					ServiceProvider: &serviceproviders.ServiceProvider{
						Provider: &serviceproviders.Service{},
						Services: services,
					},
				}
				//fmt.Printf("message: %v", message)

				//fmt.Printf("ServiceProvider %v: %v\n\n", err, serviceProvider)
				serviceProvider.insertService(message)

			}
		}

		// Do the test, i.e. elect the service
		_, err = serviceProvider.serviceElection()
		isError = err != nil

		// len(services)
		t.Run(fmt.Sprintf("%s/len(services)", n), func(t *testing.T) {
			if len(serviceProvider.serviceProvider.Services) != tc.expectedLen {
				t.Errorf("Length of service list is not %d: %d", tc.expectedLen, len(serviceProvider.serviceProvider.Services))
			}
		})

		// (err != nil)
		t.Run(fmt.Sprintf("%s/(err != nil)", n), func(t *testing.T) {
			if isError != tc.expectError {
				t.Errorf("Error is not expected: %v", err)
			}
		})

		// Get the service
		returnedService, getServiceError := serviceProvider.GetService()

		// get service returns error
		t.Run(fmt.Sprintf("%s/get service returns error", n), func(t *testing.T) {
			if (getServiceError != nil) != tc.expectGetServiceError {
				t.Errorf("Error is not expected: %v", err)
			}
		})

		// expected service
		t.Run(fmt.Sprintf("%s/expected service", n), func(t *testing.T) {
			if getServiceError == nil {
				if returnedService != serviceProvider.serviceProvider.Services[tc.expectedServices] ||
					returnedService.Status != serviceproviders.Service_SERVICE {
					t.Errorf("Not the expected service: %v", returnedService)
				}
			}
		})
	}
}
