package client

import (
	"github.com/aliworkshop/grpc_chat/client/data"
	"github.com/aliworkshop/grpc_chat/client/event"
	"github.com/aliworkshop/grpc_chat/proto"
)

func (c *client) handleMessage(msg *proto.Message) {
	c.eventChan <- &Event{
		Client: c,
		Event: event.Event{
			Type: event.TypeMessage,
			Request: &data.Data{
				Action: data.Action(msg.Action),
				Type:   data.Type(msg.Type),
				Id:     msg.Id,
				Body:   msg.Body,
			},
		},
	}
}

func (c *client) handleJoin(msg *proto.Message) {
	c.eventChan <- &Event{
		Client: c,
		Event: event.Event{
			Type: event.TypeJoin,
			Request: &data.Data{
				Action: data.Action(msg.Action),
				Type:   data.Type(msg.Type),
				Id:     msg.Id,
				Body:   msg.Body,
			},
		},
	}
}
