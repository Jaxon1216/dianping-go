# C09 Stream 异步落库与事务

> 课程性质：核心
> 返回：[主导学](../roadmap.md)；前置：[C08](C08-秒杀资格原子判定.md)、MySQL 事务、消息确认；下一课：[C10 部署与观测](C10-部署测试与可观测验证.md)
> 对应：B6；[Q8](../interview.md#q8-异步消费者如何保证订单与数据库库存一起提交)；目标掌握：看懂链路、理解机制、设计与表达。

## 本课重点

场景：Lua 已把订单消息写入 Redis Stream，消费者准备写 MySQL；处理过程中可能发生锁竞争、库存耗尽、进程崩溃或重复投递。

- Stream 消费组将 Redis 准入与 MySQL 写入分开，数据库事务包含查重、条件扣库存与订单保存。
- 事务内必须使用 tx 查询对象；stock > 0 条件与受影响行数共同判断是否扣减成功。
- Pending/重放要求幂等；当前组初始化、确认控制、恢复循环和数据库组合唯一约束均需补齐。

## 主链与职责

```text
后台 goroutine
  -> XREADGROUP g1/c1 读取 stream.orders 的新消息 ">"
  -> 解析 id、userId、voucherId
  -> 按 userId 获取 RedSync 锁
  -> MySQL Transaction：
       查询同用户同券订单
       条件 UPDATE stock = stock - 1 WHERE stock > 0
       INSERT voucher_order
  -> 正常处理成功后 XACK；失败分支还需核对确认控制

处理失败
  -> 读取同一消费者 Pending List，ID "0"
  -> 重试处理
  -> 成功后 ACK
```

事务保证的是“数据库库存更新与订单插入一起提交或一起回滚”。Stream 的确认与 Redis 预扣在事务之外，因此整个系统仍是跨存储最终一致问题。

## 源码阅读

1. [`internal/service/voucher_order.go`](../../../internal/service/voucher_order.go)：从构造器中的 goroutine 读到 `handleVoucherOrder`、`handlePendingList` 和 `createVoucherOrder`。
2. [`internal/query/gen.go`](../../../internal/query/gen.go)：`Transaction` 把事务 DB clone 成新的 `*query.Query`。
3. [`internal/model/voucher_order.gen.go`](../../../internal/model/voucher_order.gen.go)：订单字段和主键类型。
4. [`internal/model/seckill_voucher.gen.go`](../../../internal/model/seckill_voucher.gen.go)：数据库库存字段。
5. [`deploy/docker-compose/hmdp.sql`](../../../deploy/docker-compose/hmdp.sql)：订单表只有 `id` 主键，没有用户与券的组合唯一约束。
6. [`internal/service/service.go`](../../../internal/service/service.go)：RedSync 当前只使用一个 Redis client pool。

事务内的关键结构是：

```go
return s.query.Transaction(func(tx *query.Query) error {
    count, err := tx.VoucherOrder.Where(...).Count()
    info, err := tx.SeckillVoucher.
        Where(...Stock.Gt(0)).
        Update(...Stock.Sub(1))
    return tx.VoucherOrder.Save(voucherOrder)
})
```

回调返回 error 触发回滚，返回 nil 触发提交。使用 `tx` 才能让三个 SQL 共用同一事务连接。

## 机制与边界

### 条件更新解决数据库超卖

`UPDATE ... SET stock = stock - 1 WHERE stock > 0` 把“检查和扣减”交给单条 SQL 原子完成。`RowsAffected == 0` 代表没有符合条件的库存，即使 SQL 本身没有报错，也必须按库存不足处理。

这比“先 SELECT stock，再 UPDATE”更能抵抗并发，但只保护 MySQL 库存不为负，不自动让 Redis 预扣与 MySQL 一致。

### ACK 与 Pending 的语义

消费组读取后，消息进入 Pending Entries List；XACK 会从 Pending 中移除记录。正确策略应是只有业务落库成功后才 ACK，Redis 并不替应用校验这个前提。进程在落库前崩溃，消息应留在 Pending，重启后可恢复。这个模型通常是 at-least-once，因此处理逻辑必须幂等。

“恰好一次”不能仅靠 Redis Stream 宣称。消息可能重复投递，数据库提交成功但 ACK 前崩溃就是经典重放场景。

主循环处理失败后调用 `handlePendingList()`，返回后仍继续执行 XACK，没有显式成功分支；不能把当前代码概括为“严格保证成功后才确认”。恢复循环、重复订单结果与确认控制需要一起收敛。

### 当前幂等边界

- 应用层查询“同用户同券是否已有订单”。
- 用户维度 RedSync 锁降低同一用户并发处理。
- 订单 ID 主键可拒绝同 ID 重复插入。

但订单表没有 `(user_id, voucher_id)` 唯一索引。若锁失效、Redis 数据丢失或不同消息带不同订单 ID 并发到达，两个事务都可能先查到 0 再插入。生产级兜底应添加组合唯一约束，并把重复键当作幂等成功或明确业务结果处理。

当前 RedSync 只接入一个 Redis 实例，实际是带续期/所有权语义的 Redis 锁用法，不是多独立 Redis 节点上的 Redlock quorum。

### 消费者生命周期和恢复的确定问题

- 源码没有 `XGROUP CREATE ... MKSTREAM`。消费组不存在时，`XREADGROUP` 持续返回 `NOGROUP`，循环没有日志和退避，可能形成错误热循环。
- 所有进程都使用固定 consumer 名 `c1`，不利于多实例归属、Pending 转移和观测。
- 消息字段使用直接类型断言；畸形消息可能 panic，解析失败消息也不会 ACK 或进入死信。
- 主循环只有业务处理失败时才进入 Pending 恢复，服务重启后不会主动先扫描历史 Pending。
- Pending 读取遇到 `redis.Nil` 时直接 `continue`，到队列清空后可能无法跳出恢复循环，阻塞后续新消息消费。
- 没有最大重试次数、指数退避、死信队列和人工修复入口；永久业务错误可能成为 poison message。
- worker 使用 `context.Background()` 和无限循环，应用退出时不会优雅停止。

### 跨存储补偿

Lua 已扣 Redis 库存并记录用户，但 MySQL 事务可能永久失败。需要对账任务比较 Redis 资格、Stream 状态和 MySQL 订单，决定重试落库还是回补 Redis 库存/移除用户资格。当前仓库没有这条补偿链。

## 需要懂什么

- **事务保证的是同一数据库内共同提交。** 所有 SQL 使用 tx 才共享事务连接；使用外层共享 query 会让部分操作脱离事务，Redis 预扣与 Stream ACK 本来也不在该事务中。
- **条件更新把检查与扣减交给一条 SQL。** WHERE stock > 0 避免两个请求先读到相同库存再分别扣减；RowsAffected=0 虽未必有 SQL error，业务仍须视为库存不足。
- **提交后确认仍可能重放。** 数据库提交成功、ACK 前崩溃时消息还在 Pending；恢复后必须识别已完成，不能再次扣库存。当前确认分支不够严格，不能概括为完整可靠投递。
- **锁、应用查重与唯一索引的作用不同。** 锁减少同时处理，查重提供业务判断，唯一索引在最终存储守住不变量；当前没有用户与券组合唯一索引，不能保证不同订单 ID 下绝无重复。

## 需要记住什么

- 正常主链：XREADGROUP → 用户锁 → tx 查重 → 条件扣库存 → 保存订单 → ACK。
- 同一 tx 保证 MySQL 库存与订单共同提交；不覆盖 Redis 资格和库存。
- 消费组需要初始化；读新消息使用 >，恢复本消费者 Pending 使用 0。
- 当前没有完整恢复闭环：固定 consumer、启动不扫旧 Pending、无最大重试/死信/补偿，ACK 还需受业务成功控制。
