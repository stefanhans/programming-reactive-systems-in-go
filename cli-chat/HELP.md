### Help

`/help`

Usage: /help [ groups [<group>] | <command> ]

Description: help shows:
- this help text (/help)
- all groups of commands (/help groups)
- all commands of a group (/help groups <group>)
- help text of a command (/help <command>)

<br>

### Help Groups

`/help groups`

bootstrap broadcast chat cli development memberlist script test

<br>

`/help groups bootstrap`

bootstrap bootstrapconfig bootstrapjoin bootstrapleave bootstraplistlocal bootstrappeers bootstraprefill bootstrapreset

<br>

`/help groups broadcast`

broadcastadd broadcastdel broadcastlist

<br>

`/help groups chat`

chatleave chatmemberping chatmembers chatping chatstart chatstop msg

<br>

`/help groups cli`

exit log

<br>

`/help groups development`

init play quit

<br>

`/help groups memberlist`

loadconfig localnode memberlist memberlistconfigure memberlistcreate memberlistdelete memberlisthealthscore memberlistjoin memberlistleave memberlistshutdown memberlistshutdowntransport memberliststart memberlistupdate saveconfig showconfig showmemberlist

<br>

`/help groups script`

echo execute shell sleep

<br>

`/help groups test`

testfilter testfilters testlocalfilters testreset testrun testsummary testsummaryprepare

<br>

### Help Group "bootstrap"

`/help bootstrap`

Usage: /bootstrap

Description: bootstrap shows bootstrap data from remote

<br>

`/help bootstrapconfig`

Usage: /bootstrapconfig

Description: bootstrapconfig shows bootstrap configuration from remote

<br>

`/help bootstrapjoin`

Usage: /bootstrapjoin

Description: bootstrapjoin joins calling peer to bootstrap peers

<br>

`/help bootstrapleave`

Usage: /bootstrapleave

Description: bootstrapleave leaves calling peer from bootstrap peers

<br>

`/help bootstraplistlocal`

Usage: /bootstraplistlocal

Description: bootstraplistlocal lists bootstrap peers from local map

<br>

`/help bootstrappeers`

Usage: /bootstrappeers

Description: bootstrappeers lists bootstrap peers from remote

<br>

`/help bootstraprefill`

Usage: /bootstraprefill

Description: bootstraprefill refills bootstrap peers with calling peer

<br>

`/help bootstrapreset`

Usage: /bootstrapreset

Description: bootstrapreset resets the bootstrap peers at remote

<br>

### Help Group "broadcast"

`/help broadcastadd`

Usage: /broadcastadd <key> <message>

Description: broadcastadd updates a key/message at all members

<br>

`/help broadcastdel`

Usage: /broadcastdel <key>

Description: broadcastdel deletes a key at all members

<br>

`/help broadcastlist`

Usage: /broadcastlist

Description: broadcastlist lists all local key/value pairs

<br>

### Help Group "chat"

`/help chatleave`

Usage: /chatleave

Description: chatleave broadcasts deletion of this chat member

<br>

`/help chatmemberping`

Usage: /chatmemberping <member id>

Description: chatmemberping pings a chat member

<br>

`/help chatmembers`

Usage: /chatmembers <chatmember>

Description: chatmembers lists all chat members

<br>

`/help chatping`

Usage: /chatping

Description: chatping pings a member of the chat via memberlist

<br>

`/help chatstart`

Usage: /chatstart

Description: chatstart starts chat listener and broadcasts the new chat member

<br>

`/help chatstop`

Usage: /chatstop

Description: chatstop stops the chat listener

<br>

`/help msg`

Usage: /msg <string>

Description: msg sends the rest of the line as message to all other members

<br>

### Help Group "cli"

`/help exit`

Usage: /exit

Description: exit exits the application directly

<br>

`/help log`

Usage: /log (on <filename> | off)

Description: log starts or stops writing logging output in the specified file

<br>

### Help Group "development"

`/help init`

Usage: /init

Description: init simulates the init process of the application

<br>

`/help play`

Usage: /play [...]

Description: play is for developer to play

<br>

`/help quit`

Usage: /quit

Description: quit simulates the exit process of the application

<br>

### Help Group "memberlist"

`/help loadconfig`

Usage: /loadconfig file

Description: loadconfig loads the memberlist configuration from JSON file

<br>

`/help localnode`

Usage: /localnode

Description: localnode shows the local node's name and address

<br>

`/help memberlist`

Usage: /memberlist

Description: memberlist lists all members

<br>

`/help memberlistconfigure`

Usage: /memberlistconfigure

Description: memberlistconfigure creates a default memberlist configuration

<br>

`/help memberlistcreate`

Usage: /memberlistcreate

Description: memberlistcreate creates the memberlist specified by the configuration

<br>

`/help memberlistdelete`

Usage: /memberlistdelete

Description: memberlistdelete sets memberlist = nil

<br>

`/help memberlisthealthscore`

Usage: /memberlisthealthscore

Description: memberlisthealthscore shows the health score >= 0, lower numbers are better

<br>

`/help memberlistjoin`

Usage: /memberlistjoin

Description: memberlistjoin joins to memberlist

<br>

`/help memberlistleave`

Usage: /memberlistleave [<timeout in seconds, default: 1 sec>]

Description: memberlistleave broadcasts leave message until finished or timeout is reached

<br>

`/help memberlistshutdown`

Usage: /memberlistshutdown

Description: memberlistshutdown stops broadcasting to the members

<br>

`/help memberlistshutdowntransport`

Usage: /memberlistshutdowntransport

Description: memberlistshutdowntransport stops broadcasting transport to the members

<br>

`/help memberliststart`

Usage: /memberliststart

Description: memberliststart starts broadcasting to the members

<br>

`/help memberlistupdate`

Usage: /memberlistupdate [<timeout in seconds, default: 1 sec>]

Description: memberlistupdate broadcasts re-advertising the local node message until finished or timeout is reached

<br>

`/help saveconfig`

Usage: /saveconfig [file]

Description: saveconfig saves the memberlist configuration as JSON file

<br>

`/help showconfig`

Usage: /showconfig

Description: showconfig shows the memberlist configuration

<br>

`/help showmemberlist`
                                                                                      │
Usage: /showmemberlist                                                                                              │

Description: showmemberlist shows the memberlist

<br>

### Help Group "script"

`/help echo`

Usage: /echo <string>

Description: echo prints rest of line

<br>

`/help execute`

Usage: /execute file

Description: execute executes the commands in the file line by line; '#' is a comment line

<br>

`/help shell`

Usage: /shell <script>

Description: shell executes the shell script

<br>

`/help sleep`

Usage: /sleep <seconds>

Description: sleep sleeps for seconds

<br>

### Help Group "test"

`/help testfilter`

Usage: /testfilter messagesView <expected events> <expected text>

Description: testfilter adds an event filter to the messages view of the chat

<br>

`/help testfilters`

Usage: /testfilters

Description: testfilters shows all filters and events

<br>

`/help testlocalfilters`

Usage: /testlocalfilters

Description: testlocalfilters shows the local event filters

<br>

`/help testreset`

Usage: /testreset

Description: testreset resets the complete test run including the summary

<br>

`/help testrun`

Usage: /testrun

Description: testrun shows the current test run, i.e., only the not yet done commands

<br>

`/help testsummary`

Usage: /testsummary

Description: testsummary shows test summary

<br>

`/help testsummaryprepare`

Usage: /testsummaryprepare

Description: testsummaryprepare finalizes the test and prepares the test summary

