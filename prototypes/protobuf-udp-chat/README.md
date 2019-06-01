# Chat Using Protocol Buffers Via UDP

This prototype uses [Protocol Buffer](https://developers.google.com/protocol-buffers/docs/gotutorial). 
It is mainly the same prototype as `protobuf-tcp-chat`, but using UDP instead of TCP.


Before we see some details, let's do `go build` and run it in three terminals as follows:

`./protobuf-udp-chat alice 127.0.0.1 12340`

`./protobuf-udp-chat bob 127.0.0.1 12341`

`./protobuf-udp-chat charly 127.0.0.1 12342`

Write your message and send it by the return key.

`Ctrl-C` quits the chat.

### UDP Instead Of TCP

We use protocol "udp" instead of "tcp" in the `net` package functions.

```
publishingListener, err := net.ListenPacket("udp", publishingService)
conn, err := net.Dial("udp", publishingService)

displayingListener, err := net.ListenPacket("udp", displayingService)
conn, err := net.Dial("udp", recipient)
```

We take care of the buffer size according to [RFC 791](https://tools.ietf.org/html/rfc791).

```
// The maximum safe UDP payload is 508 bytes.
// This is a packet size of 576 (IPv4 minimum reassembly buffer size),
// minus the maximum 60-byte Ip header and the 8-byte UDP header.
bufferSize = 508
```

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


