package main

import (
	"context"
	"flag"
	"github.com/aliworkshop/grpc_chat/server"
	"github.com/aliworkshop/grpc_chat/user"
	"log"
	"math/rand"
	"time"
)

var (
	serverMode bool
	debugMode  bool
	userId     string
)

func init() {
	flag.BoolVar(&serverMode, "s", false, "run as the server")
	flag.BoolVar(&debugMode, "v", false, "enable debug logging")
	flag.StringVar(&userId, "n", "", "the username for the client")
	flag.Parse()
}

func init() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)
}

func main() {
	if serverMode {
		s := server.NewServer()
		if err := s.Run(context.Background()); err != nil {
			panic(err)
		}
	} else {
		c := user.Client(userId)
		if err := c.Run(context.Background()); err != nil {
			panic(err)
		}
	}
}
