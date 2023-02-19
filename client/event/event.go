package event

import "github.com/aliworkshop/grpc_chat/client/data"

type Event struct {
	Type    Type
	Request *data.Data
}

var Closed = Event{
	Type: TypeClosed,
}
