package impl

import (
	"github.com/Clay294/forum/grpc/thread"
	"github.com/Clay294/forum/ioc"
)

var _ thread.ThreadRpcServer = &controller{}

func init() {
	c := &controller{}
	err := ioc.Controllers().Registry(c)
	if err != nil {
		panic(err)
	}
}

type controller struct {
	thread.UnimplementedThreadRpcServer
}

func (c *controller) Init() error {
	return nil
}

func (c *controller) Name() string {
	return "thread"
}
