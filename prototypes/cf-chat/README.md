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


### Internal Commands

GCP's bootstrap service API:
               
- \gcp shows all GCP data, i.e., regarding bootstrap service           
- \gcpconfig shows GCP configuration     
- \gcpsubscribe subscribes this member to bootstrap service     
- \gcpunsubscribe unsubscribes this member from bootstrap service
- \gcplist shows all members from bootstrap service
- \gcpreset resets the bootstrap service
  
Application commands:

- \all shows complete data               
- \chat shows all chat data        
- \list shows list of chat members        
- \self shows data of this member            
- \message shows message object 
- \types lists all possible message types         
- \logfile shows log file of this member
- \quit quits the application