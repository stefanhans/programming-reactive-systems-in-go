# Chat Using Cloud Functions As List Of Members Service

This prototype uses [Google Cloud Functions](https://cloud.google.com/functions/docs). 
It is the first prototype who sends its messages to all members directly.

We have three parts:

- [list of members service](https://github.com/stefanhans/programming-reactive-systems-in-go/tree/master/prototypes/cf-chat/cf)
- [client's API of the list of members service](https://github.com/stefanhans/programming-reactive-systems-in-go/tree/master/prototypes/cf-chat/memberlist)
- [chat with internal commands](https://github.com/stefanhans/programming-reactive-systems-in-go/tree/master/prototypes/cf-chat/chat)

A member joins via the list of members service, which adds the new member and returns the actualized list.
Then, the client informs all old members about itself joined.


### Features

- informs about joining and leaving members
- writes a log file
- set of internal commands, to which you can add new commands quickly
- every member knows all other members to avoid central client publisher as a single point of failure
- reactive service instead of well-known address of central client publisher

For more details, please see the subdirectories.
