package main

import (
	"context"
	"fmt"

	"github.com/Clay294/forum/grpc/thread"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Name   string
	Secret string
}

func (c *Config) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"client_name": c.Name, "client_secret": c.Secret}, nil
}

func (c *Config) RequireTransportSecurity() bool {
	return false
}

func NewConifg(name, secret string) *Config {
	return &Config{
		Name:   name,
		Secret: secret,
	}
}

func main() {
	do, err := grpc.Dial(
		"127.0.0.1:4044",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(NewConifg("Barry", "Barry")),
	)
	if err != nil {
		panic(err)
	}
	defer do.Close()

	threadClient := thread.NewThreadRpcClient(do)
	res, err := threadClient.CreateThread(context.Background(), &thread.ReqCreateThread{
		Title: "thread1",
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	// tuc, err := threadClient.UploadThread(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fd, err := os.Open("rsakeys\\rsa_key_pairs_uuids.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// defer fd.Close()

	// reader := bufio.NewReader(fd)

	// for {
	// 	line, err := reader.ReadBytes('\n')
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			tuc.Send(&thread.ReqUploadThread{
	// 				Meta: map[string]string{"FileName": "uuids"},
	// 				Data: line,
	// 			})
	// 			break
	// 		}
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	tuc.Send(&thread.ReqUploadThread{
	// 		Meta: map[string]string{"FileName": "uuids"},
	// 		Data: line,
	// 	})
	// }

	// rut := new(thread.ResUploadThread)
	// rut, err = tuc.CloseAndRecv()
	// if err != nil {
	// 	if err == io.EOF {
	// 		return
	// 	}
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(rut)
}
