# 主题笔记

笔记只记录已确认、未来可复用的结论。一次性排错过程留在对话中，稳定结论再追加。

## 索引

- **Go**：package、接口、错误、context 和项目分层。
  - 专题：[Go 基础 01：指针 / 接口 / nil / error / panic](notes/go-basics-01-pointer-interface-error.md)（第 1 轮对话沉淀，STAR 格式）。
- **Gin/HTTP**：路由、中间件、请求和响应。
  - 专题：[前端请求全链路映射与 Go GMP 调度模型](notes/go-web-request-flow-and-gmp.md)（Express 映射，GMP 模型机制）。
- **MySQL/GORM**：SQL、GORM Gen、模型和查询。
- **Redis**：数据结构、缓存、Lua、Stream、动态参数和原子性。
- **并发**：goroutine、channel、锁和一致性。
- **本轮 MySQL/Redis 学习**：DDL 字段、分页、JOIN、事务回调、索引、幂等和 Lua 调用链。
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

## MySQL/GORM：事务回调与数据库对象

### 结论

`Transaction` 回调返回非 `nil` error 时，事务管理器通常回滚；返回 `nil` 时通常提交。`return err` 和 `return nil` 只是向外层报告结果，本身不直接执行 `ROLLBACK` 或 `COMMIT`。

### 项目映射

`internal/service/voucher_order.go` 的 `createVoucherOrder` 通过 `s.query.Transaction` 包裹查重、扣库存和保存订单。回调参数 `tx *query.Query` 是绑定当前事务的查询对象，事务内的查询和写入应优先使用 `tx`。

### 易错点

当前实现拿到 `tx` 后仍使用 `s.query` 访问订单和库存，是否真正加入同一事务需要运行验证，属于事务边界风险。`if err := save(); err != nil` 中的分号是 Go 的 if 初始化语句分隔符，不是 SQL 分号。

## MySQL：LIMIT/OFFSET、JOIN 与参数占位符

### 结论

`LIMIT n` 表示最多返回 n 行；`OFFSET n` 表示先跳过 n 行。`LIMIT 10 OFFSET 10` 在稳定排序下表示返回第 11 到第 20 行。`JOIN ... ON` 根据关联字段拼接多张表，`WHERE` 再过滤结果。SQL 中的 `?` 是参数占位符，由数据库驱动绑定真实值。

### 项目映射

`internal/service/shop.go` 使用 `FindByPage` 生成分页查询；`internal/service/user.go` 按手机号查询用户；博客和用户可通过 `tb_blog.user_id = tb_user.id` 进行联合查询。

### 易错点

分页最好配合 `ORDER BY`，否则结果顺序不稳定。`NULL` 不能用 `= NULL` 判断，应使用 `IS NULL`。参数占位符不能理解成需要用户输入的问号。

## MySQL：索引与秒杀幂等

### 结论

索引是辅助数据库快速定位数据的额外结构，主键和唯一索引同时具备约束作用。秒杀一人一单需要幂等设计，常见组合是 Redis Set/分布式锁减少并发，MySQL 唯一索引 `(user_id, voucher_id)` 做最终兜底。

### 项目映射

当前 DDL 中 `tb_user.phone` 有唯一索引，`tb_shop.type_id` 有普通索引，`tb_voucher_order` 主要只有 `id` 主键。秒杀代码在 Redis Lua 中判断重复购买，并在 Go 消费者中使用 RedSync 锁。

### 易错点

代码层面的“先查询再插入”不是数据库级别的最终约束；高并发下仍应考虑唯一索引和重复键错误处理。可使用 `SHOW INDEX` 和 `EXPLAIN` 验证索引是否存在、是否被使用。

## Redis：Lua 调用链、动态参数与原子性

### 结论

Go 代码通过 go-redis 的 `Script.Run` 调用 Redis，Redis 服务端执行 Lua。`Run` 的普通参数会按顺序进入 Lua 的 `ARGV[1]`、`ARGV[2]` 等位置。Redis Lua 的原子性主要表示脚本执行期间不会被其他 Redis 命令插队，不等同于 MySQL 事务的失败自动回滚。

### 项目映射

`internal/service/voucher_order.go` 启动时读取 `internal/scripts/seckill.lua`，请求到来时传入 `req.VoucherID`、`userId` 和 `orderId`。脚本检查库存、判断重复购买、扣减 Redis 库存、记录用户并写入 `stream.orders`，返回 0、1、2。

### 易错点

Redis 数据主要在内存中，所以访问延迟低；Lua 运行在 Redis 服务端，不是运行在“内存环境”里。高内聚是代码设计概念，原子性是并发执行语义，二者不是同一个概念。Redis Lua 和 MySQL 事务也不是跨系统的同一个事务。

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
