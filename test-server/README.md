# HTTP Server As Sidecar For Multi-Client Tests

The server acts as a sidecar and coordinates the execution of multi-client tests.

[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/stefanhans/programming-reactive-systems-in-go/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/stefanhans/programming-reactive-systems-in-go/test-server/cli-test-server?status.svg)](https://godoc.org/github.com/stefanhans/programming-reactive-systems-in-go/test-server/cli-test-server)

Testing the interactions of isomorphic clients is complex. 
Therefore, the testing framework has to be as reliable and straightforward as possible.

You want to test the clients and not the test framework. So, we start with a very simple approach.

- provide a test file with each line defining one command for one client
- provide a command to set a text filter on the terminal view of the chat messages

Now, you can make end-to-end tests like this:

```bash
alice init
bob init
charly init
alice testfilter messagesView 1 <bob> test
alice testfilter messagesView 1 <charly> test
bob testfilter messagesView 1 <alice> test
bob testfilter messagesView 1 <charly> test
charly testfilter messagesView 1 <alice> test
charly testfilter messagesView 1 <bob> test
alice msg test
bob msg test
charly msg test
```

In this test, we have three clients, alice, bob, and charly. All do the same. 
They initialize the chat, set a filter expecting a message from the other two, and send a message themselves.

A successful test shows the following summary:

```
│Summary of "default" (8e389cce-7301-11e9-b4f3-acde48001122)                           │
│--------------------------------------------------------------------------            │
│command  "charly init" OK                                                             │
│command  "charly testfilter messagesView 1 <alice> test" OK                           │
│event  "charly testfilter messagesView 1 <alice> test" OK                             │
│command  "charly testfilter messagesView 1 <bob> test" OK                             │
│event  "charly testfilter messagesView 1 <bob> test" OK                               │
│command  "charly msg test" OK                                                         │
│command  "alice testfilter messagesView 1 <bob> test" OK                              │
│event  "alice testfilter messagesView 1 <bob> test" OK                                │
│command  "alice testfilter messagesView 1 <charly> test" OK                           │
│event  "alice testfilter messagesView 1 <charly> test" OK                             │
│command  "alice msg test" OK                                                          │
│command  "bob testfilter messagesView 1 <alice> test" OK                              │
│event  "bob testfilter messagesView 1 <alice> test" OK                                │
│command  "bob testfilter messagesView 1 <charly> test" OK                             │
│event  "bob testfilter messagesView 1 <charly> test" OK                               │
│command  "bob msg test" OK                                                            │
│--------------------------------------------------------------------------            │
```

You see the commands with their return values and the events showing the match with their filters.

The implementation is fully serialized. All clients load the newest information from the test server in a loop.
If the first line starts with their name, they execute the rest of the line as command and send the result back.
Then, they send a request to remove the executed command.

The filter sends every matching event to the test server. In the end, the server summarizes the data.