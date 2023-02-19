package domain

import (
	"context"
	"github.com/aliworkshop/errorslib"
	"github.com/aliworkshop/grpc_chat/client"
	"github.com/aliworkshop/grpc_chat/client/data"
	"github.com/aliworkshop/grpc_chat/proto"
)

type RequestHandle func(c client.Client, request *data.Data)
type JoinHandler func(c client.Client, request *data.Data)

type ChatUc interface {
	proto.ChatServer
	Subscribe(userId string, conn proto.Chat_StreamServer) (client.Client, errorslib.ErrorModel)
	RegisterRequestHandler(t data.Type, handle RequestHandle)
	RegisterJoinHandler(t data.Type, handle JoinHandler)
	Start()
	GetClientByKey(key string) client.Client
	Run(ctx context.Context) errorslib.ErrorModel
}
