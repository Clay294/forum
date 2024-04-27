package middleware

import (
	"fmt"

	"github.com/Clay294/forum/apps/users"
	"github.com/Clay294/forum/exception"
	"github.com/Clay294/forum/flog"
	"github.com/Clay294/forum/ioc"
	"github.com/gin-gonic/gin"
)

func Authenticate(gctx *gin.Context) {
	jtRaw, err := gctx.Cookie(users.JsonTokenCookieName)

	if err != nil {
		flog.Flogger().Error().Str("unit", "middleware").Str("method", "ahthenticate").Msgf("the json token is not in cookie: %s", err)
		exception.ResponseForumException(gctx, exception.NewErrUnauthorized(fmt.Errorf("请先登录")))
		return
	}

	rvt := users.NewReqValidateToken()

	rvt.JsonToken = jtRaw
	rvt.ClientIP = gctx.ClientIP()

	res := ioc.Controllers().GetServiceController(users.UnitName)
	sc := res.(users.Service)

	tk, err := sc.ValidateToken(gctx.Request.Context(), rvt)
	if err != nil {
		exception.ResponseForumException(gctx, err)
		return
	}

	gctx.Set(users.TokenMetaContextName, tk)

	//var tokenMeta *users.TokenMeta
	//
	//switch exact := res.(type) {
	//case error:
	//	flog.Flogger().Error().Msgf("middleware authenticating, getting the service controller of unit %s failed: %s", users.UnitName, err)
	//	exception.ResponseForumException(gctx, exception.NewErrMiddlewareAuthenticate(fmt.Errorf("请先登录")))
	//	return
	//case users.Service:
	//	tokenMeta, err = exact.ValidateToken(gctx.Request.Context(), rvt)
	//
	//	if err != nil {
	//		exception.ResponseForumException(gctx, err)
	//		return
	//	}
	//
	//	gctx.Set(users.TokenMetaContextName, tokenMeta)
	//	return
	//}
}
