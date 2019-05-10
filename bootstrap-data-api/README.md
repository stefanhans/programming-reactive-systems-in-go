# Bootstrap Service API

The integrated API provides bootstrap information about which peer to contact to join the decentralized overlay network.

[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/stefanhans/programming-reactive-systems-in-go/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-api?status.svg)](https://godoc.org/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-api)


### The API definition is small and simple

There is only one JSON Data Object.

```json
{
  "Config": {
    "maxpeers": 2,
    "minrefillcandidates": 2,
    "numpeers": 2
  },
  "Peers": {
    "71c3c0ea-e1bd-4112-97ac-1cfeeaa07bc3": {
      "id": "71c3c0ea-e1bd-4112-97ac-1cfeeaa07bc3",
      "name": "bob-Stefans-MBP-94d785ce-734a-11e9-a0c2-acde48001122",
      "ip": "127.0.0.1",
      "port": "56471",
      "protocol": "tcp",
      "status": "0",
      "timestamp": "1557509977"
    },
    "a559fd29-0410-4cac-b1f1-109a9513f22f": {
      "id": "a559fd29-0410-4cac-b1f1-109a9513f22f",
      "name": "alice-Stefans-MBP-939dbaf2-734a-11e9-847f-acde48001122",
      "ip": "127.0.0.1",
      "port": "56461",
      "protocol": "tcp",
      "status": "0",
      "timestamp": "1557509975"
    }
  }
}
```

The interface itself is quite straightforward and returns the JSON object. Except `ping` which returns only "OK."

```go
// Run starts the service.
func Run() {

	http.HandleFunc("/join", Join)
	http.HandleFunc("/leave", Leave)
	http.HandleFunc("/refill", Refill)
	http.HandleFunc("/list", List)
	http.HandleFunc("/reset", Reset)
	http.HandleFunc("/ping", Ping)
	http.HandleFunc("/config", ConfigUpdate)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

But, `refill` needs some explanation. As configured by `maxpeers`, we store up to maximum.
If one of these bootstrap peers is leaving the network, all others are triggered to send a request to refill. 
The server chooses the peer with the oldest timestamp or uptime, respectively, to complete the list.


### Testing

We test the API against two different backend services:

#### A HTTP server running on localhost

Start the webserver:

```
cd ~/go/src/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-server

go build
./bootstrap-data-server
```

Now, execute the following tests:

```
cd ~/go/src/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-api

export BOOTSTRAP_DATA_SERVER="http://localhost:8080"
go test -run TestLocalhost
```

#### GCP Cloud Functions

```
cd ~/go/src/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-api

export BOOTSTRAP_DATA_SERVER="https://europe-west1-bootstrap-data-cloudfunctions.cloudfunctions.net"
go test -run TestCf
```

The cloud functions are already deployed. Nevertheless, you find the source code with instructions on how to deploy 
onto the Google Cloud Platform in `bootstrap-data-cloud-functions`.