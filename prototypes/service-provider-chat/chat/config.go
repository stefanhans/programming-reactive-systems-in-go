package main

import (
	"fmt"
	"net"

	"github.com/stefanhans/programming-reactive-systems-in-go/prototypes/service-provider-chat/chat/chat-group"
)

type ServiceConfig struct {
	Name     string
	Protocol string // "tcp"."udp"
	Ip       string // 127.0.0.1
	Port     string // 22365
}

type Config struct {
	Name            string
	ServiceProvider *ServiceConfig // well-known service "127.0.0.1:22365" as entry point
	ChatService     *ServiceConfig // first member as chat server
	ChatListener    *ServiceConfig
}

func (config *Config) String() string {
	out := fmt.Sprintf("\tName: %q\n", config.Name)
	out += fmt.Sprintf("\tServiceProvider: %v\n", *config.ServiceProvider)
	out += fmt.Sprintf("\tChatService: %v\n", *config.ChatService)
	out += fmt.Sprintf("\tChatListener: %v", *config.ChatListener)
	return out
}

func (config *Config) SetServiceProvider(ip, port string) {
	config.ServiceProvider.Ip = ip
	config.ServiceProvider.Port = port
}

func (config *Config) ServiceProviderAddress() string {
	return config.ServiceProvider.Ip + ":" + config.ServiceProvider.Port
}

func (config *Config) SetChatService(name, ip, port string) {
	config.ChatService.Name = name
	config.ChatService.Ip = ip
	config.ChatService.Port = port
}

func (config *Config) ChatServiceAddress() string {
	return config.ChatService.Ip + ":" + config.ChatService.Port
}

func (config *Config) SetChatListener(ip, port string) {
	config.ChatListener.Ip = ip
	config.ChatListener.Port = port
}

func (config *Config) ChatListenerAddress() string {
	return config.ChatListener.Ip + ":" + config.ChatListener.Port
}

var (

	// Application identity set by command args
	displayingService string
	selfMember        *chatgroup.Member

	// Publisher storage for member of chat group
	// todo refactor chatgroup.memberlist instead of []*chatgroup.Member
	cgMember []*chatgroup.Member

	//
	logfilename string
)

var requestActionMap = map[chatgroup.Message_MessageType]func(*chatgroup.Message, net.Addr) error{
	chatgroup.Message_SUBSCRIBE_REQUEST:   handleSubscribeRequest,
	chatgroup.Message_UNSUBSCRIBE_REQUEST: handleUnsubscribeRequest,
	chatgroup.Message_PUBLISH_REQUEST:     handlePublishRequest,
}

var replyActionMap = map[chatgroup.Message_MessageType]func(*chatgroup.Message) error{
	chatgroup.Message_SUBSCRIBE_REPLY:   handleSubscribeReply,
	chatgroup.Message_UNSUBSCRIBE_REPLY: handleUnsubscribeReply,
	chatgroup.Message_PUBLISH_REPLY:     handlePublishReply,
}
