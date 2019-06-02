# Pre-Chat Using 'memberlist'

This prototype uses 

- HashiCorp's ['memberlist'](https://github.com/hashicorp/memberlist) - a SWIM++ protocol implementation
- ['liner'](https://github.com/peterh/liner) - a pure Go line editor with history

It is the first prototype who uses a bootstrap service and can act on unexpected network failure.

### Bootstrap Service With API

We have three parts:

- bootstrap-data-api
- bootstrap-data-server
- bootstrap-data-cloudfunctions
- 

- [bootstrap service](https://github.com/stefanhans/programming-reactive-systems-in-go/tree/master/prototypes/cf-chat/cf)
- [client's API of the bootstrap service](https://github.com/stefanhans/programming-reactive-systems-in-go/tree/master/prototypes/cf-chat/memberlist)
- [chat with internal commands](https://github.com/stefanhans/programming-reactive-systems-in-go/tree/master/prototypes/cf-chat/chat)

A member joins via the list of members service, which adds the new member and returns the actualized list.
Then, the client informs all old members about itself joined.





Before we see some details, let's do `go build` and run it in three terminals as follows:

```
export BOOTSTRAP_DATA_SERVER="https://europe-west1-bootstrap-data-cloudfunctions.cloudfunctions.net"

./liner-memberlist-chat alice
```

```
export BOOTSTRAP_DATA_SERVER="https://europe-west1-bootstrap-data-cloudfunctions.cloudfunctions.net"

./liner-memberlist-chat bob
```

```
export BOOTSTRAP_DATA_SERVER="https://europe-west1-bootstrap-data-cloudfunctions.cloudfunctions.net"

./liner-memberlist-chat charly
```

Write your message and send it by the return key.

### Features

- informs about joining and leaving members
- writes a log file
- set of internal commands, to which you can add new commands quickly
- every member knows all other members to avoid central client publisher as a single point of failure
- reactive service instead of well-known address of central client publisher

For more details, please see the subdirectories.


### Message Types and Services

```
message Member {
    string name = 1;
    string ip = 2;
    string port = 3;
}

message MemberList {
    repeated Member member = 1;
}

message Message {
    Member  sender = 1;
    string  text = 2;
}

// Service definition for gRPC plugin to publish messages and handle subscriptions
service Publisher {
    rpc Subscribe(Member) returns (Member) {}
    rpc Unsubscribe(Member) returns (Member) {}
    rpc Publish(Message) returns (MemberList) {}
}

// Service definition for gRPC plugin to display messages
service Displayer {
    rpc DisplayText(Message) returns (Message) {}
    rpc DisplaySubscription(Member) returns (Member) {}
    rpc DisplayUnsubscription(Member) returns (Member) {}
}
```

The first client starts the service, i.e., the Publisher, which handles all messages and subscriptions.




**Create your memberlist node**

- `memberlistconfigure` creates a memberlist configuration
- `memberlistcreate` creates the memberlist specified by the configuration


**Get, and possibly join, the bootstrap peers**

- `bootstrapjoin` joins calling peer to bootstrap peers


**Join the memberlist**

- `memberliststart` starts broadcasting between the members
- `memberlistjoin` joins bootstrap peers to memberlist


**Join the chat and say hi**

- `chatjoin` start chat listener and join the chat
- `msg hi` send a message to all chat members


**Leave all and quit**

- `chatleave` 
- `memberlistleave` 
- `memberlistshutdown`
- `bootstrapleave` 
- `quit`


