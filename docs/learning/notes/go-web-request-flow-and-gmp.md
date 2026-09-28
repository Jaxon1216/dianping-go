# 前端请求的全链路映射与 Go GMP 调度模型

## 结论标题：Express 到 Gin 的洋葱模型映射与 Go 并发底层

### 问题/场景
在学习 Go 后端开发时，前端开发者常需要将 Node.js/Express 的请求链路概念迁移到 Go/Gin 项目中。此外，Go “来一个请求开一个协程”的高并发表现背后，需要理解其底层调度机制（GMP 模型）。

### 结论
1. **全链路映射（洋葱模型）**：
   - **Router**：对应 Express 的路由挂载，Gin 将请求分发到 Handler。
   - **Middleware**：前置拦截器（如鉴权、跨域），可全局挂载或局部挂载。
   - **Handler**：控制层，不写核心业务，负责 `ShouldBind` 接参，调 Service，`JSON` 返参。
   - **Service**：业务逻辑层，处理数据规则。
   - **Data/Query**：通过 GORM Gen 或 Redis 客户端与存储层交互。
   
2. **Go 并发与 GMP 调度模型**：
   - **G (Goroutine)**：轻量级协程（~2KB），处理请求的业务代码，生命周期随 Handler 返回结束。
   - **M (Machine)**：系统线程，真正执行计算的内核实体。
   - **P (Processor)**：逻辑处理器，维护 G 的本地队列，M 必须绑定 P 才能运行。
   - **调度核心**：
     - **工作窃取 (Work Stealing)**：当 P 的本地队列空时，M 会从全局队列或其他 P 的队列窃取 G，实现负载均衡。
     - **移交机制 (Handoff)**：当 M 发生系统调用阻塞时，调度器会将 P 剥离并分配给空闲的 M，防止队列阻塞。
     - **网络 I/O (Netpoller)**：底层使用 epoll/kqueue 异步网络 I/O，网络等待只挂起 G 不阻塞 M。

### 项目映射
- **Router/Middleware**：`internal/server/http.go`
- **Handler 接参与返参**：`internal/handler/user.go` (`ctx.ShouldBind`, `v1.HandleSuccess`)
- **Service 业务逻辑**：`internal/service/user.go`

### 示例或验证
```go
// 典型的 Handler 写法：接参 -> 调 Service -> 返参
func (h *UserHandler) Login(ctx *gin.Context) {
    var req v1.LoginReq
    if err := ctx.ShouldBind(&req); err != nil {
        v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
        return // Go 特色：错误立刻 return，没有 try/catch
    }
    data, err := h.userService.Login(ctx.Request.Context(), &req)
    if err != nil {
        v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
        return
    }
    v1.HandleSuccess(ctx, data)
}
```

### 易错点
1. **上下文对象混淆**：Express 中 `req/res` 分离，而在 Gin 中参数绑定和返回都在统一的 `ctx *gin.Context` 里。
2. **异常处理机制**：习惯 Express 的 `async/await + try/catch` 或全局 `next(err)`，而 Go 依赖于在每一层调用后进行 `if err != nil` 判断。
3. **阻塞理解误区**：认为 Go 写同步代码会阻塞服务器，实际上 Go 底层 GMP 模型和 Netpoller 自动实现了异步非阻塞执行，开发者只需按同步方式编写代码。

### 关联文件
- `docs/learning/notes.md` (已更新索引)
