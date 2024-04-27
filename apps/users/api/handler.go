package api

import (
	"fmt"

	"github.com/Clay294/forum/protocol/middleware"

	"github.com/Clay294/forum/apps/users"
	"github.com/Clay294/forum/ioc"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type apiHandler struct {
	service users.Service
}

func (ah *apiHandler) Init() error {
	res := ioc.Controllers().GetServiceController(ah.Name())

	switch exact := res.(type) {
	case error:
		return fmt.Errorf("initializing the api handler of unit %s failed: %s", ah.Name(), exact)
	case users.Service:
		ah.service = exact
		return nil
	}

	return fmt.Errorf("initializing the api handler of unit %s failed: the result getting from ioc controllers is unknown type", ah.Name())
}

func (ah *apiHandler) Name() string {
	return users.UnitName
}

func (ah *apiHandler) Registry(router gin.IRouter) {
	// router.Use(middleware.GetRequestInfo)
	router.POST(users.CreateUserUrl, ah.HandleCreateUser)
	router.POST(users.LoginUrl, ah.HandleLogin)
	router.GET(users.ValidateTokenUrl, ah.HandleValidateToken)
	router.Use(middleware.Authenticate)
	router.GET(users.GetUserInfoUrl, ah.HandleGetUserInfo)
	router.PUT(users.UpdateUserInfoUrl, ah.HandleUpdateUserInfo)
	router.GET(users.LogoutUrl, ah.HandleLogout)
}

func init() {
	ah := new(apiHandler)

	err := ioc.Handlers().Registry(ah)

	if err != nil {
		log.Panic().Str("unit", ah.Name()).Caller().Msgf("registering api handler failed: %s", err)
	}
}
