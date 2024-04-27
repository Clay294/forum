package protocol

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Clay294/forum/config"
	"github.com/Clay294/forum/flog"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	*http.Server
}

func NewHttpServer(engine *gin.Engine) *HttpServer {
	return &HttpServer{
		Server: &http.Server{
			Addr:         config.GlobalConf().CreateAddr(),
			Handler:      engine,
			ReadTimeout:  time.Second * time.Duration(config.GlobalConf().ReadTimeout),
			WriteTimeout: time.Second * time.Duration(config.GlobalConf().WriteTimeout),
		},
	}
}

func (hs *HttpServer) Start() error {
	err := hs.ListenAndServe()
	if err != nil {
		flog.Flogger().Error().Msgf("starting http server failed:%s", err)
		return fmt.Errorf("starting http server failed")
	}

	return nil
}
