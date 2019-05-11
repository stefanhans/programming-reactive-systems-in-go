# Terminal Chat As Isomorphic Client Application

The isomorphic client combines backend and frontend as well as an interactive command line 
with a framework for multi-client testing, and a chat as an example application.

[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/stefanhans/programming-reactive-systems-in-go/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/stefanhans/programming-reactive-systems-in-go/cli-chat?status.svg)](https://godoc.org/github.com/stefanhans/programming-reactive-systems-in-go/cli-chat)


The backend consists of three layers or services, respectively,
 
 - the bootstrap layer with the bootstrap API 
 - the memberlist layer using `memberlist`, an implementation of the SWIM++ protocol
 - the application layer providing a chat as an example 
 
 All layers provide a set of commands to interact with the frontend.
 
 The frontend is a terminal UI with access to the commands of the backend and a simple scripting engine.
 Additionally, it acts as a client for multi-client testing.   


### Usage of the application command

You can start the application in test mode or as fully initialized chat. The state after the test depends on its commands. 
Only in the test mode, you do have access on the test related commands. All other commands are always accessible.

```
usage: 	 ./cli-chat [-test [-testfile=<filename>]] [-logfile=<filename> | -logfile=/dev/null] <name>
```


### Start the chat normally

Providing only your nickname the chat starts all up silently.

```
./cli-chat alice
```


### Start the chat from scratch interactively

To start without any initialization, you use the test mode with an empty test.

```
./cli-chat -test -testfile=empty.cmd alice
```

Now, you can start step by step or play as wanted.


**Create your memberlist node**

- `/memberlistconfigure` creates a memberlist configuration
- `/memberlistcreate` creates the memberlist specified by the configuration


**Get, and possibly join, the bootstrap peers**

- `/bootstrapjoin` joins calling peer to bootstrap peers


**Join the memberlist**

- `/memberliststart` starts broadcasting between the members
- `/memberlistjoin` joins bootstrap peers to memberlist


**Join the chat and say hi**

- `/chatstart` start chat listener and join the chat
- `/msg hi` send a message to all chat members


**See all available commands**

- `/help` shows commands and descriptions 

[HELP.md](HELP.md) shows the complete help as well.


**Leave the cli-tool**

- `/quit` does a soft shutdown
- `/exit` does a hard exit


