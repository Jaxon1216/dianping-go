# 学习路线：从 Redis 到后端工程

> 原文：[《黑马点评复盘》（Java 版）](https://kcnebgoczud9.feishu.cn/wiki/A2X4wuqX6iEbkmk41L9cAtminYo)
> 提炼稿：[java-review-extracted.md](java-review-extracted.md)
> 逐课展开：[courses/](courses/)
> 源码事实：以仓库当前代码为准；模块完成度见 [project.md](project.md)

## 这份路线解决什么问题

原文是 Java 版黑马点评的复盘，共 10K 字，结构是「Redis 基础 → 7 个业务功能 → 总结」。本项目用 Go 实现了**同一套业务**，但原文有两个不足：

1. 偏结论清单，只写「用什么方案」，很少写「怎么落到代码」。
2. 直接用 Java 生态术语（`StringRedisTemplate`、`Redisson`、`synchronized`），照搬会误导。

所以这份路线做三件事：保留原文的问题演进顺序，补上 Go 源码入口，给每个阶段加可验证的完成标准。

原文的核心判断是：**整个项目的重点在优惠券秒杀**，其余功能主要是 Redis 数据结构的应用题。下面的排序遵循这个判断。

## 三条学习规则

每个知识点必须同时满足三条，否则算没学会：

1. 能用自己的话说出**它解决什么问题**，以及**不用它会怎样**。
2. 能在源码里指出**具体文件和函数**。
3. 能说出它在**并发或故障**下的失效边界。

达不到第三条时，说明还停在「会背」，没到「会判断」。

## 六个阶段

### 阶段 1：Redis 基础与数据结构选型

- 目标：看到业务需求能选出数据结构，而不是先背命令。
- 原文对应：`1.Redis基础知识`
- 必读源码：[constants/redis.go](../../internal/base/constants/redis.go)、[project.md](project.md) 的 Redis 使用位置清单
- 必须能回答：
  - 验证码为什么用 String，登录用户为什么用 Hash？
  - 关注关系为什么用 Set，点赞排行榜为什么必须换 ZSet？
  - 签到为什么用 Bitmap，UV 为什么用 HyperLogLog？
- 完成标准：给一个「记录每天是否打卡」的需求，能直接说出 key 设计、类型和位偏移含义（对照 [user.go](../../internal/service/user.go) 的签到实现验证）。

### 阶段 2：请求链路与登录态

- 目标：打通「请求怎么进来、用户身份怎么传下去」。
- 原文对应：`2.1短信登陆`、`登录拦截器优化`
- 必读源码：[server/http.go](../../internal/server/http.go) 的路由与中间件、[refresh_token.go](../../internal/middleware/refresh_token.go)、[login.go](../../internal/middleware/login.go)、[user_holder](../../internal/base/user_holder/user_holder.go)
- 必须能回答：
  - 为什么登录态放 Redis，而不是进程内存或 Session？
  - `RefreshToken` 和 `Login` 两个中间件各做什么，为什么不能合并？
  - token 续期发生在哪一步，续期失败会怎样？
- 完成标准：能画出「浏览器 → 全局中间件 → 路由分组 → Handler → Service → Redis」一次完整往返。
- 对应课程：[C01](courses/C01-启动分层与依赖组装.md)、[C02](courses/C02-前后端HTTP契约.md)、[C03](courses/C03-Redis会话与滑动续期.md)

### 阶段 3：读性能与缓存

- 目标：理解缓存不是「加一层就快了」，而是引入了一致性问题。
- 原文对应：`2.2商户查询缓存`
- 必读源码：[cache_client.go](../../internal/base/cache_client/cache_client.go)、[shop.go](../../internal/service/shop.go)
- 必须能回答：
  - 缓存穿透、击穿、雪崩的**触发条件**分别是什么，为什么不能混为一谈？
  - 互斥锁和逻辑过期两种击穿方案，谁阻塞用户、谁返回旧数据？
  - 更新时先写库还是先删缓存，为什么？
- 完成标准：能指出本项目里三种缓存策略各自的实现函数，并说明哪一种主链路实际在用。
- 对应课程：[C04](courses/C04-CacheAside与缓存一致性.md)

### 阶段 4：关系与时间序列

- 目标：理解「用 Redis 表达关系」的取舍。
- 原文对应：`2.4达人探店`、`2.5好友关注`、`2.6附近商铺`
- 必读源码：[follow.go](../../internal/service/follow.go)（`SAdd`、`SInter`）、[blog.go](../../internal/service/blog.go)（`ZRevRange`、`ZRevRangeByScoreWithScores`）、[shop.go](../../internal/service/shop.go)（`GeoSearchLocation`）
- 必须能回答：
  - 共同关注为什么用 Set 交集，而不是查两次数据库再比对？
  - Feed 为什么用 ZSet，滚动分页为什么需要「时间戳 + 同分值偏移」两个条件？
  - GEO 的 key 为什么按商户类型分组？
- 完成标准：能解释同一时间戳的多个笔记为什么会导致「重复或漏读」，并说出游标该怎么设计。
- 对应课程：[C05](courses/C05-GEO附近商户全链路.md)、[C06](courses/C06-关注与点赞的数据一致性.md)、[C07](courses/C07-Feed推送与滚动分页.md)

### 阶段 5：并发与异步（本项目核心）

- 目标：把「超卖、一人一单、异步落库」串成一条完整的并发控制链。
- 原文对应：`2.3优惠券秒杀`（含分布式锁、Redisson、秒杀优化）
- 必读源码：[seckill.lua](../../internal/scripts/seckill.lua)、[voucher_order.go](../../internal/service/voucher_order.go)、[redis_lock.go](../../internal/base/redis_lock/redis_lock.go)、[redis_worker.go](../../internal/base/redis_worker/redis_worker.go)
- 必须能回答：
  - 为什么「先查库存再扣库存」一定会超卖？乐观锁改了什么？
  - 为什么 Redis 里的资格判定要放进同一段 Lua？
  - 为什么 Stream 消费既要 ACK，又要考虑 pending 消息？
  - 「准入成功」和「落库成功」为什么是两件事？
- 完成标准：能画出「用户请求 → Lua 判定 → Stream → 消费协程 → 分布式锁 → MySQL 事务」全链路，并说出每个环节失败后的结果。
- 对应课程：[C08](courses/C08-秒杀资格原子判定.md)、[C09](courses/C09-Stream异步落库与事务.md)

### 阶段 6：证据边界

- 目标：区分「代码存在」「编译通过」「功能正确」「并发安全」「性能达标」。
- 原文对应：`3.总结`
- 必读源码：[project.md](project.md) 的模块状态表、[app_test.go](../../test/app/app_test.go)
- 必须能回答：
  - 目录里有 handler 和 model，能不能说明这个模块可用？
  - 压测通过能不能说明没有超卖？
  - 哪些模块当前只是占位或实验代码？
- 完成标准：能对着模块状态表，指出至少三个「看起来有、实际没闭环」的地方。
- 对应课程：[C10](courses/C10-部署测试与可观测验证.md)

## 原文功能与 Go 项目覆盖对照

| 原文功能 | Go 项目位置 | 状态 | 本路线阶段 |
| --- | --- | --- | --- |
| Redis 数据类型 | `internal/base/constants/redis.go` | 已实现 | 阶段 1 |
| 短信登录 | `internal/service/user.go`、`refresh_token.go` | 部分实现（未接真实短信） | 阶段 2 |
| 商户查询缓存 | `cache_client.go`、`shop.go` | 部分实现 | 阶段 3 |
| 优惠券秒杀 | `voucher_order.go`、`seckill.lua` | 部分实现（消费组初始化为空） | 阶段 5 |
| 达人探店 / 点赞 | `blog.go` | 已实现 | 阶段 4 |
| 好友关注 / 共同关注 | `follow.go` | 部分实现 | 阶段 4 |
| 附近商铺 | `shop.go` | 部分实现 | 阶段 4 |
| 用户签到 | `user.go` | 部分实现 | 阶段 1 |
| UV 统计 | `test/app/app_test.go` | 仅实验，无业务接口 | 选读 |

「部分实现」的具体疑点见 [project.md](project.md) 的模块状态表，不要把它当作已完成。

## 原文不能照抄的三处

这三处原文说法在 Go 项目里已经被更严谨的实现替代，照抄会得出错误结论。

1. **一人一单不是「给方法加悲观锁」**。原文写「给创建订单这个方法加锁」，这在集群下会失效。本项目实际是两段防护：Lua 用 `sismember` + `sadd` 做 Redis 内原子判定，数据库侧再用分布式锁兜底。见 [seckill.lua](../../internal/scripts/seckill.lua)。
2. **缓存策略要落到三个不同函数**。原文只列了「互斥锁 / 逻辑过期」两个名词。本项目对应 `QueryWithPassThrough`、`QueryWithMutex`、`QueryWithLogicalExpire`，见 [cache_client.go](../../internal/base/cache_client/cache_client.go)。
3. **Stream 不只是「异步队列」**。原文只画了流程图，没讲消费组、ACK、pending。本项目真实代码里有 `XReadGroup`、`Group: "g1"`、`Consumer: "c1"`，见 [voucher_order.go](../../internal/service/voucher_order.go)。消息重放和恢复问题都藏在这里。

## 为什么按这个顺序

原文按功能罗列（登录 → 缓存 → 秒杀 → 探店 → 关注 → GEO → 签到 → UV），8 个模块并列，看不出依赖。改成按**能力递进**排序，理由如下：

- 阶段 1 是所有后续阶段的语言，不懂数据结构选型，后面只能背结论。
- 阶段 2 打通「请求怎么被处理」，是读任何业务代码的前提。
- 阶段 3 引入第一个真正的工程问题：一致性与失效。这是从「调 API」转向「做系统」的分界。
- 阶段 4 全是数据结构的应用题，学起来最快，适合放在最难的部分之前建立信心。
- 阶段 5 是原文认定的重点，需要前面所有能力打底，所以放最后讲透。
- 阶段 6 收口，防止把「有代码」误认为「已完成」。

## 与笔记的关系

- 路线负责**顺序和目标**，不负责记住每个细节。
- 稳定结论沉淀到 [notes.md](notes.md)；需要逐条展开时看 [courses/](courses/)。
- 需要面试表达时看 [interview.md](interview.md)。
- 外部资料（原文、官方文档、社区）见 [resources.md](resources.md)。

## 前置知识按需补齐

不要一开始就把这些全学完，卡在哪补哪。

| 知识点 | 会卡住什么 | 项目位置 | 补习入口 |
| --- | --- | --- | --- |
| Go 指针、接口、error | `service` 构造器返回接口，DTO 含指针字段，错误沿层返回 | `internal/service/*.go` | [Go 基础笔记](notes/go-basics-01-pointer-interface-error.md) |
| Gin 与 `context.Context` | Gin 负责 HTTP 流程，request context 负责跨层携带用户和取消信号 | `internal/middleware` | [请求链路笔记](notes/go-web-request-flow-and-gmp.md) |
| SQL 索引与事务 | 秒杀落库需要条件更新、唯一性和事务边界 | `hmdp.sql`、`voucher_order.go` | [C09](courses/C09-Stream异步落库与事务.md) |
| Redis 数据结构 | 不同业务分别依赖 Hash、String、Set、ZSet、GEO、Bitmap、Stream | `internal/service` | 按阶段 1 和阶段 4 的场景学，不先背命令大全 |
| 项目分层判断 | 目录名容易被误读成 DDD | `internal/` 全目录 | [分层架构判断](notes/architecture-layered-vs-ddd.md) |

## 选读与进阶

主线走完再看这些，它们不进入面试主线。

| 主题 | 值得学的理由 | 当前状态 |
| --- | --- | --- |
| Bitmap 连续签到 | 学位运算与紧凑状态表达 | 有接口和实现，未纳入简历支柱 |
| HyperLogLog UV | 学近似去重与误差控制 | 只有集成实验，无业务接口 |
| 评论系统 | 学层级评论与计数一致性 | service 为占位，路由未完整接通 |
| 缓存逻辑过期 | 热点数据可接受短暂陈旧时的方案 | 有备用代码，主查询未启用 |
| 消息可靠性闭环 | 生产级秒杀必须补齐 | 待实现：消费组初始化、重试/死信、补偿对账 |
| 数据库最终幂等 | 多实例或消息重放下保证一人一单 | 待实现：`(user_id, voucher_id)` 唯一索引与重复键处理 |

## 一句话记忆提纲

- **先懂职责与状态归属（阶段 1-2）**：前端负责交互和凭据传输，服务端负责鉴权和业务规则；Wire 生成构造调用，Redis 保存共享会话。
- **再懂索引与事实源（阶段 3-4）**：MySQL 保存业务事实，Redis 分别承担缓存、空间候选、集合和时间序列；双写不自动具有事务性。
- **重点懂原子边界（阶段 5）**：Lua 保护 Redis 内步骤，事务保护 MySQL 内步骤，消息连接两者。准入成功不等于落库成功。
- **最后懂证据边界（阶段 6）**：不能把负载均衡说成高可用，也不能把编译通过说成并发正确。
