# 黑马点评 Go 学习基础设施计划

## Summary

为当前 fork 的黑马点评 Go 后端建立一套仓库内学习协作基础设施，服务于长期的“问题驱动 + 自动沉淀”学习方式：

- Agent 平时负责答疑、带读代码、排查问题、解释设计取舍。
- 路线用于定位项目地图和知识缺口，不作为强制排课。
- 稳定结论按主题增量沉淀，避免每次对话生成重复流水账。
- 首轮只创建文档骨架和工具清单，不修改业务实现、不安装插件、不把全部规则做成 Skill。

推荐采用：

1. `AGENTS.md` 作为仓库级 Agent 行为入口。
2. `docs/learning/` 作为项目事实、学习路线、用户画像、工具和笔记的静态来源。
3. 以后只有出现稳定、可复用且需要执行步骤的 Agent 行为时，才抽取为独立 Skill。

## Current State Analysis

### Learner profile

- 已有较长期经验：HTML、CSS、TypeScript、React、Express、C++。
- 初学且较浅：Docker、Go、Gin。
- 后端基础、数据库/Redis 运维工具、并发与工程化知识由 Agent 辅助补齐。
- 学习偏好：不需要一次性走完整课程；希望围绕实际问题获得答疑、排查和代码理解，并在形成稳定知识后沉淀笔记。

### Repository shape

- 启动入口：`cmd/server/main.go`
- 依赖注入：`cmd/server/wire/wire.go` 与生成文件 `cmd/server/wire/wire_gen.go`
- HTTP 路由与中间件：`internal/server/http.go`、`internal/middleware/`
- 请求处理：`internal/handler/`
- 业务服务：`internal/service/`
- Redis/锁/缓存/ID 等基础设施：`internal/base/`
- GORM Gen 生成模型和查询：`internal/model/`、`internal/query/`
- API 请求响应类型：`api/v1/`
- 数据初始化与本地依赖：`deploy/docker-compose/docker-compose.yml`、`deploy/docker-compose/hmdp.sql`
- 测试目前主要是依赖真实 Redis/MySQL 的集成/实验型测试：`test/app/app_test.go`、`test/redsync/redsync_test.go`
- 现有 `docs/` 只有 Swagger 生成物，不应和学习文档混放。

### Official-course module mapping

学习地图要显式区分“已实现、部分实现、待实现、实现疑点”：

- 短信登录：`internal/service/user.go`、`internal/middleware/refresh_token.go`、Redis Hash/TTL；已实现，需核对异常处理和 token 生命周期。
- 商户查询缓存：`internal/service/shop.go`、`internal/base/cache_client/`；已实现多种策略示例，但互斥锁/逻辑过期分支存在未验证代码。
- 用户签到：`internal/service/user.go` 的 Bit 操作；部分实现，路由/位偏移/月份语义需要核对。
- UV 统计：当前在 `test/app/app_test.go` 中有 HyperLogLog 实验；业务接口尚未形成完整模块。
- 附近商户：`ShopService.QueryShopOfType`、Redis GEO、SQL 查询；部分接通，需核对路由和经纬度命名。
- 关注/共同关注：`internal/service/follow.go`、Redis Set；已实现，但需核对关注用户字段和 Wire 注入完整性。
- Feed/达人探店：`internal/service/blog.go`、Redis ZSet；已实现主链路，需重点学习滚动分页和时间游标。
- 优惠券秒杀：`internal/service/voucher.go`、`voucher_order.go`、`internal/scripts/seckill.lua`、Redis Stream、Redsync；已有完整意图，但消费者初始化、消息组、订单落库和错误恢复需要重点审阅。
- 评论：存在 `blog_comments` 模型、handler、service 和 SQL 表，但基本是占位实现；标记为待实现。
- 生成代码/依赖注入：`internal/model`、`internal/query`、Wire；作为贯穿全项目的工程基础专题。

## Proposed Changes

### 1. `AGENTS.md`

仓库级入口，内容保持短小，只定义 Agent 的默认行为：

