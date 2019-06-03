# Pre-Chat Using 'memberlist'

This prototype uses the following external libraries:

- HashiCorp's [`memberlist`](https://github.com/hashicorp/memberlist) - a SWIM++ protocol implementation
- [`liner`](https://github.com/peterh/liner) - a pure Go line editor with history

It is the first prototype who uses a bootstrap service and can act on unexpected network failure.
### Configure The Bootstrap Service

A member joins via bootstrap service. The service adds the new peer, if not enough peers stored already, 
and returns the list in any case. The bootstrap peers are entry points to join the `memberlist`.

The bootstrap service has two implementations:

- [bootstrap-data-server](https://github.com/stefanhans/programming-reactive-systems-in-go/tree/master/prototypes/liner-memberlist-chat/bootstrap-data-server)
is an implementation as HTTP server e.g. running on localhost for development
- [bootstrap-data-cloudfunctions](https://github.com/stefanhans/programming-reactive-systems-in-go/tree/master/prototypes/liner-memberlist-chat/bootstrap-data-cloudfunctions)
is an implementation deployed as Google Cloud Functions

The [bootstrap-data-api](https://github.com/stefanhans/programming-reactive-systems-in-go/tree/master/prototypes/liner-memberlist-chat/bootstrap-data-api)
completes the service. 

Execute `. .bootstrap-switch` to switch between both implementations. In case you use the server, you have to start the server manually as described [here](https://github.com/stefanhans/programming-reactive-systems-in-go/tree/master/prototypes/liner-memberlist-chat/bootstrap-data-server).

### Run The Application

After choosing the bootstrap service implementation, let's do `go build` and run it in three terminals as follows - 
do not forget to execute `. .bootstrap-switch` in every terminal appropriately:

```
./liner-memberlist-chat alice
```

```
./liner-memberlist-chat bob
```

```
./liner-memberlist-chat charly
```

### UI And Internal Commands

The started application returns an internal prompt and has not much initialization done. The inner command line 
provides command completion and history.

Type `<Tab><Tab>` to see all commands. 

```
bootstrapjoin               memberlistconfigure
bootstrapleave              memberlistcreate
bootstraplist               memberlistdelete
bootstraplistlocal          memberlisthealthscore
bootstraprefill             memberlistjoin
bootstrapreset              memberlistleave
broadcastadd                memberlistshutdown
broadcastdel                memberlistshutdowntransport
broadcastlist               memberliststart
chatjoin                    memberlistupdate
chatleave                   msg
chatmemberlist              ping
chatmemberping              play
echo                        quit
execute                     saveconfig
loadconfig                  showconfig
localnode                   showmemberlist
log                         sleep
memberlist
```

### Start The Chat Manually

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


### Start The Chat Via Internal Script

If you want to start the chat, you can use `execute chat.txt`, and every line of `chat.txt` is executed one by one:

```
memberlistconfigure
memberlistcreate
bootstrapjoin
memberliststart
memberlistjoin
chatjoin
chatmemberlist
sleep 1
msg hello
```




