package client

import (
	"github.com/aliworkshop/grpc_chat/client/event"
)

type Event struct {
	Client Client
	Event  event.Event
}
