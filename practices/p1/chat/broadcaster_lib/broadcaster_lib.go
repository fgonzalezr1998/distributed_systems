package broadcaster_lib

import (
	"fmt"
	"log"
)

type ClientType struct {
	client_id string
	username string
	channel ClientChannelType
	private_user string
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
	/*
	 * Add a client to the connected clients list and
	 * delete it of the disconnected clients list
	 */
	
	var new_client ClientType

	new_client.client_id = id
	new_client.username = username
	new_client.channel = c
	new_client.private_user = ""

	b.clients_list.connected = append(b.clients_list.connected, new_client)
	b.clients_list.disconnected = deleteById(b.clients_list.disconnected, id)
}

func (b * BroadcastType) DeleteClient(id string) {
	/*
	 * Delete a client of the connected clients list
	 * and add it to the disconnected clients list
	 */

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
	/*
	 * Announce the entering of a new client
	 */

	var str string = "[SERVER] Connected Clients:\n"
	for _, value := range b.clients_list.connected {
		str = str + value.username + "\n"
	}

	b.SendBroadcast("", str)
}

func (b * BroadcastType) SendBroadcast(sender_id, msg string) {
	/*
	 * Send msg to all connected clients
	 */
	
	for _, value := range b.clients_list.connected {
		if (value.client_id != sender_id) {
			value.channel <- msg
		}
	}
}

func (b * BroadcastType) SendTo(receiver_user, msg string) {
	/*
	 * Send the message 'msg' to client whose username is
	 * 'receiver_user'
	 */
	
	for _, value := range b.clients_list.connected {
		if (value.username == receiver_user) {
			value.channel <- msg
			break
		}
	}
}

func (b * BroadcastType) SetPrivateChan(id, priv_user string) {
	var client ClientType
	var index int
	var exists bool

	exists = false
	for index, client = range b.clients_list.connected {
		if (client.client_id == id) {
			exists = true
			break
		}
	}

	if (exists) {
		b.clients_list.connected[index].private_user = priv_user
	}
}

func (b * BroadcastType) Exists(username string) bool {
	/*
	 * Returns if 'username' exists at connected clients list
	 */

	for _, value := range b.clients_list.connected {
		if (value.username == username) {
			return true
		}
	}

	return false
}

func (b BroadcastType) IsInPrivate(sender_id string, output_user * string) bool {
	/*
	 * Returns if 'sender' client is in a private channel and in this case,
	 * the output user is store at 'output_user'
	 */

	 var client ClientType

	 for _, client = range b.clients_list.connected {
		if (client.client_id == sender_id) {
			break
		}
	}

	if (client.private_user != "") {
		*output_user = client.private_user
		return true
	}

	return false
}

// ---Debug functions---

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

// ---End debug functions---

// Private functions:

func deleteById(l []ClientType, id string) (list []ClientType) {
	/*
	 * Delete an element of the list by its Id given by 'id'
	 */
	
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
	/*
	 * Return the value of a client whose client_id is 'id'
	 */

	for _, value := range b.clients_list.connected {
		if (value.client_id == id) {
			*client = value
			return true
		}
	}
	return false
}