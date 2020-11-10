package broadcaster_lib

type ClientType struct {
	client_id int32
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
}

func (b * BroadcastType) Init() {
	b.Entering = make(chan ClientChannelType)
	b.Leaving = make(chan ClientChannelType)
	b.Messages = make(chan string)
}