- 先读取相关学习文档和源码事实，再回答。
- 不把推测当成项目事实；涉及实现状态时给出文件路径和证据。
- 默认解释“请求入口 -> handler -> service -> DB/Redis -> 响应”的链路。
- 对初学 Go 的内容采用渐进式讲解：先直觉，再 Go 语法，再项目代码，再风险/替代方案。
- 适当关联用户已有的 TypeScript/Express/C++ 经验，但明确 Go 的差异。
- 排查问题时先复现/定位，再提出最小修复；未经要求不做大范围重构。
- 形成稳定结论时按 `docs/learning/notes/` 主题增量写入，并更新索引。
- 变更学习文档时保留历史，不覆盖已有结论。
- 未经明确要求不修改业务代码、不安装工具、不提交 commit。

### 2. `docs/learning/README.md`

学习文档总入口，说明：

- 文档用途和阅读顺序。
- 哪些内容是项目事实，哪些是学习解释，哪些是待验证假设。
- 如何提出问题、请求带读、请求排查、请求复盘。
- 笔记命名、索引更新和增量记录规则。

### 3. `docs/learning/user-profile.md`

记录用户当前基础、已知强项、薄弱点和偏好的教学方式。加入“不要重复讲解”的维护字段：

- 已熟悉：Web 基础、前端工程、Express/C++ 常见概念。
- 当前重点：Go 语法/惯用法、Gin 生命周期、GORM/SQL、Redis 数据结构与一致性、并发、测试和 Docker。
- 每次新增能力时只追加验证过的结论或代码证据。

### 4. `docs/learning/learning-contract.md`

定义 Agent 的讲解风格和协作协议：

- 先回答当前问题，再指出必要前置知识。
- 复杂问题拆成“现象、定位、原理、项目映射、验证、总结”。
- 对代码解释优先使用真实文件和函数，不凭空创建伪代码。
- 发现项目 bug 或疑点时区分“确定 bug、风险、设计选择、待验证”。
- 笔记只有在结论稳定、对未来有复用价值时自动沉淀。
- 复盘可包含 3 至 5 个自测问题，但不强制每次生成。

### 5. `docs/learning/project-map.md`

绘制当前仓库的静态地图，覆盖：

- 启动和依赖注入链路。
- 路由/中间件/handler/service/query/model 的职责。
- MySQL、Redis、Lua、Stream、GEO、ZSet、Set、Bitmap、HLL 的使用位置。
- 配置文件、Compose、SQL 初始化、Swagger、测试和生成代码。
- 每个目录的“应该先看什么”和“常见问题入口”。

### 6. `docs/learning/study-route.md`

不是按日期排课，而是按“遇到问题时的导航路线”组织：

1. 仓库启动与请求链路。
2. Go 工程基础：package、interface、pointer、error、context、goroutine/channel。
3. Gin 与 HTTP：路由、middleware、绑定、响应和认证。
4. 数据访问：SQL 表、GORM Gen、查询对象、事务、错误映射。
5. Redis 基础：String/Hash/List/Set/ZSet/Bitmap/GEO/HLL/Stream。
6. 登录与 token。
7. 缓存穿透、击穿、雪崩和一致性。
8. Feed、点赞、关注和滚动分页。
9. 秒杀：Lua 原子性、库存、幂等、分布式锁、Stream 消费和 pending list。
10. 测试、压测、日志、配置和 Docker。

每一项包含：项目文件、当前状态、核心问题、建议验证方式、前置知识和关联笔记。路线正文初始不假设所有模块已正确实现。

### 7. `docs/learning/module-status.md`

维护官方 Java 课程模块与本 Go 仓库的对照表，字段固定为：

- 模块
- 官方概念
- Go 项目文件
- 状态：已实现 / 部分实现 / 占位 / 实验代码 / 待验证
- 关键数据结构
- 当前疑点
- 入口笔记

这份文件用于避免“目录存在就误判功能完整”，并记录诸如评论占位、签到路由未完整接通、Wire 注入差异、秒杀消费者初始化等问题。

### 8. `docs/learning/tooling.md`

给出不绑定 IDE 的本地工具清单和用途：

