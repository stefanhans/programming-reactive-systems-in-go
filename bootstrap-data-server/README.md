# HTTP Server As Bootstrap Service

The server collects and coordinates the bootstrap information for the clients.

[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/stefanhans/programming-reactive-systems-in-go/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-server/bootstrap-server?status.svg)](https://godoc.org/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-server/bootstrap-server)

**Start the webserver**

```
cd ~/go/src/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-server

go build
./bootstrap-data-server
```

**Test manually**

```
curl -s http://localhost:8080/ping

curl -sd "3 1"  http://localhost:8080/config

curl -sd "alice-id alice 127.0.0.1 12340 tcp test 1"  http://localhost:8080/join
curl -sd "bob-id bob 127.0.0.1 12341 tcp test 2"  http://localhost:8080/join

curl -sd "alice-id" http://localhost:8080/leave
curl -s http://localhost:8080/list

curl -sd "2 1"  http://localhost:8080/config
curl -sd "alice-id alice 127.0.0.1 12341 tcp test 1552655444"  http://localhost:8080/refill

curl -s http://localhost:8080/reset
```

**Test via API**

```
cd ~/go/src/github.com/stefanhans/programming-reactive-systems-in-go/bootstrap-data-api

export BOOTSTRAP_DATA_SERVER="http://localhost:8080"
go test -run TestLocalhost
```
