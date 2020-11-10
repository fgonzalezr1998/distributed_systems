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
	Clients_list ClientsListType
}

func (b * BroadcastType) Init() {
	b.Entering = make(chan ClientChannelType)
	b.Leaving = make(chan ClientChannelType)
	b.Messages = make(chan string)
}

func (b * BroadcastType) AddClient(id string) {
	var new_client ClientType

	new_client.client_id = id

	b.Clients_list.connected = append(b.Clients_list.connected, new_client)

	b.Clients_list.disconnected = deleteById(b.Clients_list.disconnected, id)
}

func (b BroadcastType) PrintConnectedClients() {
	for _, value := range b.Clients_list.connected {
		fmt.Println(value.client_id)
	}
	fmt.Println("----------")
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