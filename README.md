# Isomorphic Application System Prototype Environment


Related to the master thesis _"Programming Reactive Systems in Go"_

[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/stefanhans/golang-contexting/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/stefanhans/programming-reactive-systems-in-go)](https://goreportcard.com/report/github.com/stefanhans/programming-reactive-systems-in-go)



### Purpose

Developing distributed systems is a complex project. 
Without a central instance which is always available we face some general problems:

- who to contact in the first place
- which peers are currently online
- how does the system scale

To experiment with different libraries and to gain experience about the subject, it is useful to have a command line
running in an online connected peer. 
We can now, step-by-step, investigate the libraries and experiment freely and being open to failure. 
We gain experience in how different components interact with each other in an online network of peers.
An example application gives us a realistic use case layer.

Here the major features we need:

- command line interface within the running peers
- bootstrap data service
- group membership protocol implementation

The bootstrap service has a separated API. The implementation of the group membership protocol and the CLI tool are libraries.


### CLI Chat

The CLI Chat is a system, which has no central server and every peer can join and leave the chat at any time. 
As well a peer can leave the network unforeseeable and without informing any other peer.


### CLI Tool

The interactive command line tool has these basic features:

- command completion and history 
- interactive logging
- individual function integration
- script execution


### Bootstrap Data Service

The service consists of 

- an API `bootstrap-data-api`
- a backend for production `bootstrap-data-cloudfunctions` which is already deployed as Cloud Functions on GCP
- a backend for development `bootstrap-data-server` which can be started locally

To switch between the backend you can use `source .bootstrap-switch` which is located under `cli-chat/cli-chat`


### HashiCorps ’memberlist’ Library

We use the 'memberlist' library to inform all members about the connection information of the chat layer.
And we get a notification from it about events like joining and leaving peers.


### Chat Application

The chat application is simple - join, leave and send a message.





