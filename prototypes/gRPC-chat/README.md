# Chat Using Protocol Buffers And gRPC

This prototype uses [Protocol Buffer](https://developers.google.com/protocol-buffers/docs/gotutorial) in combination with [gRPC](https://grpc.io/docs/quickstart/go/). It was the first prototype we had designed.

Before we see some details, let's do `go build` and run it in three terminals as follows:

`./gRPC-chat alice 127.0.0.1 12340`

`./gRPC-chat bob 127.0.0.1 12341`

`./gRPC-chat charly 127.0.0.1 12342`

Write your message and send it by the return key.

`Ctrl-C` quits the chat.

### Features

- informs about joining and leaving members
- writes a logfile.

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
