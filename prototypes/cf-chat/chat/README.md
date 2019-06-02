# The Chat Application

Before we see some details, let's do `go build` and run it in three terminals as follows:

```
export GCP_SERVICE_URL="https://europe-west1-cloud-functions-talk-22365.cloudfunctions.net"

./chat alice localhost
```

```
export GCP_SERVICE_URL="https://europe-west1-cloud-functions-talk-22365.cloudfunctions.net"

./chat bob localhost
```

```
export GCP_SERVICE_URL="https://europe-west1-cloud-functions-talk-22365.cloudfunctions.net"

./chat charly localhost
```

Write your message and send it by the return key.


### Distributed Testing

For distributed testing we do introduce a parallel message communication using the same methodology. 
Then we can test all members currently in the chat and its communication functionality.

We have two ways implemented:

- `TEST_PUBLISH` is a two hop round trip over every chat member to every chat member each and then back to the tester.
- `TEST_CMD` is a simple call for commands to every chat member sending the results back to the tester.


Do the testing in a new terminal:

```
export GCP_SERVICE_URL="https://europe-west1-cloud-functions-talk-22365.cloudfunctions.net"

go test -v
```

You should see in the chats, e.g. `alice`:

```
<TEST_PUBLISH "TestPublishMessage">   
<TEST_REPLY "bob">                    
<TEST_REPLY "charly">                 
<TEST_CMD "list">
```


### Message Types

It uses an enumeration in the message to appropriately handle the types regarding request and reply.

```
message Member {
    string name = 1;
    string ip = 2;
    string port = 3;
    string protocol = 4;
}

message MemberList {
    repeated Member member = 1;
}

// Services are mapped by request/reply message types
message Message {
    enum MessageType {
        // messages to handle subscriptions
        SUBSCRIBE_REQUEST = 0;
        SUBSCRIBE_REPLY = 1;

        // messages to handle unsubscriptions
        UNSUBSCRIBE_REQUEST = 2;
        UNSUBSCRIBE_REPLY = 3;

        // messages to publish chat messages to all members
        PUBLISH_REQUEST = 4;
        PUBLISH_REPLY = 5;

        // messages to test publishing chat messages
        TEST_PUBLISH_REQUEST = 6;
        TEST_PUBLISH_REPLY = 7;

        // messages to test chat commands
        TEST_CMD_REQUEST = 8;
        TEST_CMD_REPLY = 9;

        // todo: implement TEST_PING functionality, i.e. response at each step

        // messages to test a chain of message transfer in a ping like fashion
        TEST_PING_REQUEST = 10;
        TEST_PING_REPLY = 11;
    }
    MessageType msgType = 1;
    Member  sender = 2;
    string  text = 3;
}
```

We have additional message types for testing. A todo is to implement the ping tests.

### Internal Commands

GCP's bootstrap service API:
               
- `\gcp` shows all GCP data, i.e., regarding bootstrap service           
- `\gcpconfig` shows GCP configuration     
- `\gcpsubscribe` subscribes this member to bootstrap service     
- `\gcpunsubscribe` unsubscribes this member from bootstrap service
- `\gcplist` shows all members from bootstrap service
- `\gcpreset` resets the bootstrap service
  
Application commands:

- `\all` shows complete data               
- `\chat` shows all chat data        
- `\list` shows list of chat members        
- `\self` shows data of this member            
- `\message` shows message object 
- `\types` lists all possible message types         
- `\logfile` shows log file of this member
- `\quit` quits the application