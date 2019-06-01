# Pre-Chat Using 'libp2p'

This prototype uses various packages from projects of [Protocol Labs](https://protocol.ai):

- [ipfs/go-cid](https://github.com/ipfs/go-cid)
- [ipfs/go-ipfs-addr](https://github.com/ipfs/go-ipfs-addr)
- [libp2p/go-libp2p](https://github.com/libp2p/go-libp2p)
- [libp2p/go-libp2p-host](https://github.com/libp2p/go-libp2p-host)
- [libp2p/go-libp2p-kad-dht](https://github.com/libp2p/go-libp2p-kad-dht)
- [libp2p/go-libp2p-net](https://github.com/libp2p/go-libp2p-net)
- [libp2p/go-libp2p-peer](https://github.com/libp2p/go-libp2p-peer)
- [libp2p/go-libp2p-peerstore](https://github.com/libp2p/go-libp2p-peerstore)
- [multiformats/go-multihash](https://github.com/multiformats/go-multihash)

It is inspired by [chat-with-rendezvous](https://github.com/libp2p/go-libp2p-examples/tree/master/chat-with-rendezvous)

### Get and build

```
go get github.com/programming-reactive-systems-in-go/prototypes/libp2p

go build chat.go
```

### Usage

```
Usage of ./chat:
  -b string
    	Address of bootstrap peer, <multiaddr>/<peerID>
  -r string
    	Unique string to identify group of nodes. Share this with your friends to let them connect with you
```

### Start it with a unique identifier and default bootstrap peers

Create a new uuid, e.g. by uuidgen

```
uuidgen > uuid.txt
```

<br> Use different terminal windows to run

```
./chat -r $(cat uuid.txt)
```

### Start it with your bootstrap peer

[Install IPFS](https://docs.ipfs.io/introduction/usage/#install-ipfs) if needed.

Stop the daemon, if running.

```
ipfs shutdown
```

<br> Remove all bootstrap peers and start the daemon.

```
ipfs bootstrap rm all
ipfs daemon &
```

<br> Choose peer address for private network.

```
ipfs id -f="<addrs>\n"
```

<br> Create a new uuid.

```
uuidgen > uuid.txt
```

<br> Start the application, e.g. on localhost.

```
./chat -r $(cat uuid.txt) -b $(ipfs id -f="<addrs>\n" | grep 127.0.0.1/tcp/4001)
```

### Internal commands

```
<CMD USAGE>: \rendezvous
<CMD USAGE>: \chat
<CMD USAGE>: \con
<CMD USAGE>: \peer <peer.ID Qm*...>
<CMD USAGE>: \addpeer <peer.ID Qm*...>
<CMD USAGE>: \quit
```
