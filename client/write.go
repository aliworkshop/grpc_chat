package client

import (
	"encoding/json"
	"github.com/aliworkshop/errorslib"
	"github.com/aliworkshop/grpc_chat/client/data"
	"github.com/aliworkshop/grpc_chat/proto"
	"github.com/aliworkshop/loggerlib/logger"
	"github.com/gorilla/websocket"
	"time"
)

type WriteRequest struct {
	Type int
	Data []byte
}

func (c *client) Write(w *WriteRequest) errorslib.ErrorModel {
	return c.writeMessage(w.Type, w.Data)
}

func (c *client) WriteBinary(data []byte) errorslib.ErrorModel {
	return c.writeMessage(websocket.BinaryMessage, data)
}

func (c *client) WriteJson(data *data.Data) errorslib.ErrorModel {
	c.connMtx.Lock()
	defer c.connMtx.Unlock()

	b, _ := json.Marshal(data.Body)
	if err := c.conn.Send(&proto.Message{
		Action: data.Action.String(),
		Type:   data.Type.String(),
		Id:     data.Id,
		Body:   b,
	}); err != nil {
		return errorslib.HandleError(err)
	}
	return nil
}

func (c *client) writeMessage(messageType int, data []byte) errorslib.ErrorModel {
	c.connMtx.Lock()
	defer c.connMtx.Unlock()

	if c.conn == nil {
		return errorslib.Internal().
			WithMessage("connection is closed")
	}

	if err := c.conn.SendMsg(data); err != nil {
		return errorslib.HandleError(err)
	}
	return nil
}

func (c *client) write() {
	t := time.NewTicker(time.Second * 5)
	defer t.Stop()

	for c.IsAlive() {
		select {
		case w := <-c.writeChan:
			if w == nil {
				return
			}
			if err := c.writeMessage(w.Type, w.Data); err != nil {
				c.log.
					With(logger.Field{
						"error": err,
					}).
					WithId("c.conn.WriteMessage").
					ErrorF("error on conn.WriteMessage")
			}
		case <-t.C:
			//err := c.conn(websocket.PingMessage, nil)
			//if err != nil {
			//	c.Stop()
			//	c.log.WithId("c.conn.PingMessage").With(
			//		logger.Field{
			//			"error": err.Error(),
			//		}).ErrorF("")
			//}
		}
	}
}
