# Chat Using Cloud Functions As List Of Members Service

This prototype uses [Google Cloud Functions](https://cloud.google.com/functions/docs). 
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
- set of internal commands, to which you can add new commands quickly
- every member knows all other members to avoid central client publisher as a single point of failure
- reactive service instead of well-known address of central client publisher

For more details, please see the subdirectories.
