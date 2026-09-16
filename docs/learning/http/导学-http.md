# 导学：HTTP / API

> 本课是后端链路的第一课。目标不是背 Gin API，而是能从一个真实请求出发，解释它如何经过网络端口、HTTP Server、路由、中间件、Handler、Service，最后形成响应。

## 1. 前置知识（面试高频标注）

| 知识点 | 为何需要 | 在本项目中的位置 | 高频度 |
|---|---|---|---|
| HTTP 请求与响应 | 理解客户端和服务端交换了什么 | `internal/handler`、`api/v1/v1.go` | 高 |
| URL、Query、Path、JSON Body | 理解参数从哪里进入程序 | `internal/handler/user.go`、各业务 handler | 高 |
| HTTP 方法和状态码 | 区分查询、创建、修改以及错误语义 | `internal/server/http.go`、`net/http` | 高 |
| TCP 端口与监听地址 | 理解请求如何到达 Go 进程 | `config/local.yml`、`pkg/server/http/http.go` | 高 |
| 路由匹配 | 理解 URL 如何找到处理函数 | `internal/server/http.go` | 高 |
| Middleware | 理解跨域、日志、Token 如何复用 | `internal/middleware` | 高 |
| JSON 序列化 | 理解 Go struct 如何变成 API 响应 | `api/v1/v1.go`、`api/v1/user.go` | 高 |
| Context | 在请求链路中传递取消信号和用户身份 | `ctx.Request.Context()`、`internal/base/user_holder` | 中 |
| 反向代理 | 理解 Nginx 与后端服务的关系 | `deploy/docker-compose/nginx.conf` | 中 |

## 2. 重点亮点与学习顺序（先看这个）

1. **从端口到路由：理解请求入口**
   - 重要原因：这是所有后端问题的起点，端口不通、路径不匹配、方法不对都会在进入业务代码前失败。
   - 通用技术关键词：TCP 监听、HTTP Server、URL、Router、Reverse Proxy。
   - 先看：`config/local.yml`、`cmd/server/main.go`、`pkg/server/http/http.go`、`internal/server/http.go`。
   - 顺序：配置监听地址 → 启动 Server → 创建 Gin Engine → 注册路由 → 用 `curl` 验证。

2. **从 Middleware 到 Handler：理解横切逻辑与业务入口**
   - 重要原因：日志、CORS、Token 解析等逻辑不应复制到每个接口中。
   - 通用技术关键词：Middleware Chain、Request Context、Authentication、Abort、Next。
   - 先看：`internal/middleware/cors.go`、`internal/middleware/log.go`、`internal/middleware/refresh_token.go`、`internal/handler/user.go`。
   - 顺序：先理解 `Use`，再理解 `Next`/`Abort`，最后跟踪登录接口。

3. **参数绑定与统一响应：理解 API 契约**
   - 重要原因：前后端联调的核心是输入输出契约，而不是某个框架函数名。
   - 通用技术关键词：Query Binding、JSON Binding、Validation、HTTP Status、Response Envelope。
   - 先看：`api/v1/user.go`、`api/v1/v1.go`、`internal/handler/user.go`。
   - 顺序：请求 DTO → 绑定失败 → Service 调用 → 成功响应/错误响应。

4. **HTTP 请求上下文与优雅停止**
   - 重要原因：后端不能只会启动，还要能在停止时处理已有请求并释放资源。
   - 通用技术关键词：Context Cancellation、Graceful Shutdown、Timeout、Signal。
   - 先看：`pkg/app/app.go`、`pkg/server/server.go`、`pkg/server/http/http.go`。
   - 顺序：进程信号 → 取消 Context → `Shutdown` → 五秒超时兜底。

5. **Nginx 反向代理与部署边界**
   - 重要原因：浏览器访问的地址不一定就是 Go 服务监听的地址。
   - 通用技术关键词：Reverse Proxy、Path Rewrite、Upstream、Keep-Alive。
   - 先看：`deploy/docker-compose/nginx.conf`、`deploy/docker-compose/docker-compose.yml`。
   - 顺序：客户端地址 → Nginx 监听端口 → `/api` 重写 → 后端端口。

## 3. 必备知识点

