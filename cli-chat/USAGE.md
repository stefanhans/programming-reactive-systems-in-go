This information is deprecated

```
bootstrapjoin
	 bootstrapjoin joins calling peer to bootstrap peers

bootstrapleave
	 bootstrapleave leave calling peer from bootstrap peers

bootstraplist
	 bootstraplist  list bootstrap peers from remote

bootstraplistlocal
	 bootstraplistlocal list bootstrap peers from local map

bootstraprefill
	 bootstraprefill refill bootstrap peers with calling peer

bootstrapreset
	 bootstrap joins peer to bootstrap peers

broadcastadd <key> <message>
	 broadcastadd updates a key/message at all members

broadcastdel <key> <message>
	 broadcastdel deletes a key at all members

broadcastlist
	 broadcastlist lists all local key/value pairs

chatjoin
	 chatjoin starts chat listener and broadcasts new chat member

chatleave
	 chatleave broadcasts deletion of this chat member

chatmemberlist
	 chatmemberlist lists all chat members

chatmemberping
	 chatmemberping ping a chat member

echo text_w/o_linebreak
	 echo prints rest of line

execute file
	 execute execute the commands in the file line by line, '#' is comment

loadconfig file
	 loadconfig load the memberlist configuration from JSON file

localnode
	 localnode shows the local node's name and address

log (on <filename>)|off
	 log starts or stops writing logging output in the specified file

memberlist
	 memberlist lists all members

memberlistconfigure
	 memberlistconfigure creates a memberlist configuration

memberlistcreate
	 memberlistcreate creates the memberlist specified by the configuration

memberlistdelete
	 memberlistdelete sets memberlist = nil

memberlisthealthscore
	 memberlisthealthscore shows the health score >= 0, lower numbers are better

memberlistjoin [<members> ...]
	 memberlistjoin add oneself or other member(s) to the memberlist

memberlistleave [<timeout in seconds, default: 1 sec>]
	 memberlistleave broadcasts leave message until finished or timeout is reached

memberlistshutdown
	 memberlistshutdown stops broadcasting to the members

memberlistshutdowntransport
	 memberlistshutdowntransport stops broadcasting transport to the members

memberliststart
	 memberliststart starts broadcasting to the members

memberlistupdate [<timeout in seconds, default: 1 sec>]
	 memberlistupdate broadcasts re-advertising the local node message until finished or timeout is reached

msg
	 msg sends the rest of the line as message to all other members

ping <chatmember>
	 ping pings a member of the chat

play
	 for developer playing

quit
	 close the session and exit

saveconfig [file]
	 saveconfig saves the memberlist configuration as JSON file

showconfig
	 showconfig shows the memberlist configuration

showmemberlist
	 showmemberlist shows the memberlist

sleep seconds
	 sleep sleeps for seconds
```