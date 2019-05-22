package serviceprovider

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/service-provider-chat/serviceprovider/serviceproviders"
)

// insertService insert a service to the list, if name does not already exists
func (serviceProvider *ServiceProvider) insertService(message *serviceproviders.Message) error {

	// Append service to empty list
	if len(serviceProvider.serviceProvider.Services) == 0 {
		serviceProvider.serviceProvider.Services = append(serviceProvider.serviceProvider.Services, message.ServiceProvider.Services[0])
		return nil
	}

	// Return error, if service already exists
	for _, s := range serviceProvider.serviceProvider.Services {
		if s.Name == message.ServiceProvider.Services[0].Name {
			return fmt.Errorf("service %q already exists in list", s.Name)
		}
	}

	// Append service to list
	serviceProvider.serviceProvider.Services = append(serviceProvider.serviceProvider.Services, message.ServiceProvider.Services[0])
	return nil
}

// serviceElection organizes and returns one service
func (serviceProvider *ServiceProvider) serviceElection() (*serviceproviders.Service, error) {

	// Empty list returns error
	if len(serviceProvider.serviceProvider.Services) == 0 {
		return nil, fmt.Errorf("empty service list")
	}

	// Returns first service found
	for _, s := range serviceProvider.serviceProvider.Services {
		if s.Status == serviceproviders.Service_SERVICE {
			return s, nil
		}
	}

	// Shuffle service list
	for i := range serviceProvider.serviceProvider.Services {
		rand.Seed(time.Now().UnixNano())
		j := rand.Intn(i + 1)
		serviceProvider.serviceProvider.Services[i], serviceProvider.serviceProvider.Services[j] = serviceProvider.serviceProvider.Services[j], serviceProvider.serviceProvider.Services[i]
	}

	// Make first candidate to service
	for _, s := range serviceProvider.serviceProvider.Services {
		if s.Status == serviceproviders.Service_CANDIDATE {
			s.Status = serviceproviders.Service_SERVICE
			return s, nil
		}
	}
	return nil, fmt.Errorf("no service candidate found")
}
