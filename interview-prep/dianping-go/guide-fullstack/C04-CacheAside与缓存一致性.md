# C04 Cache-Aside 与缓存一致性

> 课程性质：核心
> 返回：[主导学](../导学.md)；前置：[C01](C01-启动分层与依赖组装.md)、Redis String、JSON；下一课：[C05 GEO 附近商户](C05-GEO附近商户全链路.md)
> 对应：B3；[Q3](../面经.md#q3-商户详情缓存如何处理穿透与更新一致性)；目标掌握：看懂链路、理解机制、设计与表达。

## 本课重点

场景：多个客户端查询同一个存在的商户、持续查询不存在的商户 ID，并在后台更新该商户。

- 当前商户详情主链只启用 Cache-Aside 与空值缓存；互斥锁、逻辑过期是未接入的备用实现。
- 正常缓存 30 分钟，数据库不存在时缓存空字符串 2 分钟，减少相同无效 ID 的重复回源。
- 更新先写 MySQL 再删 Redis，后续读取重建；这个顺序仍有失败和并发不一致窗口。

## 主链与职责

当前实际启用的是 Cache-Aside 加空值缓存：

```text
GET /shop/:id
  -> Redis GET cache:shop:{id}
     -> 命中 JSON：反序列化并返回
     -> 命中空字符串：返回不存在
     -> key 不存在：查询 MySQL
        -> 有数据：JSON 写 Redis，TTL 30 分钟
        -> 无数据：空字符串写 Redis，TTL 2 分钟

PUT /shop
  -> 更新 MySQL
  -> 删除 cache:shop:{id}
  -> 下一次读取回源并重建
```

空值缓存解决“相同不存在 ID 反复打数据库”的缓存穿透。它没有解决热门 key 到期瞬间的大量并发回源，也不等于布隆过滤器。

## 源码阅读

1. [`internal/service/shop.go`](../../../internal/service/shop.go)：`QueryById` 实际调用 `QueryWithPassThrough`；互斥锁和逻辑过期调用被注释。
2. [`internal/base/cache_client/cache_client.go`](../../../internal/base/cache_client/cache_client.go)：顺序阅读 `QueryWithPassThrough`、`QueryWithMutex`、`QueryWithLogicalExpire`。
3. [`internal/base/constants/redis.go`](../../../internal/base/constants/redis.go)：正常缓存 30 分钟，空值 2 分钟。
4. [`internal/handler/shop.go`](../../../internal/handler/shop.go)：详情查询和更新的 HTTP 边界。
5. [`internal/query/shop.gen.go`](../../../internal/query/shop.gen.go)：`GetByID` 的生成 SQL 与类型化字段。

主链关键判断是：

```go
shopJson, err := c.rdb.Get(ctx, key).Result()
if shopJson != "" {
    // 正常命中
}
if err == nil {
    // Redis 命中了显式缓存的空字符串
}
// redis.Nil 才表示 key 真正不存在，需要回源
```

“空字符串”和“没有 key”都表现为 `shopJson == ""`，必须结合 `err` 区分。

## 机制与边界

### Cache-Aside 的所有权

MySQL 是商户事实源，Redis 是可丢失的派生缓存。更新流程先提交数据库，再删除缓存；读取 miss 时从数据库重建。选择删除而不是直接更新缓存，是为了让写路径不必复制完整的缓存序列化和 TTL 逻辑。

### 并发不一致窗口

“更新数据库后删缓存”仍有失败窗口：

- 数据库更新成功、删除 Redis 失败：旧值会保留到 TTL。
- 读请求先 miss 并读取旧数据库值，随后更新提交并删缓存，最后旧读请求又把旧值写回：可能短暂产生脏缓存。

常见演进包括删除失败重试、消息/订阅驱动失效、延迟双删或基于版本的写入保护。选择取决于一致性要求，不能只背固定答案。

### 空值缓存的边界

- 2 分钟内新创建相同 ID 的数据可能仍被视为不存在，写路径需要主动删除空值。
- 攻击者持续请求大量不同随机 ID，仍会制造大量短期空 key；需要参数校验、限流或布隆过滤器。
- 正常 TTL 固定为 30 分钟，批量预热且同时过期时可能发生缓存雪崩；可加入随机抖动。

### 备用互斥锁实现的确定问题

当前主链未启用 `QueryWithMutex`，不能把它写成线上能力。代码本身还有这些风险：

- `ret` 是 nil 指针，却直接传给 `json.Unmarshal`，命中缓存时会反序列化失败。
- 无论是否拿到锁都注册 `defer unlock`；获取失败的请求可能删除其他请求持有的锁。
- 锁值固定为 `"1"`，删除时不校验所有者。
- 获取失败通过递归和固定 50 ms 休眠重试，没有最大次数，也没有在 sleep 前检查 context。

### 备用逻辑过期实现的取舍

逻辑过期允许请求先返回旧值，只有一个请求异步重建，适合“可接受短暂陈旧但不接受热点阻塞”的场景。当前实现未接入主查询，并且 goroutine 修改外层 `err` 变量、解锁使用原请求 context，存在数据竞争与请求取消后无法解锁的风险。

## 需要懂什么

- **缓存与事实源的职责不同。** MySQL 保存业务事实，Redis 保存可重建副本；删除缓存让写路径不必维护完整缓存值和 TTL，却不能自动保证两次操作同时成功。
- **穿透、击穿、雪崩的触发条件不同。** 不存在 ID 反复回源是穿透，热点 key 到期并发回源是击穿，大量 key 同时到期是雪崩；空值缓存主要处理第一类。
- **空字符串与 key 缺失必须结合错误区分。** GET 命中空字符串时 err 为 nil，缺失 key 时为 redis.Nil；只看返回字符串会把已缓存的不存在误判为缓存 miss。
- **先写库后删缓存仍可能读到旧值。** 删除失败会保留旧缓存，旧查询也可能晚于删除回填旧数据；可靠失效事件、版本保护等属于后续方案，不是当前能力。

## 需要记住什么

- 读取顺序：命中 JSON 返回 → 命中空字符串返回不存在 → key 缺失回源并写缓存。
- 写入顺序：MySQL 更新成功 → 删除该商户缓存 → 下次请求重建。
- TTL：正常值 30 分钟，空值 2 分钟；空值也可能短暂遮住后来新建的数据。
- 当前不能声称已解决热点击穿或强一致；备用锁/逻辑过期实现还存在正确性问题。
