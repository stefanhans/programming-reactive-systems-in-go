# Prototypes

According to the chapters of the thesis "Programming Reactive Systems in Go" we have the following prototypes:


|     |Thesis Section                                       |Prototype                                     |Chat Server |Group Membership|Bootstrap               |Internal Commands|Protocol Buffers|UI     |
|:---:|:---:                                                |:---:                                         |:---:       |:---:           |:---:                   |:---:   |:---:           |:---:  |
|5.1. |Chat With Protocol Buffers And gRPC                  |[gRPC-chat](gRPC-chat)                        |First client|-               |Well-known address      |-       |Yes             |'TUI'  |
|5.2. |Chat With Protocol Buffers Via TCP                   |[protobuf-tcp-chat](protobuf-tcp-chat)        |First client|-               |Well-known address      |Yes     |Yes             |'TUI'  |
|5.2. |Chat With Protocol Buffers Via UDP                   |[protobuf-udp-chat](protobuf-udp-chat)        |First client|-               |Well-known address      |Yes     |Yes             |'TUI'  |
|5.3. |Chat Using Cloud Functions As List Of Members Service|[cf-chat](cf-chat)                            |-           |Cloud Functions |Cloud Functions         |Yes     |Yes             |'TUI'  |
|5.4. |Pre-Chat Using 'memberlist'                          |[liner-memberlist-chat](liner-memberlist-chat)|-           |'memberlist'    |Cloud Functions         |Yes     |Yes             |'liner'|
|5.5. |Pre-Chat Using 'libp2p'                              |[libp2p-chat](libp2p-chat)                    |-           |'libp2p'        |Bootstrap peers         |Yes     | -              |'TUI'  |