- [ ] 一个 HTTP 请求至少包含方法、URL、请求头、请求体。
- [ ] Query 参数、Path 参数和 JSON Body 的位置与用途不同。
- [ ] `GET` 通常用于读取，`POST` 通常用于创建或触发动作，`PUT` 通常用于更新。
- [ ] `2xx` 表示成功，`4xx` 通常表示请求方问题，`5xx` 通常表示服务端问题。
- [ ] 监听 `127.0.0.1:8081` 与监听 `0.0.0.0:8081` 的可访问范围不同。
- [ ] 路由只负责匹配请求，Handler 负责 HTTP 适配，Service 负责业务用例。
- [ ] 全局 Middleware 会包住后续路由；`Next` 放行，`Abort` 终止当前链路。
- [ ] 请求 Body 被读取后可能无法再次读取，因此日志中间件需要恢复 Body。
- [ ] `ctx.Request.Context()` 是 Go HTTP 请求上下文；Gin Context 是 Gin 对请求的封装。
- [ ] Nginx 的路径重写可能改变后端实际收到的 URL。
- [ ] 优雅停止需要给正在执行的请求一个有限的完成时间。

## 4. 推荐阅读（结合仓库）

| 主题 | 通用技术点 | 建议阅读位置 | 预计时间 | 读完能回答什么 |
|---|---|---|---:|---|
| 服务启动 | 配置加载、依赖注入、进程入口 | `cmd/server/main.go`、`cmd/server/wire/wire.go` | 20 分钟 | 程序从哪里启动，配置如何进入 Server？ |
| 监听 HTTP 端口 | `net/http.Server`、Listen、Shutdown | `pkg/server/http/http.go` | 20 分钟 | Go 进程如何开始接收请求？ |
| 路由注册 | Router、Group、HTTP Method | `internal/server/http.go` | 25 分钟 | `POST /user/login` 如何找到 `UserHandler.Login`？ |
| 请求参数 | Query/JSON/URI Binding | `internal/handler/user.go`、`api/v1/user.go` | 25 分钟 | 验证码和登录参数分别从哪里读取？ |
| 统一响应 | JSON 序列化、状态码、响应契约 | `api/v1/v1.go` | 15 分钟 | 为什么成功和失败都返回统一结构？ |
| 日志中间件 | Middleware Chain、请求/响应观测 | `internal/middleware/log.go` | 25 分钟 | 如何记录请求参数和响应耗时？Body 被读取后为什么要恢复？ |
| 跨域处理 | CORS、OPTIONS 预检请求 | `internal/middleware/cors.go` | 15 分钟 | 浏览器跨域请求为什么会先发 OPTIONS？ |
| Token 中间件 | Bearer Token、Request Context、TTL 刷新 | `internal/middleware/refresh_token.go`、`internal/middleware/login.go` | 30 分钟 | 登录态如何从请求头传到业务代码？ |
| 反向代理 | Rewrite、Upstream、Keep-Alive | `deploy/docker-compose/nginx.conf` | 25 分钟 | Nginx 监听 8080 时，后端为什么监听 8081？ |
| 真实验证 | curl、Swagger、日志定位 | `README.md`、`docs/swagger.yaml` | 30 分钟 | 如何证明请求确实经过了预期链路？ |

## 5. 自学提醒

先不要试图一次理解 Gin 源码。只要能沿着一个接口回答“请求从哪里来、参数在哪里、经过哪些中间件、调用哪个 Service、响应如何写回”，第一阶段就达标。若某文件或原理看不懂，请继续追问 AI；本 skill 负责给学习路径与题目，不提供逐行讲解。

## 6. 项目技术定位

- **倾向：后端**
- 依据：项目使用 Go、Gin、`net/http`、MySQL 和 Redis 实现 API 服务，并包含请求处理、认证中间件、日志、反向代理和优雅停止。

## 7. 核心原理解析

### 7.1 请求如何进入 Go 服务

- **问题**：客户端发送的 URL 如何到达某个 Go 函数？
- **机制**：客户端先连接目标 IP 和端口；HTTP Server 接收请求；Gin Engine 根据方法和路径匹配路由；匹配后执行 Handler。
- **项目落点**：`config/local.yml` 配置 `127.0.0.1:8081`，`pkg/server/http/http.go` 调用 `ListenAndServe`，`internal/server/http.go` 注册 `/user/login`。

### 7.2 Middleware 为什么可以复用横切逻辑

- **问题**：日志、跨域和 Token 校验如何避免写进每个 Handler？
- **机制**：Middleware 是包裹请求处理函数的链。它可以在 `Next` 前做前置处理，在 `Next` 后做后置处理，也可以用 `Abort` 提前结束。
- **项目落点**：`NewHTTPServer` 注册 CORS、响应日志、请求日志和 Token 刷新；需要登录的路由额外挂载 `middleware.Login()`。

