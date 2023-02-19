package server

import (
	"github.com/aliworkshop/grpc_chat/client"
	"github.com/aliworkshop/grpc_chat/client/data"
)

func (uc _server) HandlerPrivateMsg(c client.Client, request *data.Data) {
	c.WriteJson(request)
	if cl, ok := uc.clients[request.Id]; ok {
		cl.WriteJson(request)
	}
}

func (uc _server) HandlerChannelMsg(c client.Client, request *data.Data) {
	if channel, ok := uc.channels[request.Id]; ok {
		if c.GetKey() != channel.Admin {
			c.Write(&client.WriteRequest{
				Data: []byte("you have not access to write to channel"),
			})
			return
		}
		for _, userId := range channel.Members {
			uc.clients[userId].WriteJson(request)
		}
	}
}

func (uc _server) HandlerGroupMsg(c client.Client, request *data.Data) {
	if group, ok := uc.groups[request.Id]; ok {
		for _, userId := range group.Members {
			uc.clients[userId].WriteJson(request)
		}
	}
}
