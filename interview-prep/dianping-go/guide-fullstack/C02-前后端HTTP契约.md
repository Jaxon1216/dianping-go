# C02 前后端 HTTP 契约

> 课程性质：核心
> 返回：[主导学](../导学.md)；前置：[C01](C01-启动分层与依赖组装.md)、HTTP 请求/响应基础；下一课：[C03 Redis 会话](C03-Redis会话与滑动续期.md)
> 对应：B1；[Q1](../面经.md#q1-项目中的一次请求如何从浏览器走到数据层)；目标掌握：看懂链路、理解机制、设计与表达。

## 本课重点

场景：用户在商户列表页发起 `GET /api/shop/of/type?typeId=1&current=1`。

- `/api` 是浏览器到代理的路径前缀；Vite/Nginx 去掉前缀后转发给 Gin。
- API 模块定义路径和参数，Axios 拦截器处理凭据与响应解包，页面维护交互状态。
- 前端路由守卫改善导航体验，服务端鉴权决定请求能否访问受保护资源。

## 主链与职责

```text
页面 loadMore
  -> frontend API 函数
  -> Axios request interceptor：可选附加 Bearer token
  -> /api/shop/of/type
  -> Vite proxy 或 Nginx rewrite 去掉 /api
  -> Gin 全局 middleware
  -> 路由匹配 ShopHandler
  -> ShouldBindQuery
  -> ShopService
  -> Redis / MySQL
  -> v1.Response JSON
  -> Axios response interceptor 取 data.data
  -> Vue 更新响应式状态并渲染
```

这条链路里有三个不同契约：

- 浏览器和代理之间：`/api` 是外部统一前缀。
- 代理和 Gin 之间：转发前去掉 `/api`，所以 Gin 注册裸业务路径。
- Gin 和页面之间：统一响应 `{ success, errorMsg?, data?, total? }`，Axios 再解包 `data`。

## 源码阅读

1. [`frontend/src/views/shop-list/index.vue`](../../../frontend/src/views/shop-list/index.vue)：`loadMore` 维护页码并消费数组。
2. [`frontend/src/api/shop.js`](../../../frontend/src/api/shop.js)：API 模块只表达路径和参数，不持有 token。
3. [`frontend/src/utils/request.js`](../../../frontend/src/utils/request.js)：请求拦截器附加 token，响应拦截器检查 `success` 并返回业务 `data`。
4. [`frontend/vite.config.js`](../../../frontend/vite.config.js)：开发环境把 `/api` 转发到 `127.0.0.1:8081` 并去前缀。
5. [`deploy/docker-compose/nginx.conf`](../../../deploy/docker-compose/nginx.conf)：部署入口做相同 rewrite，并把请求转发到两个宿主机 Go 端口。
6. [`internal/server/http.go`](../../../internal/server/http.go)：全局中间件顺序、公开路由和登录路由边界。
7. [`internal/handler/shop.go`](../../../internal/handler/shop.go)：参数绑定、400/500 映射与统一响应。
8. [`api/v1/v1.go`](../../../api/v1/v1.go)：响应 envelope 的真实字段。

## 机制与边界

### 前端路由守卫不等于服务端鉴权

Vue Router 的 `requireAuth` 只阻止正常 UI 导航。用户仍可用 curl、浏览器控制台或其他客户端直接请求 API。因此：

- 路由守卫负责体验：无 token 时提前跳登录页。
- Axios 拦截器负责传输：有 token 时附加请求头，401 时清理本地状态。
- Gin 登录中间件负责安全边界：服务端依据恢复出的用户决定是否放行。

三层职责不能互相替代。

### Gin Context 与 request Context

全局中间件接收 `*gin.Context`，但项目把用户写进 `ctx.Request.Context()`。handler 随后把这个标准库 context 传给 service。这样业务层不必依赖 Gin，也能获得取消信号和值。

注意：把值写入新 context 后，必须把新 request 重新赋给 `ctx.Request`；`context.WithValue` 不会原地修改旧 context。

### 错误契约现状

- 参数绑定失败通常是 400，并返回统一 JSON。
- 多数业务错误被 handler 统一映射为 500，即使“验证码错误”“库存不足”属于可预期业务结果。
- 登录中间件直接返回空的 401；Axios 通过 HTTP error 分支处理，不依赖统一 JSON。
- `Response.Total` 带 `omitempty`，总数为 0 时字段会消失；列表接口中的 `total` 目前常是当前页长度，不是数据库总量。

### 已确认的契约风险

- Swagger 把 `BasePath` 设为 `/v1`，真实 Gin 路由却没有 `/v1` 前缀；文档请求路径可能与运行接口不一致。
- Nginx upstream 配置 8081 和 8082，但 Compose 不启动 Go 服务；需要人工启动两个不同配置的进程。
- *Axios 只有 10 秒总超时；*后端访问 MySQL/Redis 时大多没有从请求层显式设置更短的依赖超时。
- 响应错误码映射较粗，客户端只能依赖错误文本，不利于稳定分支处理。

## 需要懂什么

- **三层客户端职责的变化原因不同。** 接口路径改变应修改 API 模块，token 或统一错误策略改变应修改拦截器，加载与渲染变化应由页面处理；混在一起会让跨页面策略重复。
- **本地凭据存在不等于服务端登录有效。** sessionStorage 中的 token 可能已经过期或被撤销，只有服务端恢复出有效会话，受保护路由才能放行。
- **HTTP 状态与业务结果是两个维度。** 当前接口同时存在统一 JSON、空 401 和粗粒度 500 映射；稳定业务 code 是可选演进，不能当作已有契约。
- **构建体积只是性能线索。** chunk 大小不能单独说明首屏是否慢，还取决于网络、缓存、解析执行和设备；构建通过同样不代表运行功能正确。

## 需要记住什么

- 主链：页面 → API 模块 → Axios → Vite/Nginx → Gin → handler → service → 存储 → 响应解包。
- 外部 `/api/shop/...` 对应后端 `/shop/...`；排查 404 时先分清代理前缀与 Gin 路由。
- 常见响应为 `{ success, errorMsg?, data?, total? }`，Axios 通常只返回业务 data；total 当前不总是数据库总量。
- Gin Context 管 HTTP 流程，标准 request context 跨层传值与取消；WithValue 返回新 context，必须重新装回 Request。