### 7.3 Handler 与 Service 为什么分层

- **问题**：为什么不在路由函数里直接查询数据库？
- **机制**：Handler 处理 HTTP 细节，例如绑定参数、选择状态码和写 JSON；Service 表达业务用例，便于复用、测试和替换传输层。
- **项目落点**：`UserHandler.Login` 绑定 `LoginReq` 后调用 `UserService.Login`，Service 再访问 Redis 和数据库。

### 7.4 请求上下文如何贯穿链路

- **问题**：登录用户和取消信号如何传给后续业务？
- **机制**：Go 的 `context.Context` 用于携带请求范围的数据、截止时间和取消信号；中间件可以创建带用户信息的新 Context，并替换 Request。
- **项目落点**：`RefreshToken` 从 Redis 读取用户后，通过 `user_holder.WithUser` 写入 Request Context；`Login` 中间件再从 Context 判断是否放行。

### 7.5 反向代理改变了什么

- **问题**：为什么外部访问地址和 Go 的监听地址可以不同？
- **机制**：Nginx 作为入口接收请求，再通过 upstream 转发给后端；rewrite 可能移除外部路径前缀。
- **项目落点**：Nginx 监听 `8080`，`/api` 请求被重写后转发到 `host.docker.internal:8081` 或 `8082`。当前本地配置只明确启动了一个 Go 端口，因此多 upstream 是否可用需要运行验证。

## 8. 关键设计决策

| 决策点 | 备选 | 取舍 | 风险 | 验证 |
|---|---|---|---|---|
| 使用 Gin 作为路由层 | 直接使用 `net/http` | Gin 提供路由分组、绑定和 Middleware，学习成本较低 | 过度依赖框架，忽略 HTTP 原理 | 对比 `gin.Context` 与 `http.Handler`，用 curl 验证路由 |
| 使用统一响应结构 | 每个接口自由返回 JSON | 客户端处理方式一致 | 错误语义可能被过度简化 | 检查成功、参数错误、未登录三类响应 |
| 全局挂载日志中间件 | 在每个 Handler 手动打日志 | 覆盖面一致，减少重复代码 | 记录敏感参数，响应体过大 | 使用测试 Token 和假数据检查日志内容 |
| Token 放 Redis | 无状态 JWT | 可以主动过期、续期和保存用户信息 | Redis 故障会影响登录态；需要控制 TTL | `redis-cli TTL`、过期后请求 `/user/me` |
| Nginx 做反向代理 | 客户端直接访问 Go 端口 | 统一入口，后续可扩容和做负载均衡 | rewrite、端口和 upstream 配置错误 | 分别访问 `8080` 与 `8081`，比较日志和状态码 |
| 优雅停止设置五秒超时 | 立即退出或无限等待 | 尽量完成请求，同时避免进程无法退出 | 长请求可能被强制终止 | 发送慢请求后发送 SIGTERM，观察日志 |

## 9. 量化与验证（含待测）

建议按以下顺序做一次最小实验：

1. 启动依赖和服务：`docker compose -f deploy/docker-compose/docker-compose.yml up -d`，再运行 `go run ./cmd/server -conf config/local.yml`。
2. 验证服务入口：`curl -i http://127.0.0.1:8081/`，记录状态码、响应体和服务日志。
3. 验证参数绑定：向 `/user/login` 发送缺少字段的请求，观察 `400` 或当前实际错误响应。
4. 验证登录链路：先调用 `/user/code?phone=...`，再用测试验证码 `123456` 调用登录接口，检查 Redis 中的登录 Hash 和 TTL。
5. 验证认证链路：携带 `Authorization: Bearer <token>` 请求 `/user/me`，再删除或等待 Token 过期后重复请求。
6. 验证代理链路：访问 `http://127.0.0.1:8080/api/...`，确认 Nginx rewrite 后后端是否能正确匹配路由。

建议记录这些指标：

- 单接口响应状态码分布和错误原因（待测）。
- 请求耗时，至少记录平均值和 P95（待测）。
- Nginx 到 Go 的转发成功率（待测）。
- Token 过期前后认证结果（待测）。
- 优雅停止时已完成请求数与被取消请求数（待测）。

注意：当前 `nginx.conf` 中包含 `8081` 和 `8082` 两个 upstream 地址，但 `docker-compose.yml` 没有启动 Go 应用容器，且 `config/local.yml` 只配置了 `8081`。这是本课需要通过运行验证的配置事实，不能直接假设双实例已经可用。
