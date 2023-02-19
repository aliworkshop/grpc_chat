package client

import (
	"github.com/aliworkshop/grpc_chat/client/data"
	"github.com/aliworkshop/grpc_chat/proto"
	"github.com/gorilla/websocket"
	"log"
)

func (c *client) read() {
	for c.conn != nil {
		b, err := c.conn.Recv()
		if err != nil {
			log.Println("error on read message")
			switch e := err.(type) {
			case *websocket.CloseError:
				if e.Code >= websocket.CloseNormalClosure && e.Code <= websocket.CloseTLSHandshake {
					c.Stop()
				}
			default:
				log.Println(e)
			}
			return
		}
		c.handleTextMsg(b)
	}
}

func (c *client) handleTextMsg(msg *proto.Message) {
	switch msg.Action {
	case data.Message.String():
		c.handleMessage(msg)
	case data.Join.String():
		c.handleJoin(msg)
	}
}
