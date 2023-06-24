package server

import (
	"github.com/aliworkshop/grpc_chat/client"
	"github.com/aliworkshop/grpc_chat/client/data"
)

func (uc *_server) HandleJoinGroup(c client.Client, request *data.Data) {
	if group, ok := uc.groups[request.Id]; ok {
		if !keyExistsInArray(c.GetKey(), group.Members) {
			group.Members = append(group.Members, c.GetKey())
			uc.groups[request.Id] = group
		}
	}
}

func (uc *_server) HandleJoinChannel(c client.Client, request *data.Data) {
	if channel, ok := uc.channels[request.Id]; ok {
		if !keyExistsInArray(c.GetKey(), channel.Members) {
			channel.Members = append(channel.Members, c.GetKey())
			uc.channels[request.Id] = channel
		}
	}
}
