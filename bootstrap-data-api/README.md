[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/stefanhans/programming-reactive-systems-in-go/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-api?status.svg)](https://godoc.org/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-api)

We test the API against two different backend services:

- a webserver running on localhost:8080
- GCP Cloud Functions

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

export BOOTSTRAP_DATA_SERVER="https://europe-west1-bootstrap-data-cloudfunctions.cloudfunctions.net"
go test -run TestCf
```