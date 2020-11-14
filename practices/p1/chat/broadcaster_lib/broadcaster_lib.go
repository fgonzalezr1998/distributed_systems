package broadcaster_lib

import (
	"fmt"
	"log"
)

type ClientType struct {
	client_id string
	username string
	channel ClientChannelType
}

type ClientsListType struct {
	connected []ClientType
	disconnected []ClientType
}

type ClientChannelType chan<- string // Write-only channel

type BroadcastType struct {
	clients_list ClientsListType
}

// Public Functions:

func (b * BroadcastType) AddClient(id string, username string, c chan<- string) {
	var new_client ClientType

	new_client.client_id = id
	new_client.channel = c
	new_client.username = username

	b.clients_list.connected = append(b.clients_list.connected, new_client)
	b.clients_list.disconnected = deleteById(b.clients_list.disconnected, id)
}

func (b * BroadcastType) DeleteClient(id string) {
	var client ClientType

	if (!b.getClient(id, &client)) {
		log.Print("[ERROR] Client doesn't exist!\n")
	}

	// Close its channel before to delete it:

	close(client.channel)

	b.clients_list.connected = deleteById(b.clients_list.connected, id)
	b.clients_list.disconnected = append(b.clients_list.disconnected, client)
}

func (b * BroadcastType) AnnounceClients() {
	var str string = "Connected Clients:\n"
	for _, value := range b.clients_list.connected {
		str = str + value.username + "\n"
	}

	b.SendBroadcast("", str)
}

func (b * BroadcastType) SendBroadcast(sender, msg string) {
	for _, value := range b.clients_list.connected {
		if (value.client_id != sender) {
			value.channel <- msg
		}
	}
}

func (b BroadcastType) PrintConnectedClients() {
	for _, value := range b.clients_list.connected {
		fmt.Println(value.client_id + ":" + value.username)
	}
}

func (b BroadcastType) PrintDisconnectedClients() {
	for _, value := range b.clients_list.disconnected {
		fmt.Println(value.client_id)
	}
}

// Private functions:

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

func (b BroadcastType) getClient(id string, client * ClientType) bool {
	for _, value := range b.clients_list.connected {
		if (value.client_id == id) {
			*client = value
			return true
		}
	}
	return false
}