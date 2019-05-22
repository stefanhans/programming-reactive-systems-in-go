### Preparations

```
cd ~/go/src/github.com/stefanhans/programming-reactive-systems-in-go/prototypes/cf-chat/chat
go build
```

### Execution

```
export GCP_SERVICE_URL="https://europe-west1-cloud-functions-talk-22365.cloudfunctions.net"

./chat alice localhost
```