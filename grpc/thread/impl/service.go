package impl

import (
	"context"
	"github.com/Clay294/forum/grpc/thread"
	"time"
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
