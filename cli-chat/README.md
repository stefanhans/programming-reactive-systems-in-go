# Terminal Application As Isomorphic Chat

[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/stefanhans/programming-reactive-systems-in-go/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/stefanhans/programming-reactive-systems-in-go/cli-chat?status.svg)](https://godoc.org/github.com/stefanhans/programming-reactive-systems-in-go/cli-chat)


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

- `memberlistconfigure` creates a memberlist configuration
- `memberlistcreate` creates the memberlist specified by the configuration


**Get, and possibly join, the bootstrap peers**

- `bootstrapjoin` joins calling peer to bootstrap peers


**Join the memberlist**

- `memberliststart` starts broadcasting between the members
- `memberlistjoin` joins bootstrap peers to memberlist


**Join the chat and say hi**

- `chatstart` start chat listener and join the chat
- `msg hi` send a message to all chat members


**See all available commands**

- `<tab><tab>` lists all commands using code completion
- `help` or other not existing command shows commands and descriptions


**Leave the cli-tool**

- `quit`


