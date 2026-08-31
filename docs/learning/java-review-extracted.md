# 黑马点评原版复盘提炼

> 来源： [黑马点评复盘（Java 版）](https://kcnebgoczud9.feishu.cn/wiki/A2X4wuqX6iEbkmk41L9cAtminYo)
>
> 用途：保存原版复盘中的核心知识，并映射到当前 Go 项目。原文使用 Java、Spring 和 `StringRedisTemplate`，当前项目使用 Go、Gin、GORM 和 `go-redis`。涉及当前项目的结论以源码为准，未核对的地方标记为“待验证”。

## 1. 项目目标与学习重点

黑马点评是一个类似大众点评的 Web 项目，主要通过 Redis 解决缓存访问、登录状态、排行榜、附近商户、统计和高并发秒杀问题。

原版复盘将功能分成：

1. 短信登录
2. 商户查询缓存
3. 优惠券秒杀
4. 达人探店
5. 好友关注
6. 附近商铺
7. 用户签到
8. UV 统计

学习重点不是记住某个 Java API，而是理解：

```text
请求入口
  -> 业务判断
  -> Redis/MySQL 状态变化
  -> 并发与一致性处理
  -> 返回结果
```

其中优惠券秒杀最集中地体现了库存一致性、幂等、一人一单、分布式锁、Lua 原子性和异步消息处理。

## 2. Redis 基础

Redis 是基于内存的键值型 NoSQL 数据库，也可以理解为一个远程数据结构服务器。它的速度快，原因之一是主要在内存中读写；代价是需要考虑容量、持久化、过期和数据丢失风险。

### 2.1 常用数据结构

| 类型 | 直觉 | 本项目中的用途 |
| --- | --- | --- |
| String | 一个 key 对应一个值 | 验证码、缓存、库存、过期时间 |
| Hash | 一个 key 下有多个字段 | 登录用户信息 |
| List | 有序链表/队列 | 商户类型缓存 |
| Set | 无序且不重复的集合 | 关注关系、共同关注、重复购买判断 |
| ZSet/Sorted Set | 带 score 的有序集合 | 点赞排行榜、Feed 时间排序 |
| Bitmap | 按 bit 保存状态 | 用户签到 |
| GEO | 保存经纬度并进行距离查询 | 附近商户 |
| HyperLogLog | 低内存估算去重数量 | UV 统计实验 |
| Stream | 带消费组和确认机制的消息流 | 秒杀订单异步处理 |

Java 原版通过类似下面的 API 操作 Redis：

```java
stringRedisTemplate.opsForValue().set(key, value, ttl);
stringRedisTemplate.opsForValue().get(key);
stringRedisTemplate.expire(key, ttl);
```

Go 项目使用 `github.com/redis/go-redis/v9`，调用方式不同，但 Redis 数据结构和并发原理相同。

相关代码：

- `internal/base/constants/redis.go`
- `internal/service/`
- `internal/base/cache_client/`
- `internal/base/redis_lock/`
- `internal/base/redis_worker/`

### 2.2 选择数据结构的记忆方式

```text
需要一个值             -> String
需要多个字段            -> Hash
需要去重                -> Set
需要按分数或时间排序     -> ZSet
需要按日期记录 true/false -> Bitmap
需要按距离查询           -> GEO
需要估算独立人数         -> HyperLogLog
需要异步消费消息         -> Stream
```

## 3. 短信登录与 Token

### 3.1 Session 的集群问题

单机应用可以把登录用户放在 Session 中，但多节点服务存在问题：

```text
请求 1 -> Tomcat A，Session 在 A
请求 2 -> Tomcat B，B 没有 A 的 Session
```

如果没有共享 Session，用户切换节点后可能被判断为未登录。

### 3.2 Redis 登录状态

常见方案是：

1. 用户输入手机号。
2. 服务生成验证码并写入 Redis，设置短 TTL。
3. 用户提交手机号和验证码。
4. 服务查询数据库中的用户，不存在就创建。
5. 生成随机 Token。
6. 以 Token 为 key，将用户信息写入 Redis Hash，并设置 TTL。
7. 后续请求携带 `Bearer Token`。
8. 中间件根据 Token 查询 Redis，恢复当前用户。
9. 每次请求刷新 Token TTL。

当前 Go 项目对应链路：

```text
internal/server/http.go
  -> middleware.RefreshToken
  -> Redis HGetAll
  -> internal/base/user_holder
  -> 需要登录的路由使用 middleware.Login
  -> handler -> service
```

相关文件：

- `internal/middleware/refresh_token.go`
- `internal/middleware/login.go`
- `internal/base/user_holder/user_holder.go`
- `internal/service/user.go`
- `internal/base/constants/redis.go`

原版中的 `ThreadLocal` 可以理解为“当前请求/线程上下文中的用户”。Go 不直接照搬 ThreadLocal，而是通过 `context.Context` 保存请求范围内的用户信息。

## 4. 商户查询缓存

### 4.1 基本读缓存流程

```text
请求商户
  -> 查询 Redis
  -> 命中：直接返回
  -> 未命中：查询 MySQL
  -> 写入 Redis
  -> 返回商户
```

缓存的目标是减少数据库访问，让热点商户读取更快。

当前 Go 项目主要位于：

- `internal/service/shop.go`
- `internal/base/cache_client/cache_client.go`

### 4.2 缓存与数据库一致性

缓存数据来自数据库，但数据库会变化，因此缓存可能过期或与数据库不一致。

常见策略：

- 内存淘汰：Redis 内存不足时淘汰数据。
- TTL 过期：设置过期时间，定期让缓存失效。
- 主动更新：通常先更新数据库，再删除或更新缓存。

最常见的写入原则是：

```text
先写数据库
再删除缓存
```

实际选择要结合并发、失败重试和数据重要性判断，不能机械套用。

### 4.3 缓存穿透

定义：请求的数据在缓存和数据库中都不存在，大量请求会反复查询数据库。

```text
请求不存在的 id
  -> Redis 没有
  -> MySQL 也没有
  -> 每次请求都访问 MySQL
```

措施：

- 缓存空对象，并设置较短 TTL。
- 使用布隆过滤器提前判断数据是否可能存在。
- 增强 ID 复杂度。
- 对热点参数限流。

### 4.4 缓存雪崩

定义：大量缓存 key 在同一时间失效，或者 Redis 故障，导致大量请求同时访问数据库。

措施：

- 给不同 key 设置不同的 TTL。
- Redis 主从、哨兵或集群提高可用性。
- 限流和业务降级。
- 使用多级缓存。

### 4.5 缓存击穿

定义：一个高并发热点 key 突然失效，大量请求同时重建同一个缓存。

常见方案：

- 互斥锁：只允许一个请求重建缓存，其他请求等待或重试。
- 逻辑过期：缓存物理上不立即删除，过期后由后台或一个线程异步重建。

当前 Go 项目的 `shop.go` 已有缓存穿透、互斥锁和逻辑过期相关分支，但具体正确性仍应结合测试和运行结果验证。

## 5. 全局唯一 ID

秒杀订单需要唯一 ID，不能只依赖单机自增数字，因为多实例服务可能同时生成相同 ID。

常见设计要求：

- 全局唯一。
- 趋势递增，便于数据库索引。
- 尽量不暴露业务规模。
- 生成速度快。

可使用数据库自增、UUID、雪花算法或自定义时间戳加序列号方案。当前项目的具体 ID 生成实现需要继续沿 `internal/base` 和订单 service 核对。

## 6. 优惠券秒杀与并发

### 6.1 基本业务流程

```text
用户抢券
  -> 判断秒杀是否开始
  -> 判断库存是否充足
  -> 判断用户是否已经购买
  -> 扣减库存
  -> 创建订单
  -> 返回订单 ID
```

高并发下不能简单地先查库存再更新库存。

### 6.2 超卖问题

错误示例：

```text
库存 = 1
线程 A 查询库存，得到 1
线程 B 查询库存，得到 1
线程 A 扣减为 0
线程 B 也扣减，出现超卖
```

一种数据库方案是乐观锁 CAS：

```sql
UPDATE voucher
SET stock = stock - 1
WHERE id = ? AND stock > 0;
```

通过更新条件保证库存仍然大于 0。也可以比较版本号：

```sql
UPDATE voucher
SET stock = stock - 1, version = version + 1
WHERE id = ? AND version = ?;
```

乐观锁适合冲突相对少、读多写少的场景；悲观锁在操作前先锁住数据，适合冲突频繁但会降低并发。

### 6.3 一人一单

业务规则是同一用户只能购买同一优惠券一次。

仅仅先查询“是否购买过”仍有竞态：

```text
线程 A 查询：没有订单
线程 B 查询：没有订单
线程 A 创建订单
线程 B 也创建订单
```

需要组合使用：

- 数据库唯一索引作为最终约束。
- 事务保证查询和创建的一致性。
- 分布式锁减少同一用户的并发进入。
- Redis Set 或 Lua 在高并发入口快速判断。

## 7. 分布式锁

### 7.1 为什么需要分布式锁

进程内的互斥锁只能保护同一个进程。多台服务实例之间需要所有进程都能看到的锁。

分布式锁通常要求：

1. 多进程可见。
2. 互斥。
3. 高可用。
4. 性能可接受。
5. 释放安全。

### 7.2 Redis SETNX

Redis 的 `SETNX` 只有在 key 不存在时才设置成功，可以用于抢锁。

安全写法应将设置值和过期时间放在一个原子命令中：

```text
SET lock:key unique-token EX 10 NX
```

不能简单拆成：

```text
SETNX lock:key token
EXPIRE lock:key 10
```

因为两条命令之间 Redis 可能宕机，导致锁没有过期时间。

### 7.3 安全释放锁

释放锁前必须确认锁的值属于当前线程/请求：

```text
if GET(lock:key) == my-token:
    DEL(lock:key)
```

判断和删除必须是一个原子操作，否则会出现：

```text
线程 A 的锁过期
线程 B 获得同一个 key
线程 A 醒来后删除了线程 B 的锁
```

通常使用 Lua：

```lua
if redis.call('get', KEYS[1]) == ARGV[1] then
    return redis.call('del', KEYS[1])
end
return 0
```

当前 Go 项目对应：

- `internal/base/redis_lock/redis_lock.go`
- `internal/scripts/unlock.lua`
- `internal/service/voucher_order.go`

本项目使用 `redsync` 辅助分布式锁，不是 Java 原版中的 Redisson。两者都是工具层，仍然需要理解锁的超时、唯一标识和安全释放原理。

## 8. 秒杀异步化：Redis Stream

### 8.1 为什么异步

秒杀入口最重要的是快速判断资格，而不是让大量请求同时执行复杂的数据库下单逻辑。

优化后：

```text
请求
  -> Redis/Lua 判断资格和扣减库存
  -> 成功后写入消息队列
  -> 立即返回订单 ID
  -> 后台消费者异步创建数据库订单
```

消息队列的三个核心价值：

- 解耦生产者和消费者。
- 异步处理，缩短请求耗时。
- 削峰，避免数据库瞬时承受全部请求。

Redis 可用于消息队列的方式：

- List：简单队列。
- Pub/Sub：发布订阅，但消息可靠性较弱。
- Stream：支持消费组、确认和 pending 消息，更适合本项目。

### 8.2 Lua 负责什么

秒杀入口中的 Lua 脚本应将以下操作放在 Redis 内原子执行：

1. 判断库存是否充足。
2. 判断用户是否已经购买。
3. 扣减库存。
4. 记录用户购买资格。
5. 写入订单消息。

只有全部条件满足，才允许进入异步下单。

当前 Go 项目对应：

- `internal/scripts/seckill.lua`
- `internal/service/voucher_order.go`
- `internal/base/redis_worker/redis_worker.go`

需要继续验证的工程问题：

- Stream 消费组是否初始化。
- 消息解析失败如何处理。
- 消费失败是否进入 pending list。
- 服务重启后如何恢复未确认消息。
- 数据库订单失败时是否有补偿。

## 9. 达人探店、点赞与 Feed

### 9.1 发布和查询笔记

基本业务：

- 发布博客/探店笔记。
- 保存文字、图片和作者信息。
- 查询详情。
- 查询热门内容。

当前 Go 项目主要位于：

- `internal/handler/blog.go`
- `internal/service/blog.go`
- `api/v1/blog.go`

### 9.2 点赞

用 Set 保存点赞用户 ID，可以快速判断用户是否点过赞：

```text
blog:liked:{blogId} -> {userId1, userId2, ...}
```

数据库保存点赞总数，Redis Set 保存明细，二者用途不同：

- 数据库计数适合持久化展示。
- Set 适合快速判断重复点赞和取消点赞。

### 9.3 点赞排行榜

使用 ZSet：

```text
key: blog:liked:{blogId}
score: 点赞时间
member: 用户 ID
```

按 score 排序即可查询最早或最新点赞用户。

### 9.4 Feed 流

关注用户发布内容后，将博客 ID 和发布时间写入粉丝的 Feed ZSet：

```text
feed:{userId}
score: 发布时间
member: 博客 ID
```

查询 Feed 时按时间倒序读取，并使用时间游标实现滚动分页，而不是简单使用固定页码。

当前 Go 项目使用：

- `internal/service/blog.go`
- `internal/base/constants/redis.go` 中的 Feed key

## 10. 好友关注与共同关注

关注关系通常同时保存两份：

```text
MySQL: 关注关系的持久化记录
Redis Set: 快速判断和集合运算
```

共同关注可以求两个用户关注集合的交集：

```text
SINTER follow:{userA} follow:{userB}
```

数据库是长期事实来源，Redis 是高性能查询结构。Redis 写入失败、数据库写入失败或两者顺序不一致时，需要考虑补偿和一致性。

当前 Go 项目：

- `internal/service/follow.go`
- `internal/handler/follow.go`
- `internal/query/follow.gen.go`

路由、Wire 注入和字段是否完全接通，当前项目地图中仍标记为待验证。

## 11. 附近商铺与 GEO

可以按商户类型建立 GEO 集合：

```text
shop:type:{typeId}
```

商户新增时写入经纬度：

```text
GEOADD shop:type:{typeId} longitude latitude shopId
```

查询时：

1. 根据商户类型找到 GEO key。
2. 以用户经纬度为中心查询附近商户。
3. 按距离排序。
4. 再根据商户 ID 查询完整商户信息。

当前 Go 项目主要位于 `internal/service/shop.go` 和 `internal/server/http.go`。路由参数、经纬度字段命名和分页细节需要以当前源码与运行验证为准。

## 12. 用户签到与 Bitmap

Bitmap 用 bit 表示某一天是否签到：

```text
key: sign:{userId}:{yyyyMM}
offset: dayOfMonth - 1
value: 1
```

例如：

```text
第 1 天 -> offset 0
第 15 天 -> offset 14
```

签到：

```text
SETBIT sign:{userId}:{yyyyMM} {day-1} 1
```

查询连续签到天数时，可以读取当月 bit，再通过位运算统计连续的 1。

Redis 中 Bitmap 实际复用了 String 的底层存储，因此 Go Redis 客户端通常通过 String/bit 操作 API 使用它。

当前项目有签到模型和 handler/service 代码，但具体使用 Redis Bitmap 还是数据库签到表、月份边界和连续签到实现，需要继续以 `internal/service/user.go`、`internal/handler/user.go` 和路由为准，原版 Java 方案不能直接视为当前 Go 实现。

## 13. UV 与 PV

### 13.1 概念区别

- UV：Unique Visitor，独立访客数。同一天同一访客多次访问，通常只计 1 次。
- PV：Page View，页面访问量。每打开一次页面都可以计数。

### 13.2 HyperLogLog

HyperLogLog 适合在内存占用很小的情况下估算去重数量：

```text
PFADD uv:202608 user_1 user_2 user_3
PFCOUNT uv:202608
```

它是近似统计，不适合要求绝对精确的业务数据。

当前 Go 项目在 `test/app/app_test.go` 中有 HLL 实验，但目前没有正式 UV 接口，项目地图将它标记为实验代码。

## 14. Java 版到 Go 版的阅读映射

| Java 版概念 | 当前 Go 项目 |
| --- | --- |
| Controller | `internal/handler` |
| Spring MVC 路由 | `internal/server/http.go` 中的 Gin 路由 |
| Interceptor | `internal/middleware` |
| ServiceImpl | `internal/service` |
| Entity/DO | `internal/model` |
| MyBatis Mapper | `internal/query` 和 GORM |
| `StringRedisTemplate` | `github.com/redis/go-redis/v9` |
| ThreadLocal 用户 | request context + `internal/base/user_holder` |
| Redisson | `redsync`、`redis_lock` |
| MQ 消费者 | `internal/base/redis_worker` |
| Lua 脚本 | `internal/scripts` |
| Spring 配置 | `config/*.yml` + Viper |
| Spring Boot 启动 | `cmd/server/main.go` + Wire + `pkg/app` |

阅读 Java 原版时，优先提取“为什么这样设计”，再寻找 Go 项目中的等价实现，不要逐行寻找同名类。

## 15. 推荐学习顺序

1. Redis 数据结构和 TTL。
2. `cmd/server/main.go`、Wire 和 Gin 路由。
3. 登录 Token 与 request context。
4. 商户缓存、穿透、雪崩、击穿。
5. 分布式锁和 Lua 安全释放。
6. 秒杀 Lua、Stream 和异步下单。
7. 点赞、关注和 Feed 的 Set/ZSet。
8. GEO 附近商户。
9. Bitmap 签到。
10. HyperLogLog UV 实验。

每个模块都按下面的方式学习：

```text
找到请求入口
  -> 跟到 handler
  -> 跟到 service
  -> 找出 MySQL/Redis 的读写
  -> 说明并发和一致性问题
  -> 用测试、日志或 Redis/MySQL 查询验证
```

## 16. 重要提醒

原版复盘是 Java 学习笔记，适合作为概念地图，不是当前 Go 项目的实现文档。当前项目存在部分未接通或待验证模块，尤其是评论、关注注入、签到、秒杀 Stream 恢复和 HLL 正式业务接口。学习时应以源码、配置、SQL、测试和运行结果共同确认结论。
