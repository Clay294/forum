package impl

import (
	"context"
	"time"

	"github.com/Clay294/forum/grpc/thread"
)

func (c *controller) CreateThread(ctx context.Context, rct *thread.ReqCreateThread) (*thread.Thread, error) {
	now := time.Now().UnixMilli()
	return &thread.Thread{
		ThreadBase: rct,
		ThreadMeta: &thread.ThreadMeta{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, nil
}

// func (c *controller) UploadThread(stream thread.ThreadRpc_UploadThreadServer) error {
// 	for {
// 		data, err := stream.Recv()
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			fmt.Println(err)
// 			return err
// 		}
// 		fmt.Println(data.Meta)
// 		fmt.Println(string(data.Data))
// 	}

// 	stream.SendAndClose(&thread.ResUploadThread{Message: "上传成功"})
// 	return nil
// }
