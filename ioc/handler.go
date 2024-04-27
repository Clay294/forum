package ioc

import (
	"fmt"
	"github.com/Clay294/forum/protocol"
	"net/url"

	"github.com/Clay294/forum/flog"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Controller
	Registry(gin.IRouter)
}

type HANDLERSCONTAINER map[string]Handler

var handlers = make(HANDLERSCONTAINER, 64)

func Handlers() HANDLERSCONTAINER {
	return handlers
}

func (hc HANDLERSCONTAINER) Registry(ah Handler) error {
	if _, ok := hc[ah.Name()]; ok {
		flog.Flogger().Error().Msgf("the api handler of %s already exists", ah.Name())
		return fmt.Errorf("the api handler of unit %s already exists", ah.Name())
	}

	hc[ah.Name()] = ah
	return nil
}

func (hc HANDLERSCONTAINER) Init(router gin.IRouter) error {
	for ahName, ah := range hc {
		err := ah.Init()
		if err != nil {
			flog.Flogger().Error().Msgf("initializing the api handler of unit %s failed: %s", ahName, err)
			return fmt.Errorf("initializing the api handler of unit %s failed: %s", ahName, err)
		}

		unitUrl, err := url.JoinPath(protocol.ApiUrl, ahName)
		if err != nil {
			flog.Flogger().Error().Msgf("initializing the api handler of unit %s failed:%s", ahName, err)
			return fmt.Errorf("initializing the api handler of unit %s failed: %s", ahName, err)
		}

		ah.Registry(router.Group(unitUrl))
	}
	return nil
}
