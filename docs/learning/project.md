# 项目地图与模块状态

## 项目地图

```text
cmd/server/main.go
  -> pkg/config.NewConfig
  -> pkg/log.NewLog
  -> cmd/server/wire.NewWire
  -> internal/server.NewHTTPServer
  -> pkg/app.App.Run
```

HTTP 请求通常经过：

```text
路由 -> 公共 middleware -> 登录 middleware（如需要）
  -> handler -> service -> query/Redis -> api/v1 响应
```

| 目录 | 职责 | 首先阅读 |
| --- | --- | --- |
| `cmd/server` | 启动和依赖注入 | `main.go`、`wire/wire.go` |
| `internal/server` | Gin 实例、路由、中间件挂载 | `http.go` |
| `internal/handler` | HTTP 参数、状态码、响应 | `handler.go` |
| `internal/service` | 业务用例和数据访问协调 | `service.go`、目标 service |
| `internal/middleware` | 登录、token 刷新、日志、CORS | `refresh_token.go` |
| `internal/base` | 缓存、锁、ID、用户上下文 | `cache_client.go` |
| `internal/model`、`internal/query` | GORM Gen 模型和查询 | 目标表对应的生成文件 |
| `api/v1` | DTO、错误和统一响应 | `v1.go` |
| `internal/scripts` | Redis Lua 脚本 | `seckill.lua`、`unlock.lua` |
| `pkg` | 应用、配置、日志和 server 抽象 | `pkg/app/app.go` |
| `deploy/docker-compose` | MariaDB、Redis、Nginx、SQL | `docker-compose.yml` |
| `test` | 集成实验和 Redis 行为验证 | `test/app/app_test.go` |

Redis 使用位置：

- String/TTL：验证码、token、缓存、秒杀库存。
- Hash：登录用户信息。
- List：商户类型缓存。
- Set：关注关系和共同关注。
- ZSet：点赞用户、Feed、时间排序。
- Bitmap：用户签到。
- GEO：附近商户。
- HyperLogLog：`test/app/app_test.go` 中的 UV 实验。
- Stream：秒杀订单异步消费。
- Lua：库存/重复购买判断和锁安全释放。

## 学习路线

学习顺序、每阶段目标和完成标准统一维护在 [roadmap.md](roadmap.md)，本文只负责「项目长什么样」和「哪些模块真的可用」。

需要逐条展开某一课，进入 [courses/](courses/)。

## 官方模块与仓库状态

状态含义：**已实现** 表示主链路存在；**部分实现** 表示有代码但仍有接通或正确性问题；**实验代码** 表示仅用于验证；**占位** 表示结构存在但业务未完成。

| 模块 | Go 项目文件 | 状态 | 关键数据结构 | 当前疑点 |
| --- | --- | --- | --- | --- |
| 短信登录 | `internal/service/user.go`、`internal/middleware/refresh_token.go` | 部分实现 | String、Hash、TTL | 生产短信未接入；异常和 token 生命周期需验证 |
| 商户缓存 | `internal/service/shop.go`、`internal/base/cache_client/` | 部分实现 | String、TTL、锁 | 互斥锁分支存在反序列化和重试风险 |
| 用户签到 | `internal/service/user.go` | 部分实现 | Bitmap | 路由、月份和位偏移需验证 |
| UV 统计 | `test/app/app_test.go` | 实验代码 | HLL | 尚无业务接口 |
| 附近商户 | `internal/service/shop.go`、`internal/server/http.go` | 部分实现 | GEO、距离、分页 | 路由和经纬度命名需验证 |
| 达人探店/点赞 | `internal/service/blog.go` | 已实现 | ZSet、DB 计数 | 并发一致性和返回顺序需验证 |
| 关注/共同关注 | `internal/service/follow.go` | 部分实现 | Set、SQL 关系 | Wire/server 注入和字段需核对 |
| Feed | `internal/service/blog.go` | 已实现 | ZSet、时间游标 | 游标和同时间分组逻辑需验证 |
| 优惠券秒杀 | `voucher.go`、`voucher_order.go`、`seckill.lua` | 部分实现 | Lua、String、Redsync | Stream group、消息解析、落库和恢复需审阅 |
| 商户类型 | `internal/service/shop_type.go` | 已实现 | List、TTL | JSON 错误处理需验证 |
| 评论 | `internal/service/blog_comments.go`、`internal/handler/blog_comments.go` | 占位 | SQL 表 | service 只返回空模型，路由未完整接通 |
| 生成与注入 | `cmd/generate/`、`internal/model/`、`internal/query/`、`cmd/server/wire/` | 已实现 | 生成代码、provider set | 修改 provider 后需重新生成并核对 |

目录存在不代表模块完成；状态升级必须有源码、测试、日志或客户端查询证据。
