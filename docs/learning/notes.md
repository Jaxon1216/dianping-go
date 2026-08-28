# 主题笔记

笔记只记录已确认、未来可复用的结论。一次性排错过程留在对话中，稳定结论再追加。

## 索引

- **Go**：package、接口、错误、context 和项目分层。
- **Gin/HTTP**：路由、中间件、请求和响应。
- **MySQL/GORM**：SQL、GORM Gen、模型和查询。
- **Redis**：数据结构、缓存、Lua 和 Stream。
- **并发**：goroutine、channel、锁和一致性。
- **Docker/工具**：Compose、客户端和验证。
- **项目疑点**：需要继续验证的实现状态。

## Go：package 与项目分层

### 结论

Go 按 package 组织代码，目录通常对应 package；本项目把启动、HTTP、handler、service、base、model 和 query 分开。

### 项目映射

从 `cmd/server/main.go` 跟到 `cmd/server/wire/wire.go`，再进入 `internal/server`。

### 易错点

目录名不是完整架构证明，仍需看 import 和调用关系。

## Gin/HTTP：请求链路

### 结论

路由在 `internal/server/http.go` 注册，公共 middleware 先处理请求，需要登录的路由再使用 `middleware.Login()`，最终进入对应 handler。

### 项目映射

可沿 `POST /user/login` 阅读 `NewHTTPServer`、`UserHandler.Login` 和 `UserService.Login`。

### 易错点

Gin context 与 `ctx.Request.Context()` 不是同一个概念；本项目把用户写入 request context。

## MySQL/GORM：生成模型和查询

### 结论

`internal/model` 保存 GORM Gen 生成的模型，`internal/query` 提供类型安全的查询对象，数据库结构来源于 `deploy/docker-compose/hmdp.sql`。

### 项目映射

`service.NewDB` 创建连接，`service.NewQuery` 调用 `query.Use(db)`；业务通过 `s.query.Shop` 等对象查询。

### 易错点

生成文件不应手工长期维护；重新生成后要检查 diff。

## Redis：数据结构按业务选择

### 结论

String 适合验证码、库存和缓存；Hash 适合用户字段；Set 适合关注关系；ZSet 适合点赞和 Feed；Bitmap 适合签到；GEO 适合距离查询；Stream 适合异步订单。

### 项目映射

key 常量在 `internal/base/constants/redis.go`，操作分布在 `internal/service/`。

### 易错点

Redis 写成功不代表 MySQL 已同步，需要明确数据源、过期策略和失败补偿。

## 并发：秒杀是多个一致性问题的组合

### 结论

Lua 负责在 Redis 内原子判断库存和重复购买，Stream 将订单处理异步化，分布式锁限制同一用户的订单处理并发。

### 项目映射

`internal/service/voucher_order.go`、`internal/scripts/seckill.lua`、`internal/base/redis_worker/redis_worker.go`。

### 易错点

代码存在不等于消息组、消费恢复和服务关闭流程完整，这些需要单独验证。

## Docker/工具：Compose 是本地依赖入口

### 结论

应用依赖配置中的 MySQL 和 Redis 地址；Compose 暴露 3306、6379 和 8080，并通过 SQL 文件初始化数据库。

### 项目映射

`config/local.yml`、`deploy/docker-compose/docker-compose.yml`。

### 易错点

数据库初始化脚本通常只在首次创建数据目录时执行。

## 项目疑点：代码状态需要证据

### 结论

不能只看目录或接口是否存在；需要同时检查路由、Wire 注入、service 实现、外部依赖初始化和测试/运行证据。

### 项目映射

评论 service 当前是占位；HLL 位于实验测试；秒杀 Stream 和部分签到流程仍需验证。

### 易错点

不要把“有模型/handler”误判为“业务可用”。

## 新增笔记模板

```md
## 结论标题
### 问题/场景
### 结论
### 项目映射
### 示例或验证
### 易错点
### 关联文件
```
