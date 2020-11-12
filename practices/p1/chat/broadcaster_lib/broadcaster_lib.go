package broadcaster_lib

import (
	"fmt"
)

type ClientType struct {
	client_id string
}

type ClientsListType struct {
	connected []ClientType
	disconnected []ClientType
}

type ClientChannelType chan<- string // Write-only channel

type BroadcastType struct {
	Entering chan ClientChannelType
	Leaving chan ClientChannelType
	Messages chan string
	clients_list ClientsListType
}

func (b * BroadcastType) Init() {
	b.Entering = make(chan ClientChannelType)
	b.Leaving = make(chan ClientChannelType)
	b.Messages = make(chan string)
}

func (b * BroadcastType) AddClient(id string) {
	var new_client ClientType

	new_client.client_id = id

	b.clients_list.connected = append(b.clients_list.connected, new_client)
	b.clients_list.disconnected = deleteById(b.clients_list.disconnected, id)
}

func (b * BroadcastType) DeleteClient(id string) {
	var new_client ClientType

	new_client.client_id = id

	b.clients_list.connected = deleteById(b.clients_list.connected, id)
	b.clients_list.disconnected = append(b.clients_list.disconnected, new_client)
}

func (b BroadcastType) PrintConnectedClients() {
	for _, value := range b.clients_list.connected {
		fmt.Println(value.client_id)
	}
}

func (b BroadcastType) PrintDisconnectedClients() {
	for _, value := range b.clients_list.disconnected {
		fmt.Println(value.client_id)
	}
}

func deleteById(l []ClientType, id string) (list []ClientType) {
	list = l
	for index, value := range l {
		if (value.client_id == id) {
			list = removeIndex(l, index)
			break
		}
	}

	return list
}

func removeIndex(l []ClientType, index int) (list [] ClientType) {
	list = append(l[:index], l[index+1:]...)
	return list
}