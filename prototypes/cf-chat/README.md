# Chat Using Cloud Functions As List Of Members Service

This prototype uses [Protocol Buffer](https://developers.google.com/protocol-buffers/docs/gotutorial). 
It is the first prototype who sends its messages to all members directly.

We have three parts:

- bootstrap service, which holds a list of all members
- client's API of the bootstrap service
- chat with internal commands

A member joins via bootstrap service, which adds the new member and returns the actualized list.
Then, the client informs all old members about itself joined.


### Features

- informs about joining and leaving members
- writes a log file
- set of internal commands, to which you can add new commands quickly.
- reactive bootstrap service instead of well-known address for central client publisher
- every member knows all other members to avoid single point of failure

