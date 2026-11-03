# Go 后端验收标准：命令 / Gin / GORM Gen / Redis / Go 基础

> 用途：把「需要学会并记住什么」变成可勾选的验收清单，避免只停留在“看过”。
> 来源：2026-11-03 对话，对照本项目真实源码。
> 验收级别：
> - **L1 能看懂**：知道这行代码在做什么。
> - **L2 能默写**：不看资料写出项目里的常见写法（当前目标）。
> - **L3 能改造/排错**：能解释为什么，能定位问题并改对。
>
> 判定方式：每条都要能同时做到「写出来 + 讲清楚为什么 + 在项目里指出出处」。

这份清单用于课后复习和验收。首次学习时，先理解一个小例子，再回到项目写法；每次只推进一个知识点。

接口、结构体和构造函数还看不懂时，先读 [小课 01：谁规定能力，谁真正干活？](../lessons/0001-user-service-interface.html)，按「TS 对象 → Go 类型与方法 → 创建并返回 → 项目依赖」逐步展开，配有即时反馈练习。复习时看 [速查页](../reference/user-service-interface.html)。

想看每一步执行动作、传入值，以及 s 如何读写实例，用 [发短信执行演示](../lessons/0002-sms-execution.html)。大小写放到执行链理解后再单独看。

---

## 0. 总原则

只背「骨架 + 项目真实用法」，不背生成代码和完整 API 手册。

```text
Gin 的 7 个固定写法 > GORM Gen 的链式查询 > Redis 的“结构对应场景” > Go 基础语法
```

---

## 1. 命令验收

### 必须能独立执行

- [ ] 启动依赖：`cd deploy/docker-compose && docker compose up -d`
- [ ] 启动服务：`go run ./cmd/server -conf config/local.yml`
- [ ] 构建：`go build -o ./bin/server ./cmd/server`
- [ ] 跑测试：`go test ./...`
- [ ] 静态检查：`go vet ./...`
- [ ] 重新生成依赖注入：`wire ./cmd/server/wire`
- [ ] 重新生成 API 文档：`swag init -g cmd/server/main.go -o ./docs`
- [ ] 重新生成 Model/Query：`go run ./cmd/generate`

### 调试命令

- [ ] `redis-cli -a 123321`
- [ ] `redis-cli -a 123321 MONITOR`（排查缓存命中/写入最有效）
- [ ] `curl -X POST "http://127.0.0.1:8081/user/code?phone=13800138000"`

### 验收问题

- [ ] 能说明 `wire`、`swag`、`cmd/generate` 分别在什么情况下需要重新执行。
- [ ] 能说明 `config/local.yml` 里的 `127.0.0.1:3306/6379` 为什么不是容器服务名。

> 待验证：`Makefile` 中 `bootstrap` 目标使用了 `nunu run`，但 `go.mod` 未引入该工具，默认环境未必安装。

---

## 2. Go 基础验收

### 2.1 指针、值与 nil

- [ ] 能区分 `T` 与 `*T`：值变量自己装数据，指针变量装地址。
- [ ] 能说明**只有指针/接口/slice/map 等引用类型可以是 nil，普通 struct 值不能**。
- [ ] 能解释为什么 `QueryById` 返回 `*model.Shop`：既避免复制，也能用 `nil` 表示“没查到”。
- [ ] 能解释 `shop, err := First()` 查不到时返回 `(nil, gorm.ErrRecordNotFound)`。
- [ ] 知道对 nil 指针取字段会 panic。

验收代码：能默写

```go
shop, err := s.query.Shop.Where(s.query.Shop.ID.Eq(id)).First()
if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, nil // 用 nil 表达“没有”
}
if err != nil {
    return nil, err
}
```

项目锚点：`internal/service/shop.go`、`internal/service/blog.go`。

### 2.2 隐式接口实现

