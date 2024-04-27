package api

import (
	"fmt"

	"github.com/Clay294/forum/apps/threads"
	"github.com/Clay294/forum/flog"
	"github.com/Clay294/forum/ioc"
	"github.com/Clay294/forum/protocol/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type apiHandler struct {
	service threads.Service
}

func (ah *apiHandler) Init() error {
	res := ioc.Controllers().GetServiceController(threads.UnitName)
	switch exact := res.(type) {
	case error:
		flog.Flogger().Error().Str("unit", threads.UnitName).Bool("init", true).Msgf("getting the serivce controller of unit %s failed: %s", threads.UnitName, exact)
		return exact
	case threads.Service:
		ah.service = exact
		return nil
	}
	flog.Flogger().Error().Str("unit", threads.UnitName).Bool("init", true).Msgf("getting the serivce controller of unit %s failed: %s", threads.UnitName, "unknown type")

	return fmt.Errorf("getting the serivce controller of unit %s failed: %s", threads.UnitName, "unknown type")
}

func (ah *apiHandler) Name() string {
	return threads.UnitName
}

func (ah *apiHandler) Registry(router gin.IRouter) {
	router.Use(middleware.Authenticate)
	router.POST(threads.CreateThreadUrl, ah.HandleCreateThread)
}

func init() {
	ah := new(apiHandler)

	err := ioc.Handlers().Registry(ah)
	if err != nil {
		log.Panic().Str("unit", ah.Name()).Caller().Msgf("registering api handler failed: %s", err)
	}
}
