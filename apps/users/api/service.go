package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Clay294/forum/apps/users"

	"github.com/Clay294/forum/protocol"

	"github.com/Clay294/forum/exception"
	"github.com/Clay294/forum/flog"
	"github.com/gin-gonic/gin"
)

func (ah *apiHandler) HandleCreateUser(gctx *gin.Context) {
	rcu := users.NewReqCreateUser()

	err := gctx.ShouldBind(rcu)

	if err != nil {
		flog.Flogger().Error().Msgf("unit %s, parsing parameters from user creation request failed: %s", ah.Name(), err)
		exception.ResponseForumException(gctx, exception.NewErrInvalidRequest(fmt.Errorf("解析请求参数失败")))
		return
	}

	rcu.CreatedIn = gctx.ClientIP()

	err = ah.service.CreateUser(gctx.Request.Context(), rcu)

	if err != nil {
		exception.ResponseForumException(gctx, err)
		return
	}

	gctx.JSON(http.StatusCreated, "创建成功")
}

func (ah *apiHandler) HandleGetUserInfo(gctx *gin.Context) {
	rgui := new(users.ReqGetUserInfo)

	if res := gctx.Query("user_id"); res != "" {
		userId, err := strconv.Atoi(res)
		if err != nil {
			flog.Flogger().Error().Msgf("unit %s, parsing parameters from getting user information request failed: %s", ah.Name(), err)
			exception.ResponseForumException(gctx, exception.NewErrInvalidRequest(fmt.Errorf("解析请求参数失败")))
			return
		}

		rgui.UserId = userId
	}

	user, err := ah.service.GetUserInfo(gctx.Request.Context(), rgui)

	if err != nil {
		exception.ResponseForumException(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, user)
}

func (ah *apiHandler) HandleLogin(gctx *gin.Context) {
	rl := users.NewReqLogin()

	err := gctx.ShouldBind(rl)

	if err != nil {
		flog.Flogger().Error().Msgf("unit %s, parsing parameters from user login request failed: %s", ah.Name(), err)
		exception.ResponseForumException(gctx, exception.NewErrInvalidRequest(fmt.Errorf("解析请求参数失败")))
		return
	}

	rl.LoggedinIn = gctx.ClientIP()

	err = rl.SetLoginBy()

	if err != nil {
		exception.ResponseForumException(gctx, exception.NewErrLogin(fmt.Errorf("登陆失败")))
		return
	}

	token, err := ah.service.Login(gctx.Request.Context(), rl)

	if err != nil {
		exception.ResponseForumException(gctx, err)
		return
	}

	gctx.SetCookie(
		users.JsonTokenCookieName,
		token.JsonToken,
		protocol.CookieMaxAge,
		protocol.CookiePath,
		protocol.Server,
		false,
		true,
	)

	// gctx.JSON(http.StatusOK, token.Subject)
	gctx.JSON(http.StatusOK, "登录成功")
}

func (ah *apiHandler) HandleValidateToken(gctx *gin.Context) {
	rvt := users.NewReqValidateToken()

	jtRaw, err := gctx.Cookie(users.JsonTokenCookieName)

	if err != nil {
		flog.Flogger().Error().Str("unit", ah.Name()).Str("request_url", users.ValidateTokenUrl).Str("request_method", users.ValidateTokenMethod).Msgf("the json token is not in cookie:%s", err)
		gctx.JSON(
			http.StatusUnauthorized,
			map[string]interface{}{
				"is_valid": false,
				"error":    exception.NewErrUnauthorized(fmt.Errorf("请先登录")),
			},
		)
		// exception.ResponseForumException(gctx, exception.NewErrUnauthorized(fmt.Errorf("请先登录")))
		return
	}

	rvt.JsonToken = jtRaw
	rvt.ClientIP = gctx.ClientIP()

	_, err = ah.service.ValidateToken(gctx.Request.Context(), rvt)

	if err != nil {
		gctx.JSON(
			http.StatusUnauthorized,
			map[string]interface{}{
				"is_valid": false,
				"error":    err,
			},
		)
		return
	}

	gctx.JSON(
		http.StatusOK, map[string]interface{}{
			"is_valid": true,
			"error":    nil,
		},
	)
}

func (ah *apiHandler) HandleLogout(gctx *gin.Context) {
	res, ok := gctx.Get(users.TokenMetaContextName)
	if !ok {
		flog.Flogger().Error().Str("unit", ah.Name()).Str("request_url", users.LogoutUrl).Str("request_method", users.LogoutMethod).Msgf("the token meta is not in context")
		exception.ResponseForumException(gctx, exception.NewErrInternalServerError(fmt.Errorf("等出失败")))
	}

	tokenMeta := res.(*users.TokenMeta)

	rl := &users.ReqLogout{JsonTokenId: tokenMeta.ID}

	err := ah.service.Logout(gctx.Request.Context(), rl)
	if err != nil {
		exception.ResponseForumException(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, "已登出")
}

func (ah *apiHandler) HandleUpdateUserInfo(gctx *gin.Context) {
	res, ok := gctx.Get(users.TokenMetaContextName)
	if !ok {
		flog.Flogger().Error().Str("unit", ah.Name()).Str("request_url", users.UpdateUserInfoUrl).Str("request_method", users.UpdateUserInfoMethod).Msgf("the token meta is not in context")
		exception.ResponseForumException(gctx, exception.NewErrInternalServerError(fmt.Errorf("更改失败")))
		return
	}

	tokenMeta := res.(*users.TokenMeta)

	//if !ok {
	//	flog.Flogger().Error().Str("unit", ah.Name()).Str("request_url", users.UpdateUserInfoUrl).Str("request_method", users.UpdateUserInfoMethod).Msgf("the token meta is not in context")
	//	flog.Flogger().Error().Msgf("unit %s, the token meta from updating userinfo request is invalid", ah.Name())
	//	exception.ResponseForumException(gctx, exception.NewErrUpdateUserInfo(fmt.Errorf("更改失败: 令牌元数据无效")))
	//	return
	//}

	ruui := users.NewReqUpdateUserInfo(tokenMeta)

	err := gctx.ShouldBind(ruui)
	if err != nil {
		flog.Flogger().Error().Str("unit", ah.Name()).Str("request_url", users.UpdateUserInfoUrl).Str("request_method", users.UpdateUserInfoMethod).Msgf("parseing request failed")
		exception.ResponseForumException(gctx, exception.NewErrInvalidRequest(fmt.Errorf("解析请求参数失败")))
		return
	}

	ruui.SetReqUserInfoMeta(tokenMeta)

	err = ah.service.UpdateUser(gctx.Request.Context(), ruui)
	if err != nil {
		exception.ResponseForumException(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, "更新成功")
}