首次阅读入口：[小课 01](../lessons/0001-user-service-interface.html#familiar)。先能区分「接口约定、具体实现、创建对象」三个角色，再做本节默写和 `&` 的自查题。

先用大白话：

- `implements` 是 TypeScript / Java 的**显式声明关键字**，例如 `class X implements I`，意思是“我声明要实现 I，请编译器检查我方法写全没有”。
- Go 没有这个关键字，也没有 class。只要给 struct 挂上同名同签名的方法，把它赋给接口变量时编译器自动核对；方法齐全就算实现。
- 所以“Go 不写 `implements`”= 少了“主动声明”这一步，不是少了“检查”这一步。

**interface 和 struct 的区别：**

| 对比项 | struct | interface |
| --- | --- | --- |
| 内容 | 字段 + 方法，描述“有什么数据、能做什么” | 只有方法签名，描述“必须能做什么” |
| 能否直接创建实例 | 能，`T{...}` | 不能，只能装某个具体实现 |
| 能否有字段 | 能 | 不能（TS 的 interface 可以描述字段，Go 不行） |
| 用途 | 存数据、写实现 | 抽象、解耦、多态 |

**大小写不是类型区别，是跨包可见性：**

- `UserService` 大写开头 = 导出（exported），其他包能写 `service.UserService`。
- `userService` 小写开头 = 未导出（unexported），只有 `service` 包内部能用。
- 项目约定：接口用大写 `UserService`，具体实现用小写 `userService`，避免同名冲突并隐藏实现。

- [ ] 能说明 Go 没有 `implements` 关键字。
- [ ] 能说明“结构体的方法集凑齐接口要求的方法”后，编译器自动认定它实现了接口。
- [ ] 能指出检查发生在什么时候：把具体类型赋给接口变量 / 作为接口返回值返回时。
- [ ] 能解释 Handler 为什么只依赖 `service.UserService` 接口，而不是具体结构体。

验收代码：能默写一个接口 + 实现 + 构造函数

```go
type UserService interface {
    SendCode(ctx context.Context, req *v1.SendCodeReq) error
}

type userService struct {
    *Service
}

func (s *userService) SendCode(ctx context.Context, req *v1.SendCodeReq) error {
    return nil
}

func NewUserService(service *Service) UserService {
    return &userService{Service: service}
}
```

自查题：把上面 `return &userService{...}` 的 `&` 去掉，为什么不通过编译？

项目锚点：`internal/service/user.go` 的 `NewUserService`、`internal/handler/user.go` 的 `UserHandler` 字段。

### 2.3 方法接收者与方法集

先纠正一个常见理解：`func (s *T) M()` 不是“函数写在 T 里面”，而是**这个方法被登记到 `*T` 的方法集里**。它和 `T` 有关，但 `T` 本身不包含这个方法。

“方法集”= 这个类型身上能调用、并用于接口匹配的方法清单。以本项目为例：

```text
方法集(userService)  = {}
方法集(*userService) = {SendCode, Login, Me, ...}
```

接口 `UserService` 要求 `SendCode` 等方法，所以 `userService` 不满足，`*userService` 满足。

- 调用 `p.M()` 时，如果 `M` 是 `*T` 接收者且 `p` 可寻址，编译器会自动改写成 `(&p).M()`，所以看起来像“T 也有这个方法”。
- 但接口匹配不会自动加 `&`：`var i I = p` 不成立，`var i I = &p` 才成立。

- [ ] 能说出 `func (s *T) M()` 中的 `s` 相当于 TS 的 `this`。
- [ ] 能说明指针接收者的三个作用：避免复制、可修改原对象、保持方法集一致。
- [ ] 记住方法集规则：

| 方法定义 | `T` 是否实现接口 | `*T` 是否实现接口 |
| --- | --- | --- |
| `func (t T) M()` | ✅ | ✅ |
| `func (t *T) M()` | ❌ | ✅ |

- [ ] 能解释为什么构造函数统一返回 `&xxx{}`。

项目锚点：`internal/service/shop.go`、`internal/service/user.go`。

### 2.4 error / context / defer / goroutine

- [ ] 能默写 `if err != nil { return ..., err }`，并说明 Go 没有 try/catch。
- [ ] 能说明 `errors.Is` 用来判断特定错误（如 `gorm.ErrRecordNotFound`、`redis.Nil`）。
- [ ] 能说明 `context.Context` 贯穿 handler → service → DB/Redis，用于超时、取消和请求级数据。
- [ ] 能说明 `ctx.Request.Context()` 与 `*gin.Context` 不是同一个东西。
- [ ] 能默写中间件里“写入用户到 request context”的写法。
- [ ] 能说明 `defer` 用于释放锁、关闭资源，且在函数返回前执行。
- [ ] 能说明 `go func(){}()` 启动后台 goroutine，`chan` 用于队列。

项目锚点：`internal/middleware/refresh_token.go`、`internal/service/voucher_order.go`、`internal/base/user_holder/user_holder.go`。

---

## 3. Gin 验收

### 3.1 Handler 签名

- [ ] 能说出 Gin handler 的统一形状：`func(ctx *gin.Context)`，**一个参数、无返回值**。
- [ ] 能解释为什么没有返回值：响应通过 `ctx.JSON` / `ctx.String` 写出，`return` 只用于提前结束。
- [ ] 能对比 Express：`(req, res)` 被合并成 `*gin.Context`。

验收代码：能默写

```go
func (h *UserHandler) SendCode(ctx *gin.Context) {
    var req v1.SendCodeReq
    if err := ctx.ShouldBindQuery(&req); err != nil {
        v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
        return
    }
    if err := h.userService.SendCode(ctx.Request.Context(), &req); err != nil {
        v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
        return
    }
    v1.HandleSuccess(ctx, nil)
}
```

自查题：`func (h *UserHandler) SendCode(ctx *gin.Context) error` 为什么不能注册到 Gin？

### 3.2 逐字理解方法定义

```go
func (h *UserHandler) SendCode(ctx *gin.Context) {
```

- [ ] `(h *UserHandler)` 是接收者，`h` 是这次调用的 “this”。
- [ ] `SendCode` 首字母大写表示可被其他包引用。
- [ ] `(ctx *gin.Context)` 是唯一参数。
- [ ] 没有返回类型。
- [ ] 路由注册写 `userHandler.SendCode` 时，`h` 已经绑定，Gin 调用时只传 `ctx`。

### 3.3 路由与分组

- [ ] 能默写 `gin.Default()`、`Group("/xxx")`、`GET/POST/PUT`。
- [ ] 能说明 `Group` 只拼路径前缀，不执行逻辑。
- [ ] 能说明全局中间件 `Use` 与分组中间件 `Use` 的差别。
- [ ] 能指出 `/user/code` 与 `/user/me` 在 `internal/server/http.go` 里为什么一个有 `Login()`、一个没有。

### 3.4 参数绑定

- [ ] 能默写三种绑定：

| 参数位置 | 写法 | 项目例子 |
| --- | --- | --- |
| Query | `ctx.ShouldBindQuery(&req)` | `SendCode` |
| JSON Body | `ctx.ShouldBind(&req)` | `Login` |
| Path | `ctx.ShouldBindUri(&req)` | `QueryUserByID` |

- [ ] 能写结构体 tag：`json:"..."`、`uri:"id"`、`binding:"required"`。

### 3.5 响应

- [ ] 能默写 `ctx.JSON(http.StatusOK, resp)`、`ctx.String(...)`。
- [ ] 能说明项目统一响应结构 `v1.Response` 和 `HandleSuccess` / `HandleError`。

### 3.6 中间件

- [ ] 能默写中间件模板，并说明 `ctx.Next()` 与 `ctx.AbortWithStatus()` 的区别。

```go
func Login() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        if user_holder.GetUser(ctx.Request.Context()) == nil {
            ctx.AbortWithStatus(401)
            return
        }
        ctx.Next()
    }
}
```

- [ ] 能说明中间件和 handler 的签名相同，靠是否调用 `Next()` 区分角色。

项目锚点：`internal/middleware/login.go`、`internal/middleware/cors.go`、`internal/middleware/refresh_token.go`、`internal/server/http.go`。

---

## 4. GORM / GORM Gen 验收

- [ ] 能说明 `internal/model`、`internal/query` 是生成代码，不手工长期维护。
- [ ] 能默写连接与查询入口：

```go
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
q := query.Use(db)
```

- [ ] 能默写以下链式操作：

| 操作 | 写法 |
| --- | --- |
| 按 ID 查一条 | `.Where(q.Shop.ID.Eq(id)).First()` |
| 条件查询 | `.Where(...).Find()` |
| 分页 | `.FindByPage(offset, size)` |
| 模糊查询 | `.Where(q.Shop.Name.Like("%"+name+"%"))` |
| 排序 | `.Order(q.Shop.Liked.Desc())` |
| 批量 ID | `.Where(q.Shop.ID.In(ids...))` |
| 更新 | `.Updates(&shop)` / `.Update(field, field.Add(1))` |
| 计数 | `.Count()` |
| 保存 | `.Save(&obj)` |
| 事务 | `q.Transaction(func(tx *query.Query) error { ... })` |

- [ ] 能说明 `ErrRecordNotFound` 与 `RowsAffected` 的用途。
- [ ] 能说明事务回调返回 `error` 决定回滚/提交，事务内应使用 `tx` 而不是原 `q`。

项目锚点：`internal/service/service.go`、`internal/service/shop.go`、`internal/service/blog.go`、`internal/service/voucher_order.go`、`cmd/generate/main.go`。

---

## 5. Redis 验收

### 5.1 数据结构与场景

- [ ] 能默写“结构 → 方法 → 项目场景”映射：

| 结构 | 必会命令 | go-redis 方法 | 项目场景 |
| --- | --- | --- | --- |
| String | `SET/GET/EXPIRE/SETNX` | `Set/Get/Expire/SetNX` | 验证码、缓存、库存 |
| Hash | `HSET/HGETALL` | `HSet/HGetAll().Scan(&obj)` | 登录用户信息 |
| Set | `SADD/SMEMBERS/SINTER` | `SAdd/SMembers/SInter` | 关注、共同关注 |
| ZSet | `ZADD/ZSCORE/ZREM/ZREVRANGE` | `ZAdd/ZScore/ZRem/ZRevRange` | 点赞、热度 |
| ZSet 分页 | `ZREVRANGEBYSCORE` | `ZRevRangeByScoreWithScores` | Feed 滚动分页 |
| Bitmap | `SETBIT/BITCOUNT/BITFIELD` | `SetBit/BitField` | 签到 |
| GEO | `GEOADD/GEOSEARCH` | `GeoSearchLocation` | 附近商户 |
| Lua | `EVAL/SCRIPT LOAD` | `redis.NewScript(...).Run(...)` | 秒杀原子判断 |
| Stream | `XREADGROUP/XACK` | `XReadGroup/XAck` | 异步订单 |

### 5.2 固定习惯

- [ ] 所有操作带 `ctx`：`rdb.Set(ctx, key, value, ttl).Err()`。
- [ ] 区分 `redis.Nil`（key 不存在，不是故障）与真实错误。
- [ ] Hash 转结构体：`rdb.HGetAll(ctx, key).Scan(&obj)`。
- [ ] key 统一从 `internal/base/constants/redis.go` 取，不散落魔法字符串。
- [ ] 能说明缓存穿透、击穿、逻辑过期三种策略各自解决什么问题。

项目锚点：`internal/base/cache_client/cache_client.go`、`internal/service/user.go`、`internal/service/blog.go`、`internal/service/shop.go`、`internal/service/voucher_order.go`。

---

## 6. 综合验收：3 个闭环练习

每个练习都要能独立完成，并用命令验证结果。

### 练习 1：Gin 闭环

- [ ] 新增 `GET /ping`，返回 `{"success":true}`。
- [ ] 新增一个记录耗时的中间件并挂到该路由。
- [ ] 用 `curl` 验证，并在日志中看到耗时。

通过标准：能解释路由注册、`ctx`、中间件 `Next()` 各自的作用。

### 练习 2：GORM Gen 闭环

- [ ] 仿照 `QueryShopByName`，实现“按评分区间查商户”。
- [ ] 支持分页，并按评分倒序。
- [ ] 用接口或测试验证返回结果。

通过标准：能解释 `Where`、`Order`、`FindByPage` 生成的 SQL 大致长什么样。

### 练习 3：Redis 闭环

- [ ] 仿照 `LikeBlog`，实现“收藏 / 取消收藏”。
- [ ] Redis 用 ZSet 存用户，数据库维护计数。
- [ ] 用 `redis-cli MONITOR` 观察实际命令。

通过标准：能解释 Redis 状态与数据库计数为什么要双写，以及失败时可能不一致。

---

## 7. 明确不需要背

- `internal/query/*.gen.go`、`internal/model/*.gen.go`
- `cmd/server/wire/wire_gen.go`
- Swagger 注解的全部细节
- Dockerfile、Nginx 配置细节
- `lancet`、`lo`、`copier` 等工具库的完整 API

这些用的时候查即可，重点是知道它们“负责什么”。

---

## 8. 关联文件

- 启动与装配：`cmd/server/main.go`、`cmd/server/wire/wire.go`
- 路由与中间件：`internal/server/http.go`、`internal/middleware/`
- Handler：`internal/handler/`
- Service：`internal/service/`
- 缓存与锁：`internal/base/cache_client/`、`internal/base/redis_lock/`
- 模型与查询：`internal/model/`、`internal/query/`
- 配置：`config/local.yml`
- 依赖：`deploy/docker-compose/docker-compose.yml`
