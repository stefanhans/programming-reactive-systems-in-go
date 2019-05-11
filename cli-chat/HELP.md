### Help

`/help`

Usage: /help [ groups [<group>] | <command> ]

Description: help shows:
- this help text (/help)
- all groups of commands (/help groups)
- all commands of a group (/help groups <group>)
- help text of a command (/help <command>)


### Help Groups

`/help groups`

bootstrap broadcast chat cli development memberlist script test

---

`/help groups bootstrap`

bootstrap bootstrapconfig bootstrapjoin bootstrapleave bootstraplistlocal bootstrappeers bootstraprefill bootstrapreset

---

`/help groups broadcast`

broadcastadd broadcastdel broadcastlist

---

`/help groups chat`

chatleave chatmemberping chatmembers chatping chatstart chatstop msg

---

`/help groups cli`

exit log

---
`/help groups development`

init play quit

---
`/help groups memberlist`

loadconfig localnode memberlist memberlistconfigure memberlistcreate memberlistdelete memberlisthealthscore memberlistjoin memberlistleave memberlistshutdown memberlistshutdowntransport memberliststart memberlistupdate saveconfig showconfig showmemberlist

---

`/help groups script`

echo execute shell sleep

---

`/help groups test`

testfilter testfilters testlocalfilters testreset testrun testsummary testsummaryprepare

### Help Group "bootstrap"

`/help bootstrap`

Usage: /bootstrap


Description: bootstrap shows bootstrap data from remote

`/help bootstrapconfig`

Usage: /bootstrapconfig

Description: bootstrapconfig shows bootstrap configuration from remote

`/help bootstrapjoin`

Usage: /bootstrapjoin

Description: bootstrapjoin joins calling peer to bootstrap peers

`/help bootstrapleave`

Usage: /bootstrapleave

Description: bootstrapleave leaves calling peer from bootstrap peers

`/help bootstraplistlocal`

Usage: /bootstraplistlocal

Description: bootstraplistlocal lists bootstrap peers from local map

`/help bootstrappeers`

Usage: /bootstrappeers

Description: bootstrappeers lists bootstrap peers from remote

`/help bootstraprefill`

Usage: /bootstraprefill

Description: bootstraprefill refills bootstrap peers with calling peer

`/help bootstrapreset`

Usage: /bootstrapreset

Description: bootstrapreset resets the bootstrap peers at remote


### Help Group "broadcast"

`/help broadcastadd`

Usage: /broadcastadd <key> <message>

Description: broadcastadd updates a key/message at all members

`/help broadcastdel`

Usage: /broadcastdel <key>

Description: broadcastdel deletes a key at all members

`/help broadcastlist`

Usage: /broadcastlist

Description: broadcastlist lists all local key/value pairs


### Help Group "chat"

`/help chatleave`

Usage: /chatleave

Description: chatleave broadcasts deletion of this chat member

`/help chatmemberping`

Usage: /chatmemberping <member id>

Description: chatmemberping pings a chat member

`/help chatmembers`

Usage: /chatmembers <chatmember>

Description: chatmembers lists all chat members

`/help chatping`

Usage: /chatping

Description: chatping pings a member of the chat via memberlist

`/help chatstart`

Usage: /chatstart

Description: chatstart starts chat listener and broadcasts the new chat member

`/help chatstop`

Usage: /chatstop

Description: chatstop stops the chat listener

`/help msg`

Usage: /msg <string>

Description: msg sends the rest of the line as message to all other members


### Help Group "cli"

`/help exit`

Usage: /exit

Description: exit exits the application directly

`/help log`

Usage: /log (on <filename> | off)

Description: log starts or stops writing logging output in the specified file


### Help Group "development"

`/help init`

Usage: /init

Description: init simulates the init process of the application

`/help play`

Usage: /play [...]

Description: play is for developer to play

`/help quit`

Usage: /quit

Description: quit simulates the exit process of the application


### Help Group "memberlist"

`/help loadconfig`

Usage: /loadconfig file

Description: loadconfig loads the memberlist configuration from JSON file

`/help localnode`

Usage: /localnode

Description: localnode shows the local node's name and address

`/help memberlist`

Usage: /memberlist

Description: memberlist lists all members

`/help memberlistconfigure`

Usage: /memberlistconfigure

Description: memberlistconfigure creates a default memberlist configuration

`/help memberlistcreate`

Usage: /memberlistcreate

Description: memberlistcreate creates the memberlist specified by the configuration

`/help memberlistdelete`

Usage: /memberlistdelete

Description: memberlistdelete sets memberlist = nil

`/help memberlisthealthscore`

Usage: /memberlisthealthscore

Description: memberlisthealthscore shows the health score >= 0, lower numbers are better

`/help memberlistjoin`

Usage: /memberlistjoin

Description: memberlistjoin joins to memberlist

`/help memberlistleave`

Usage: /memberlistleave [<timeout in seconds, default: 1 sec>]

Description: memberlistleave broadcasts leave message until finished or timeout is reached

`/help memberlistshutdown`

Usage: /memberlistshutdown

Description: memberlistshutdown stops broadcasting to the members

`/help memberlistshutdowntransport`

Usage: /memberlistshutdowntransport

Description: memberlistshutdowntransport stops broadcasting transport to the members

`/help memberliststart`

Usage: /memberliststart

Description: memberliststart starts broadcasting to the members

`/help memberlistupdate`

Usage: /memberlistupdate [<timeout in seconds, default: 1 sec>]

Description: memberlistupdate broadcasts re-advertising the local node message until finished or timeout is reached

`/help saveconfig`

Usage: /saveconfig [file]

Description: saveconfig saves the memberlist configuration as JSON file

`/help showconfig`

Usage: /showconfig

Description: showconfig shows the memberlist configuration

`/help showmemberlist`
                                                                                      │
Usage: /showmemberlist                                                                                              │

Description: showmemberlist shows the memberlist


### Help Group "script"

`/help echo`

Usage: /echo <string>

Description: echo prints rest of line

`/help execute`

Usage: /execute file

Description: execute executes the commands in the file line by line; '#' is a comment line

`/help shell`

Usage: /shell <script>

Description: shell executes the shell script

`/help sleep`

Usage: /sleep <seconds>

Description: sleep sleeps for seconds


### Help Group "test"

`/help testfilter`

Usage: /testfilter messagesView <expected events> <expected text>

Description: testfilter adds an event filter to the messages view of the chat

`/help testfilters`

Usage: /testfilters

Description: testfilters shows all filters and events

`/help testlocalfilters`

Usage: /testlocalfilters

Description: testlocalfilters shows the local event filters

`/help testreset`

Usage: /testreset

Description: testreset resets the complete test run including the summary

`/help testrun`

Usage: /testrun

Description: testrun shows the current test run, i.e., only the not yet done commands

`/help testsummary`

Usage: /testsummary

Description: testsummary shows test summary

`/help testsummaryprepare`

Usage: /testsummaryprepare

Description: testsummaryprepare finalizes the test and prepares the test summary


---