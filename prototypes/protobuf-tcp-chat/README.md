# Chat Using Protocol Buffers Via TCP

This prototype uses [Protocol Buffer](https://developers.google.com/protocol-buffers/docs/gotutorial). 
It is the first prototype with interactive commands.


Before we see some details, let's do `go build` and run it in three terminals as follows:

`./protobuf-tcp-chat alice 127.0.0.1 12340`

`./protobuf-tcp-chat bob 127.0.0.1 12341`

`./protobuf-tcp-chat charly 127.0.0.1 12342`

Write your message and send it by the return key.

`Ctrl-C` quits the chat.

### Features

- informs about joining and leaving members
- writes a log file
- set of internal commands, to which you can add new commands quickly.


### Message Types

It uses an enumeration in the message to appropriately handle the types regarding request and reply.

```
message Member {
    string name = 1;
    string ip = 2;
    string port = 3;
    bool leader = 4;
}

message MemberList {
    repeated Member member = 1;
}

// Services are mapped by request/reply message types
message Message {
    enum MessageType {
        // messages to handle subscriptions at the publishing service
        SUBSCRIBE_REQUEST = 0;
        SUBSCRIBE_REPLY = 1;

        // messages to handle unsubscriptions at the publishing service
        UNSUBSCRIBE_REQUEST = 2;
        UNSUBSCRIBE_REPLY = 3;

        // messages to publish chat messages via the publishing service
        PUBLISH_REQUEST = 4;
        PUBLISH_REPLY = 5;
    }
    MessageType msgType = 1;
    Member  sender = 2;
    string  text = 3;
    MemberList memberList = 4;
}
```

The first client starts the service running under a well-known address, and its property 'Member.leader' is true. 
The service handles all messages and subscriptions.


### Internal Commands


- `\self` shows information of this member

- `\list` lists all members; only applicable to publishing service

- `\publisher` shows well-known publisher address

- `\logfile` shows log file of this member


