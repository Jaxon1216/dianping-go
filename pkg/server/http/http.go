package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-dianping/pkg/log"
	"net/http"
	"time"
)

type Server struct {
	*gin.Engine
	httpSrv *http.Server
	host    string
	port    int
	logger  *log.Logger
}
type Option func(s *Server)

func NewServer(engine *gin.Engine, logger *log.Logger, opts ...Option) *Server {
	s := &Server{
		Engine: engine,
		logger: logger,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
func WithServerHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}
func WithServerPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) Start(context.Context) error {
	// 【启动阶段：真正启动 HTTP 服务】
	// App.Run 会调用这个 Start 方法。前面的路由注册只是“准备规则”，
	// 到这里才创建标准库 HTTP Server 并开始监听端口。

	// 这里创建的是 Go 标准库的 HTTP Server。
	// s.host 和 s.port 来自 config/local.yml，拼成例如 127.0.0.1:8081。
	s.httpSrv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		// Server 嵌入了 *gin.Engine，所以它可以作为 HTTP Handler 接收请求。
		Handler: s,
	}

	// ListenAndServe 来自 Go 标准库 net/http 包，不是本项目自定义的方法。
	// 它会绑定 Addr 端口并持续等待客户端请求；收到请求后交给 Handler。
	if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Sugar().Fatalf("listen: %s\n", err)
	}

	return nil
}
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Sugar().Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		s.logger.Sugar().Fatal("Server forced to shutdown: ", err)
	}

	s.logger.Sugar().Info("Server exiting")
	return nil
}
