# 外部资料索引

记录学习过程中真正用到的外部资料，以及它们的用途和使用边界。只保留能回答「为什么」和「怎么验证」的来源。

## 原始复盘

| 资料 | 用途 | 边界 |
| --- | --- | --- |
| [《黑马点评复盘》（Java 版）](https://kcnebgoczud9.feishu.cn/wiki/A2X4wuqX6iEbkmk41L9cAtminYo) | 问题演进来源：Redis 基础、7 个业务功能、项目总结 | Java/Spring 生态，不能直接照搬到 Go。三处已知偏差写在 [roadmap.md](roadmap.md) |
| [java-review-extracted.md](java-review-extracted.md) | 上述原文的提炼稿，已补 Go 文件映射 | 部分结论仍需对照源码验证 |
| [CSDN 笔记（原文作者推荐）](https://blog.csdn.net/weixin_52152676/article/details/144680041) | 查漏补缺，比原文更细 | 未逐条核对，引用前先对照源码 |

## 本项目内的事实来源

优先看这些，而不是外部文章。

| 资料 | 用途 |
| --- | --- |
| [project.md](project.md) | 项目地图、模块完成度、当前疑点 |
| [roadmap.md](roadmap.md) | 学习顺序、每阶段完成标准 |
| [notes.md](notes.md) | 已确认的稳定结论 |
| [courses/](courses/) | 按调用链展开的 10 节课 |
| [wire_gen.go](../../cmd/server/wire/wire_gen.go) | 运行时真正执行的依赖装配顺序 |

## Go 与框架

| 资料 | 用途 |
| --- | --- |
| [Google Wire README](https://github.com/google/wire/blob/main/README.md) | Wire 是无反射的代码生成工具；项目已停止维护，只作为本仓库既有方案学习 |
| [Gin Quickstart](https://gin-gonic.com/en/docs/quickstart/) | 对照 `gin.Default`、路由注册、监听启动 |
| [go-redis 文档](https://redis.gg/guide-go/) | 对照 `Set`、`HSet`、`ZAdd`、`XReadGroup`、`GeoSearch` 的参数语义 |
| [GORM 文档](https://gorm.io/docs/) | 对照事务、条件更新、预加载和错误处理 |

## 标准库

| 资料 | 用途 |
| --- | --- |
| [net/http Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown) | 区分「停止接收新连接」与「等待活动请求完成」 |
| [context](https://pkg.go.dev/context) | 理解请求取消、超时和跨层传值 |
| [os/signal](https://pkg.go.dev/os/signal) | `SIGINT`、`SIGTERM` 与优雅停机 |
| [builtin panic](https://pkg.go.dev/builtin#panic) | 区分普通 error 与不可恢复终止 |

## Redis

| 资料 | 用途 |
| --- | --- |
| [Redis 命令文档](https://redis.io/docs/latest/commands/) | 查 `SETNX`、`SISMEMBER`、`ZREVRANGEBYSCORE`、`GEOADD`、`SETBIT`、`PFADD` 的精确语义 |
| [Redis Streams 介绍](https://redis.io/docs/latest/develop/data-types/streams/) | 理解消费组、pending 消息、ACK 机制 |
| [Redis Lua 脚本](https://redis.io/docs/latest/develop/interact/programmability/eval-intro/) | 理解脚本原子性边界（不等于事务回滚） |

## 社区

| 资料 | 用途 |
| --- | --- |
| [Go Forum](https://forum.golangbridge.org/) | 验证非项目特有的 Go 实践问题 |
| [Gin Discussions](https://github.com/gin-gonic/gin/discussions) | 查路由、中间件、运行行为的真实案例 |

## 使用约定

1. 外部资料只用来**解释**，不作为项目事实。项目事实以源码、配置、SQL、测试输出为准。
2. 引用外部结论时标注来源；与源码冲突时以源码为准，并记录差异。
3. 不把整篇文章复制进仓库，只保留链接和一句话用途。