- Go：`go test`、`go test -race`、`go vet`、`gofmt`、`go mod`、Wire、Swag。
- Docker：Compose、容器日志、端口、卷、初始化 SQL。
- MySQL/MariaDB：命令行客户端或 GUI 客户端，用于查表、执行 SQL、查看索引和事务。
- Redis：`redis-cli`、GUI 客户端、TTL/Hash/Stream/GEO/Bitmap/HLL 检查方式。
- HTTP：curl、Postman/Apifox 或 IDE HTTP Client，用于登录拿 token 和接口回归。
- 代码阅读：IDE Go 语言支持、跳转定义、调用层级、调试器。
- 压测：现有 `scripts/秒杀抢购.jmx` 的 JMeter 使用方式。

只记录推荐能力、安装判断和常用命令，不假定必须安装某个插件，也不在首轮修改环境。

### 9. `docs/learning/notes/README.md` 与主题笔记

建立索引和最小主题文件：

- `docs/learning/notes/index.md`
- `docs/learning/notes/go.md`
- `docs/learning/notes/gin-http.md`
- `docs/learning/notes/mysql-gorm.md`
- `docs/learning/notes/redis.md`
- `docs/learning/notes/concurrency.md`
- `docs/learning/notes/docker-tooling.md`
- `docs/learning/notes/project-pitfalls.md`

每条笔记使用统一结构：

```md
## 结论标题
### 问题/场景
### 结论
### 项目映射
### 示例或验证
### 易错点
### 关联文件
```

首轮只写入项目已确认的少量事实和待验证项，不虚构学习成果；后续由 Agent 按主题增量追加。

### 10. `docs/learning/changelog.md`

只记录学习基础设施和学习认知的增量变化，例如：

- 新增/调整了哪份学习文档。
- 项目模块状态从什么变成什么，依据是什么。
- 新增了什么稳定笔记。
- 工具建议或协作规则为何调整。

不复制每次对话全文，也不替代 Git commit history。

## Assumptions & Decisions

- 文档放在 `docs/learning/`，与现有 Swagger 生成物隔离。
- 首轮不创建定制 Skill。`AGENTS.md` + 静态文档足以覆盖用户当前的答疑、排查和沉淀需求；只有当后续出现稳定工作流（例如自动生成模块状态、按模板归档笔记）时，再评估 Skill。
- 首轮不修改 `internal/`、`api/`、`config/`、`deploy/`、测试或生成文件。
- 不安装 Trae/VS Code 插件。工具文档只提供选择依据和命令，实际环境按用户需要再配置。
- 项目事实以当前源码为准；官方 Java 课程图只用于建立概念对照，不将图片中所有模块强行认定为当前仓库功能。
- 自动沉淀默认是“稳定结论增量追加”，不是无条件写入全部对话。
- 不在首轮引入外部资料或网络依赖；若后续需要对照官方课程版本，再单独补充来源和版本信息。

## Verification Steps

完成文档建设后执行只读检查：

1. `rg --files AGENTS.md docs/learning` 确认所有入口和主题文件存在。
2. 检查 `AGENTS.md` 是否能在一屏内说明 Agent 默认行为，并链接到学习文档入口。
3. 检查每个主题笔记是否有统一模板，索引是否能反向找到主题文件。
4. 检查 `module-status.md` 是否覆盖用户图片中的主要模块，并为每项提供实际源码路径和状态。
5. 检查 `project-map.md` 是否覆盖启动链路、业务链路、基础设施、数据层、测试和部署。
6. 检查文档中没有把“待验证”写成“已完成”，没有引导执行未确认的 destructive 命令。
7. 运行 `git diff --check`，确认文档没有空白错误。
8. 不运行需要启动数据库/Redis、安装依赖或修改环境的命令；文档建设的验证不应依赖外部服务。

## Acceptance Criteria

- Agent 能从 `AGENTS.md` 找到完整学习协作规则和文档入口。
- 用户画像、讲解风格、项目地图、路线、模块状态、工具清单、笔记索引、变更记录均有明确归属。
- 后续提问可以按真实文件定位，而不是只得到泛化课程回答。
- 新笔记可以按主题增量追加并被索引发现。
- 官方课程模块与 Go 仓库实现差异被显式记录。
- 没有改动现有业务代码或安装外部工具。
