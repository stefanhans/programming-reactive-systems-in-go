# Prototypes

According to the chapters of the thesis "Programming Reactive Systems in Go" we have the following prototypes:

- 5.1. Chat With Protocol Buffers And gRPC - [gRPC-chat](gRPC-chat)
- 5.2. Chat With Protocol Buffers Via TCP - [protobuf-tcp-chat](protobuf-tcp-chat)
- 5.2. Chat With Protocol Buffers Via UDP - [protobuf-udp-chat](protobuf-udp-chat)
- 5.3. Chat Using Cloud Functions As List Of Members Service - [cf-chat](cf-chat)
- 5.4. Pre-Chat Using 'memberlist' - [liner-memberlist-chat](liner-memberlist-chat)
- 5.5. Pre-Chat Using 'libp2p' - [libp2p-chat](libp2p-chat)




|Name                 |Chat Server |Membership     |Bootstrap               |Commands|Protocol Buffers|UI     |
|:---:                |:---:       |:---:          |:---:                   |:---:   |:---:           |:---:  |
|gRPC-chat            |First client|-              |Well-known address      |-       |Yes             |'TUI'  |
|protobuf-tcp-chat    |First client|-              |Well-known address      |Yes     |Yes             |'TUI'  |
|protobuf-udp-chat    |First client|-              |Well-known address      |Yes     |Yes             |'TUI'  |
|cf-chat              |-           |Cloud Functions|Cloud Functions         |Yes     |Yes             |'TUI'  |
|liner-memberlist-chat|-           |'memberlist'   |Cloud Functions         |Yes     |Yes             |'liner'|
|libp2p-chat          |-           |'libp2p'       |Official bootstrap peers|Yes     | -              |'TUI'  |
