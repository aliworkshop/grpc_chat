package user

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliworkshop/grpc_chat/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"time"
)

type ClientUc interface {
	Run(ctx context.Context) error
}

type _client struct {
	conn   proto.ChatClient
	UserId string
}

const userHeader = "x-user"

func Client(userId string) ClientUc {
	return &_client{
		UserId: userId,
	}
}

func (c *_client) Run(ctx context.Context) error {
	connCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(connCtx, "localhost:8060", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return errors.WithMessage(err, "unable to connect")
	}
	defer conn.Close()

	c.conn = proto.NewChatClient(conn)

	err = c.stream(ctx)

	return errors.WithMessage(err, "stream error")
}

func (c *_client) stream(ctx context.Context) error {
	md := metadata.New(map[string]string{userHeader: c.UserId})
	ctx = metadata.NewOutgoingContext(ctx, md)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := c.conn.Stream(ctx)
	if err != nil {
		return err
	}
	defer client.CloseSend()

	log.Println(time.Now(), "connected to stream")

	go c.send(client)
	return c.receive(client)
}

func (c *_client) receive(sc proto.Chat_StreamClient) error {
	for {
		res, err := sc.Recv()

		if s, ok := status.FromError(err); ok && s.Code() == codes.Canceled {
			log.Println("stream canceled (usually indicates shutdown)")
			return nil
		} else if err == io.EOF {
			log.Println("stream closed by server")
			return nil
		} else if err != nil {
			return err
		}
		fmt.Printf("%s: %s\n", c.UserId, string(res.Body))
	}
}

func (c *_client) send(client proto.Chat_StreamClient) {
	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanLines)

	for {
		select {
		case <-client.Context().Done():
			log.Println("client send loop disconnected")
		default:
			if sc.Scan() {
				type Message struct {
					Action string
					Type   string
					Id     string
					Body   Byte
				}
				var msg Message
				err := json.Unmarshal(sc.Bytes(), &msg)
				if err != nil {
					log.Println(time.Now(), "failed to unmarshal message: ", err)
					return
				}
				if err = client.Send(&proto.Message{
					Action: msg.Action,
					Type:   msg.Type,
					Id:     msg.Id,
					Body:   msg.Body,
				}); err != nil {
					log.Println(time.Now(), "failed to send message: ", err)
					return
				}
			} else {
				log.Println(time.Now(), "input scanner failure: %v", sc.Err())
				return
			}
		}
	}
}

type Byte []byte

func (p *Byte) UnmarshalJSON(b []byte) error {
	*p = make(Byte, len(b))
	copy(*p, b)
	return nil
}
