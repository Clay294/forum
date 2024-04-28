package main

import (
	"context"
	"fmt"

	"github.com/Clay294/forum/grpc/thread"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	do, err := grpc.Dial("127.0.0.1:4044", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer do.Close()

	threadClient := thread.NewThreadRpcClient(do)

	res, err := threadClient.CreateThread(context.Background(), &thread.ReqCreateThread{
		Title: "thread1",
	})

	fmt.Println(res)
	fmt.Println(err)
}
