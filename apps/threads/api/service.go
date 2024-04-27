package api

import (
	"fmt"
	"net/http"

	"github.com/Clay294/forum/apps/threads"
	"github.com/Clay294/forum/apps/users"
	"github.com/Clay294/forum/exception"
	"github.com/Clay294/forum/flog"
	"github.com/gin-gonic/gin"
)

func (ah *apiHandler) HandleCreateThread(gctx *gin.Context) {
	res, exists := gctx.Get(users.TokenMetaContextName)
	if !exists {
		flog.Flogger().Error().Str("unit", ah.Name()).Str("request_url", threads.CreateThreadUrl).Str("request_method", threads.CreateThreadMethod).Msgf("the token meta is not in context")
		exception.ResponseForumException(gctx, exception.NewErrInternalServerError(fmt.Errorf("创建失败")))
		return
	}

	tk := res.(*users.TokenMeta)

	rct := threads.NewReqCreateThread()

	err := gctx.ShouldBindJSON(rct)
	if err != nil {
		flog.Flogger().Error().Str("unit", ah.Name()).Str("request_url", threads.CreateThreadUrl).Str("request_method", threads.CreateThreadMethod).Msgf("parsing the request failed: %s", err)
		exception.ResponseForumException(gctx, exception.NewErrInvalidRequest(fmt.Errorf("解析请求参数失败")))
		return
	}

	rct.UserId, rct.UserName = tk.Subject, tk.UserName

	err = ah.service.CreateThread(gctx.Request.Context(), rct)
	if err != nil {
		exception.ResponseForumException(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, "创建成功")
}

func (ah *apiHandler) HandleSearchByMainHome(gctx *gin.Context) {
	// TODO api逻辑
}

// TODO HandleSearchBySection
