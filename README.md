# A System Prototype For Isomorphic Client Applications


This software project is related to the master thesis _"Programming Reactive Systems in Go"_

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

- command line interface within the running peer
- bootstrap data service
- group membership protocol implementation

The bootstrap service has a separated API. The implementation of the group membership protocol and the CLI tool are libraries.

We try to have seamless transitions from getting familiar with libraries, developing, testing, and production ready. 

Additional to some internal CLI features, we do provide script execution and multi-client testing.


### Bootstrap Data Service

The service consists of 

- an API `bootstrap-data-api`
- a backend for production `bootstrap-data-cloudfunctions` which is already deployed as Cloud Functions on GCP
- a backend for development `bootstrap-data-server` which can be started locally


### HashiCorps ’memberlist’ Library

We use the 'memberlist' library to inform all members about the connection information of the chat layer.
And we get a notification from it about events like joining and leaving peers.


### Chat Application

The chat application is simple - join, leave and send a message.


### CLI Features

During the development of isomorphic client applications, we need a simple way to organize investigation, implementation, and manually testing.
Therefore, we can quickly provide any function to the internal command line. This is our essential CLI feature. 

Following the [KISS](https://en.wikipedia.org/wiki/KISS_principle) and [YAGNI](https://en.wikipedia.org/wiki/You_aren%27t_gonna_need_it) 
principles, we have now:

- advanced logging 
- execution of internal scripts and shell scripts
- saving and loading of internal configurations

All can easily be enhanced.


### Multi-Client Testing

For the tests with interacting clients, we have `test-server` doing the coordination.
The clients have an integrated API and an own set of internal commands.
The tests itself are scripts provided to the test server.



