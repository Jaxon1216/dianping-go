# Go 点评项目学习资源

## Knowledge

- [本项目 C01：启动、分层与依赖组装](./interview-prep/dianping-go/guide-fullstack/C01-启动分层与依赖组装.md)
  本课主线。用于建立启动顺序、依赖方向和生命周期边界。
- [本项目生成后的 Wire 代码](./cmd/server/wire/wire_gen.go)
  项目运行时真正执行的普通 Go 构造调用。判断初始化顺序时优先看这里。
- [Google Wire README](https://github.com/google/wire/blob/main/README.md)
  Wire 的官方说明；确认它是无运行时状态、无反射的代码生成工具。项目已停止维护，因此只把它当作本仓库既有方案来学习。
- [Gin Quickstart](https://gin-gonic.com/en/docs/quickstart/)
  Gin 官方最小示例。用于对照 `gin.Default`、路由注册和启动监听。
- [Go `net/http.Server.Shutdown`](https://pkg.go.dev/net/http#Server.Shutdown)
  标准库优雅停机契约。用于区分停止接收新连接与等待活动请求完成。
- [Go 内置函数 `panic`](https://pkg.go.dev/builtin#panic)
  官方的 panic、defer 展开与非零退出说明。用于区分普通 error 与不可恢复终止。
- [Go `os/signal`](https://pkg.go.dev/os/signal)
  标准库进程信号说明。用于理解 `SIGINT`、`SIGTERM` 与退出流程。

## Wisdom (Communities)

- [Go Forum](https://forum.golangbridge.org/)
  面向 Go 实践问题的社区。用于验证工程取舍和排查非项目特有的问题。
- [Gin GitHub Discussions](https://github.com/gin-gonic/gin/discussions)
  Gin 官方仓库讨论区。用于查找路由、中间件和运行行为的真实案例。
