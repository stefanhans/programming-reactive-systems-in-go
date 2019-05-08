[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/stefanhans/programming-reactive-systems-in-go/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/stefanhans/programming-reactive-systems-in-go/cli-chat)](https://goreportcard.com/report/github.com/stefanhans/programming-reactive-systems-in-go/cli-chat)
[![GoDoc](https://godoc.org/github.com/stefanhans/programming-reactive-systems-in-go/cli-chat?status.svg)](https://godoc.org/github.com/stefanhans/programming-reactive-systems-in-go/cli-chat)

###Start the chat interactively

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

