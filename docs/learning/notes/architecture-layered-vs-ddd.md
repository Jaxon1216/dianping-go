# 项目架构判断：分层架构不等于 DDD

## 问题/场景

看到 `handler`、`service`、`model`、`query` 等目录时，容易把项目直接判断为 DDD。

## 结论

本项目当前不是严格的 DDD 架构，更准确的描述是：

> 使用依赖注入的分层单体应用，业务逻辑主要采用 Service + Transaction Script 组织。

它已经按 HTTP、业务协调和数据访问进行了分层，但没有形成独立的领域层。

## 项目映射

真实请求链路是：

```text
router -> middleware -> handler -> service -> GORM Query / Redis -> response
```

- `internal/handler` 负责 HTTP 参数绑定和响应。
- `internal/service` 负责业务流程，但也直接调用 GORM Query、Redis 和分布式锁。
- `internal/model` 主要是由 GORM Gen 按数据库表生成的数据结构。
- `internal/query` 是 GORM Gen 生成的数据访问 API。
- Wire 负责在启动阶段组装 DB、Redis、service、handler 和 server。

## 为什么不是严格 DDD

DDD 不由目录名称决定，关键要看模型和依赖方向。本项目缺少以下核心结构：

1. 独立于 Gin、GORM、Redis 的领域模型和领域规则。
2. 聚合、值对象、领域服务、领域事件和限界上下文等明确建模。
3. 由领域层定义、基础设施层实现的 Repository 接口。
4. 应用层只编排用例、领域层承载核心业务规则的清晰边界。

当前 `Service` 直接持有 `*query.Query`、`*redis.Client` 和 `*redsync.Redsync`，说明业务层和基础设施层仍然紧密耦合。`model.Shop` 等类型同时承担数据库映射、HTTP 输入输出和业务数据的角色，也不是独立的领域对象。

## 与 Express 的类比

它更像常见的 Express 项目：

```text
router/controller -> service -> ORM/Redis
```

加入接口和依赖注入可以提高可替换性与可测试性，但这本身不会自动变成 DDD。

## 易错点

- 有 `service` 接口不等于 DDD。
- 有 `model` 目录不等于存在领域模型。
- 有分层和依赖注入不等于依赖倒置已经覆盖领域与基础设施边界。
- 项目不采用 DDD 不代表架构差；对当前体量，简单分层更容易学习和维护。

## 关联文件

- `internal/server/http.go`
- `internal/handler/shop.go`
- `internal/service/service.go`
- `internal/service/shop.go`
- `internal/service/voucher_order.go`
- `internal/model/shop.gen.go`
- `internal/query/shop.gen.go`
- `cmd/server/wire/wire.go`
