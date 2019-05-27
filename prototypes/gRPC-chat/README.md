# Chat Using Protocol Buffers And gRPC

This prototype uses [Protocol Buffer](https://developers.google.com/protocol-buffers/docs/gotutorial) in combination with [gRPC](https://grpc.io/docs/quickstart/go/).


Before we see some details, let's do `go build` and run it in three terminals as follows:

`./gRPC-chat alice 127.0.0.1 12340`

`./gRPC-chat bob 127.0.0.1 12341`

`./gRPC-chat charly 127.0.0.1 12342`

`Ctrl-C` quits the chat.

