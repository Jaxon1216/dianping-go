# C07 Feed 推送与滚动分页

> 课程性质：核心
> 返回：[主导学](../roadmap.md)；前置：[C06](C06-关注与点赞的数据一致性.md)、Redis ZSet；下一课：[C08 秒杀原子判定](C08-秒杀资格原子判定.md)
> 对应：B4；[Q6](../interview.md#q6-feed-为什么使用时间游标加同分值偏移)；目标掌握：看懂链路、理解机制、设计与表达。

## 本课重点

场景：作者发布笔记，粉丝在关注页按时间倒序滚动加载；两篇笔记可能拥有相同毫秒时间戳。

- Feed 的设计采用写扩散：发布时向粉丝收件箱写入博客 ID，读取时按时间倒序取候选。
- 滚动分页用时间戳与同分值偏移组成游标，时间戳单独不足以处理同毫秒记录。
- 当前推送目标字段和 offset=0 参数校验有实现问题，不能将设计机制表述为端到端可用保证。

## 主链与职责

设计意图是写扩散：

```text
作者发布博客
  -> MySQL 保存博客
  -> MySQL 查询作者的粉丝关系
  -> 对每个粉丝执行 ZADD feed:{followerId}
       member = blogId
       score  = publishTimeMillis

粉丝读取
  -> ZREVRANGEBYSCORE feed:{currentUser}
       max = 上次最小时间
       offset = 上页在该时间已消费数量
       count = 2
  -> 批量查询博客详情并恢复 ZSet 顺序
  -> 补作者与点赞态
  -> 返回 list + minTime + offset
```

前端保存服务端返回的 `minTime` 和 `offset`，下一次滚动时原样带回。

## 源码阅读

1. [`internal/service/blog.go`](../../../internal/service/blog.go)：按顺序阅读 `SaveBlog` 和 `QueryBlogOfFollow`。
2. [`internal/model/follow.gen.go`](../../../internal/model/follow.gen.go)：再次确认关注关系中发起者与被关注者字段。
3. [`api/v1/blog.go`](../../../api/v1/blog.go)：滚动结果的 `list/min_time/offset`。
4. [`internal/handler/blog.go`](../../../internal/handler/blog.go)：查询参数 `lastId/offset` 的绑定。
5. [`frontend/src/api/blog.js`](../../../frontend/src/api/blog.js)：把后端 `min_time` 兼容转换为 `minTime`。
6. [`frontend/src/views/profile/index.vue`](../../../frontend/src/views/profile/index.vue)：滚动参数的保存、下拉刷新与列表追加。
7. [`frontend/src/utils/useScrollLoad.js`](../../../frontend/src/utils/useScrollLoad.js)：加载中、结束和错误状态如何控制重复请求。

## 机制与边界

### 时间游标为什么还需要 offset

ZSet 的 score 不是唯一值。假设倒序数据为：

```text
1000: A, B, C
 999: D
```

每页取 2 条，第一页拿到 A、B，返回 `minTime=1000, offset=2`。第二页仍以 1000 为上界，但跳过该 score 已消费的 2 条，先拿 C，再拿 D。第二页最低 score 变成 999，所以返回新的 offset=1。

只用时间戳会在包含边界时重复，排除边界时漏掉 C。`max + offset` 是处理相同 score 的组合游标。

### 同分值跨多页的偏移边界

当前返回的 `os` 只统计本页最低 score 的条数，没有在最低 score 仍等于请求上界时累加传入 offset。例如 6 条记录同为 1000、每页 2 条：第一页 offset=2，第二页仍返回 2，第三页就会重复第二页。正确语义是同一边界时间戳下累计已消费数量，最低时间戳变化时再重置；当前不能宣称同分值分页已完全不重不漏。

### 确定 bug：推送目标字段错误

查询条件 `follow_user_id = 作者 ID` 得到的每一行中：

- `user_id` 是粉丝；
- `follow_user_id` 是作者。

当前循环却用 `follow_user_id` 组装 Feed key，所以博客被反复写入作者自己的收件箱，而不是每个粉丝的收件箱。ZSet member 相同，重复 ZAdd 最终通常只保留一条，但真实粉丝收不到内容。

正确目标应基于每行的 `user_id`，修复后需要用两个粉丝账号做端到端验证。

### 确定 bug：首个 offset 的参数校验

handler 给 `offset int64` 同时设置默认值 0 和 `binding:"required"`。Gin 使用的 validator 会把数字 0 视为 required 校验失败，因此首次请求的合法 `offset=0` 可能直接返回 400。这里应该允许 0，并单独校验 `offset >= 0`。

### 写扩散的失败边界

- 博客先写 MySQL，Feed 后推 Redis；任一推送失败时博客已存在，且此前粉丝可能已经收到，形成部分成功。
- 推送在请求 goroutine 中串行执行，粉丝数量越大，发布延迟和失败概率越高。
- 没有 outbox、重试任务或推送进度，无法从断点恢复。
- 取消关注只改关系，不会清理已进入收件箱的历史内容。
- 新关注不会自动补发作者旧内容，这是产品语义，不一定是 bug。

更可靠的演进是：博客事务提交时写 outbox，异步 worker 分批 fan-out；记录任务幂等键和进度，失败可重试。大 V 场景可混合写扩散与读扩散，避免一次发布写入海量收件箱。

### 读取链路的性能边界

读取后仍对每篇博客逐条查作者与点赞态，存在 C06 已分析的 N+1。Feed 游标解决的是分页稳定性，不解决详情组装成本。

## 需要懂什么

- **写扩散把读取成本转移到发布。** 普通用户读多写少时收件箱读取简单，但粉丝多会放大发布写入；大 V 的读时合并是替代方案，当前未实现混合策略。
- **同分值偏移补充时间游标。** 假设 score=1000 有 A/B/C，每页两条，第一页消费两条后下一页仍以 1000 为上界、跳过已消费数量，才不会遗漏剩余项。跨多页同分值时偏移应累计，不能只记本页数量。
- **同步推送可能部分成功。** 博客已经入库、只有部分收件箱写入成功时，返回错误不能撤销之前状态；outbox 与分批幂等重试属于可选演进。
- **分页机制与完整交付链路是不同层次。** 收件人字段或参数校验出错时，游标算法本身正确也不能让粉丝实际读到内容。

## 需要记住什么

- ZSet：member=blogId，score=发布时间毫秒；返回 list、minTime、offset。
- 下一页上界为上页 minTime，offset 表示该边界分值已消费的数量；同分值跨页需要累计。
- 关注关系中 user_id 才是粉丝；当前代码误用 follow_user_id 将消息写回作者收件箱。
- offset=0 是合法起点；当前 required 校验与此冲突。写扩散、补发历史和取消关注清理需分别定义产品语义。
