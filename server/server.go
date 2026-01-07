package server

import (
	"context"
	"net/http"
	"pionex-administrative-sys/server/handler"
	"pionex-administrative-sys/static"
	"pionex-administrative-sys/utils/logger"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	srv    *http.Server
	addr   string
}

func New(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Init() {
	// 接管 Gin 内部日志输出
	gin.DefaultWriter = logger.InfoWriter()
	gin.DefaultErrorWriter = logger.ErrorWriter()
	gin.SetMode(gin.ReleaseMode)

	s.engine = gin.New()
	static.Register(s.engine)
	handler.Register(s.engine)
	s.srv = &http.Server{
		Addr:    s.addr,
		Handler: s.engine,
	}
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	_ = s.srv.Shutdown(ctx)
}
