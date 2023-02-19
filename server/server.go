package server

import (
	"context"
	"github.com/aliworkshop/errorslib"
	"github.com/aliworkshop/grpc_chat/client"
	"github.com/aliworkshop/grpc_chat/client/data"
	"github.com/aliworkshop/grpc_chat/client/event"
	"github.com/aliworkshop/grpc_chat/domain"
	"github.com/aliworkshop/grpc_chat/proto"
	"github.com/aliworkshop/loggerlib/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type _server struct {
	proto.UnimplementedChatServer
	log             logger.Logger
	started         bool
	eventChan       chan *client.Event
	requestHandlers map[data.Type][]domain.RequestHandle
	joinHandlers    map[data.Type][]domain.JoinHandler

	clients  map[string]client.Client
	groups   map[string]domain.Group
	channels map[string]domain.Channel
}

const userHeader = "x-user"

func NewServer() domain.ChatUc {
	uc := &_server{
		eventChan:       make(chan *client.Event),
		requestHandlers: make(map[data.Type][]domain.RequestHandle),
		joinHandlers:    make(map[data.Type][]domain.JoinHandler),
		clients:         make(map[string]client.Client),
		groups:          make(map[string]domain.Group),
		channels:        make(map[string]domain.Channel),
	}
	uc.groups["1"] = domain.Group{
		Name: "first group",
		Id:   "1",
	}
	uc.channels["1"] = domain.Channel{
		Name:  "first channel",
		Id:    "1",
		Admin: "15",
	}
	uc.RegisterRequestHandler(data.User, uc.HandlerPrivateMsg)
	uc.RegisterRequestHandler(data.Group, uc.HandlerGroupMsg)
	uc.RegisterRequestHandler(data.Channel, uc.HandlerChannelMsg)

	uc.RegisterJoinHandler(data.Group, uc.HandleJoinGroup)
	uc.RegisterJoinHandler(data.Channel, uc.HandleJoinChannel)
	return uc
}

func (uc *_server) RegisterRequestHandler(t data.Type, handle domain.RequestHandle) {
	if uc.started {
		panic("can not register handler while websocket is started")
	}
	handlers := uc.requestHandlers[t]
	if handlers == nil {
		handlers = make([]domain.RequestHandle, 0)
	}
	handlers = append(handlers, handle)
	uc.requestHandlers[t] = handlers
}

func (uc *_server) RegisterJoinHandler(t data.Type, handle domain.JoinHandler) {
	if uc.started {
		panic("can not register handler while websocket is started")
	}
	handlers := uc.joinHandlers[t]
	if handlers == nil {
		handlers = make([]domain.JoinHandler, 0)
	}
	handlers = append(handlers, handle)
	uc.joinHandlers[t] = handlers
}

func (uc *_server) Start() {
	if uc.started {
		return
	}
	uc.started = true
	go uc.start()
}

func (uc *_server) GetClientByKey(key string) client.Client {
	return uc.clients[key]
}

func (uc *_server) start() {
	for uc.started {
		select {
		case e := <-uc.eventChan:
			if e == nil {
				return
			}
			switch e.Event.Type {
			case event.TypeMessage:
				handlers := uc.requestHandlers[e.Event.Request.Type]
				for _, h := range handlers {
					go h(e.Client, e.Event.Request)
				}
				break
			case event.TypeClosed:
				//delete(uc.clients, e.Client.GetKey())
				break
			case event.TypeJoin:
				handlers := uc.joinHandlers[e.Event.Request.Type]
				for _, h := range handlers {
					go h(e.Client, e.Event.Request)
				}
				break
			}
			break
		}
	}
}

func (uc *_server) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
	return &proto.LoginResponse{}, nil
}

func (uc *_server) Stream(srv proto.Chat_StreamServer) error {
	userId, ok := uc.extractToken(srv.Context())
	if !ok {
		return status.Error(codes.Unauthenticated, "missing user header")
	}
	uc.Subscribe(userId, srv)

	/*for {
		req, err := srv.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		go uc.handleTextMsg(req, srv)
	}*/

	<-srv.Context().Done()
	return srv.Context().Err()
}

/*
func (uc *_server) handleTextMsg(msg *proto.Message, srv proto.Chat_StreamServer) {
	switch msg.Action {
	case data.Message.String():
		uc.handleMessage(msg, srv)
	case data.Join.String():
		uc.handleJoin(msg, srv)
	}
}

func (uc *_server) handleMessage(msg *proto.Message, srv proto.Chat_StreamServer) {
	uc.eventChan <- &client.Event{
		Client: srv,
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

func (uc *_server) handleJoin(msg *proto.Message, srv proto.Chat_StreamServer) {
	uc.eventChan <- &client.Event{
		Client: srv,
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
}*/

func (uc *_server) Subscribe(userId string, conn proto.Chat_StreamServer) (client.Client, errorslib.ErrorModel) {
	c := client.New(nil, conn, userId, uc.eventChan)
	uc.clients[c.GetKey()] = c
	c.Start()

	return c, nil
}

func keyExistsInArray(key string, array []string) bool {
	for _, elm := range array {
		if elm == key {
			return true
		}
	}
	return false
}

func (uc *_server) Run(ctx context.Context) errorslib.ErrorModel {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	srv := grpc.NewServer()
	proto.RegisterChatServer(srv, uc)

	l, err := net.Listen("tcp", "0.0.0.0:8060")
	if err != nil {
		return errorslib.Internal(err).WithMessage("server unable to bind on provided host")
	}

	log.Println("server is listening on 8060")

	go func() {
		_ = srv.Serve(l)
		cancel()
	}()

	go uc.Start()

	<-ctx.Done()
	return nil
}

func (uc *_server) extractToken(ctx context.Context) (tkn string, ok bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md[userHeader]) == 0 {
		return "", false
	}

	return md[userHeader][0], true
}
