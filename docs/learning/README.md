# 黑马点评 Go 学习总览

这套文档服务于长期的问题驱动学习，不是固定课程，也不替代项目源码。

## 阅读入口

1. 先看本文的学习者画像、讲解契约和工具清单。
2. 需要理解仓库、路线或模块完成度时看 [project.md](project.md)。
3. 需要复习已确认结论时看 [notes.md](notes.md)。
4. 基础设施变化看 [changelog.md](changelog.md)。
5. 需要参考 Java 版黑马点评复盘时看 [java-review-extracted.md](java-review-extracted.md)。

## 学习者画像

- 较熟悉：HTML、CSS、TypeScript、React、Express、C++。
- 正在补强：Go、Gin、Docker。
- 由 Agent 辅助：后端工程结构、MySQL/GORM、Redis、并发、测试、部署和排障。
- 偏好以实际问题为入口，不要求按固定课程一次性走完整链路。
- 讲解可关联 Express、TypeScript 和 C++，但必须明确 Go 的差异。

当前重点：

- Go：package、interface、pointer、error、context、goroutine、channel。
- Web：Gin 路由、中间件、请求绑定、响应和认证。
- 数据：SQL、GORM Gen、事务、错误处理和索引。
- Redis：数据结构、TTL、缓存一致性、Lua、Stream、GEO、Bitmap、HLL。
- 工程：测试、日志、配置、Wire、Docker 和接口验证。

能力记录只追加已经通过代码阅读、修改、测试或解释确认的内容。

## 讲解契约

1. 先直接回答当前问题，再指出必要前置知识。
2. 代码问题优先引用真实文件、函数和调用关系。
3. 复杂问题按“现象、定位、原理、项目映射、验证、总结”组织。
4. 排障先确认复现条件，再沿入口、参数、中间件、数据库/Redis 和响应逐层定位。
5. 区分确定 bug、风险、设计选择和待验证项。
6. 修复优先保持范围小，并说明为什么能解决问题。
7. 只有稳定、可复用且有证据的结论才进入主题笔记。
8. 笔记按主题增量追加，不覆盖历史，不记录每次对话全文。

## 工具清单

| 工具 | 用途 |
| --- | --- |
| `go test ./...` | 编译并运行测试 |
| `go test -race ./...` | 检查数据竞争 |
| `go vet ./...`、`gofmt` | 静态检查和格式化 |
| `wire ./cmd/server/wire` | 重新生成依赖注入代码 |
| `docker compose` | 启动 MariaDB、Redis、Nginx 和查看日志 |
| MySQL/MariaDB 客户端 | 查询表、索引、事务和初始化数据 |
| Redis CLI/GUI | 检查 key、TTL、Hash、Set、ZSet、Stream、GEO、Bitmap、HLL |
| curl、Postman、Apifox | 登录、携带 Bearer token 和接口回归 |
| IDE Go 支持 | 跳转定义、查找引用、调用层级和调试 |
| JMeter | 使用 `scripts/秒杀抢购.jmx` 做秒杀压测 |

不要求安装特定插件。需要真实 MySQL/Redis 的命令，先检查 `config/local.yml` 和 Compose 状态；涉及写入、清理或重置时先确认影响。

## 内容分类

- **项目事实**：来自源码、配置、SQL、测试或命令输出，并带文件路径。
- **学习解释**：为了理解事实而补充的概念、类比和推导。
- **待验证假设**：尚未通过运行、测试或进一步阅读确认的内容。

Agent 回答问题时应尽量区分这三类内容。

## 协作方式

可以直接提出以下类型的问题：

- “带我读懂某个文件/函数。”
- “这个请求从哪里进入，最后如何访问 Redis？”
- “这个报错怎么定位？”
- “这段 Go 和我熟悉的 TypeScript/Express 有什么不同？”
- “把这次确认的结论沉淀到笔记。”

笔记按主题增量追加。只有稳定、未来可能复用的结论才进入笔记；一次性的对话过程保留在聊天中。
