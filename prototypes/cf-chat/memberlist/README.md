# API For List Of Members Service

`memberlist` has the API to the service and stores the data about other members for the chat.

- `Subscribe()`
- `Unsubscribe()`
- `List()`

```go
// Memberlist is the core struct for all relevant data
type Memberlist struct {
	ServiceUrl string
	Uuid       string
	Self       *IpAddress
}

// IpAddress is the struct for the Firestore
type IpAddress struct {
	Name     string `firestore:"name,omitempty"`
	Ip       string `firestore:"ip,omitempty"`
	Port     string `firestore:"port,omitempty"`
	Protocol string `firestore:"protocol,omitempty"` // "tcp" or "udp"
}
```

### Testing

```
export GCP_SERVICE_URL="https://europe-west1-cloud-functions-talk-22365.cloudfunctions.net"

go test -v
```