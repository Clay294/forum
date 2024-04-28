package main

import (
	"fmt"
	"net"

	"github.com/Clay294/forum/grpc/thread"
	_ "github.com/Clay294/forum/grpc/thread/impl"
	"github.com/Clay294/forum/ioc"
	"google.golang.org/grpc"
)

var grpcServer = grpc.NewServer(grpc.ChainUnaryInterceptor(Authentication))

func main() {
	res := ioc.Controllers().GetServiceController("thread")
	fmt.Println(res)
	thread.RegisterThreadRpcServer(grpcServer, res.(thread.ThreadRpcServer))

	sc, err := net.Listen("tcp", "127.0.0.1:4044")
	if err != nil {
		panic(err)
	}

	err = grpcServer.Serve(sc)
	if err != nil {
		panic(err)
	}
}
