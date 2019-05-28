# Prototypes

According to the chapters of the thesis "Programming Reactive Systems in Go" we have the following prototypes:

- 5.1. [Chat With Protocol Buffers And gRPC](gRPC-chat)
- 5.2. [Chat With Protocol Buffers Via TCP](protobuf-tcp-chat)
- 5.2. [Chat With Protocol Buffers Via UDP](protobuf-udp-chat)
- 5.3. [Chat Using Cloud Functions As Bootstrap Service](cf-chat)
- 5.4. [Pre-Chat Using 'memberlist'](liner-memberlist-chat)
- 5.5. [Pre-Chat Using 'libp2p'](libp2p-chat)




|Name                 |Chat Server      |Bootstrap               |Commands|Protocol Buffers|'TUI'|'liner'|'memberlist'|'libp2p'|
|:---:                |:---:            |:---:                   |:---:   |:---:           |:---:|:---:  |:---:       |:---:   |
|gRPC-chat            |First application|Well-known address      |-       |Yes             |Yes  |-      |-           |-       |
|protobuf-tcp-chat    |First application|Well-known address      |Yes     |Yes             |Yes  |-      |-           |-       |
|protobuf-udp-chat    |First application|Well-known address      |Yes     |Yes             |Yes  |-      |-           |-       |
|cf-chat              |-                |Cloud Functions         |Yes     |Yes             |Yes  |-      |-           |-       |
|liner-memberlist-chat|-                |Cloud Functions         |Yes     |Yes             |-    |Yes    |Yes         |-       |
|libp2p-chat          |-                |Official bootstrap peers|Yes     | -              |-    |-      |-           |Yes     |
