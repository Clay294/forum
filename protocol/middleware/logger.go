package middleware

import (
	"github.com/Clay294/forum/protocol"
	"github.com/gin-gonic/gin"
)

type RequestInfo struct {
	Url    string `json:"url"`
	Method string `json:"method"`
}

func GetRequestInfo(gctx *gin.Context) {
	reqInfo := &RequestInfo{
		Url:    gctx.Request.RequestURI,
		Method: gctx.Request.Method,
	}

	gctx.Set(protocol.RequestInfoContextName, reqInfo)
}
