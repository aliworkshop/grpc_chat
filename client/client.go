package client

import (
	"github.com/aliworkshop/errorslib"
	"github.com/aliworkshop/grpc_chat/client/data"
	"github.com/aliworkshop/grpc_chat/client/event"
	"github.com/aliworkshop/grpc_chat/proto"
	"github.com/aliworkshop/loggerlib/logger"
	"sync"
)

type Client interface {
	Start()
	Stop()
	IsAlive() bool

	GetKey() string

	SetTemp(key string, value interface{})
	GetTemp(key string) interface{}

	Write(w *WriteRequest) errorslib.ErrorModel
	WriteBinary(data []byte) errorslib.ErrorModel
	WriteJson(data *data.Data) errorslib.ErrorModel
}

type client struct {
	log logger.Logger

	conn    proto.Chat_StreamServer
	connMtx *sync.Mutex

	writeChan chan *WriteRequest
	eventChan chan *Event

	started bool
	// key unique key of connection
	key string

	values    map[string]interface{}
	valuesMtx *sync.RWMutex
}

func New(log logger.Logger, conn proto.Chat_StreamServer, userId string, eventChan chan *Event) Client {
	c := &client{
		log:       log,
		conn:      conn,
		connMtx:   new(sync.Mutex),
		writeChan: make(chan *WriteRequest),
		eventChan: eventChan,
		key:       userId,
		values:    make(map[string]interface{}),
		valuesMtx: new(sync.RWMutex),
	}
	return c
}

func (c *client) Start() {
	if c.started {
		return
	}
	c.started = true

	go c.read()
	go c.write()
}

func (c *client) Stop() {
	if !c.started {
		return
	}

	c.connMtx.Lock()
	defer c.connMtx.Unlock()

	close(c.writeChan)

	c.eventChan <- &Event{
		Client: c,
		Event:  event.Closed,
	}

	c.started = false
}

func (c *client) IsAlive() bool {
	return c.conn != nil
}

func (c *client) GetKey() string {
	return c.key
}
