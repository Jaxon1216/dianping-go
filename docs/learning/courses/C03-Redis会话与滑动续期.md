# C03 Redis 会话与滑动续期

> 课程性质：核心
> 返回：[主导学](../roadmap.md)；前置：[C02](C02-前后端HTTP契约.md)、Redis Hash/TTL、`context.Context`；下一课：[C04 Cache-Aside](C04-CacheAside与缓存一致性.md)
> 对应：B2；[Q2](../interview.md#q2-为什么把登录态放在-redis而不是进程内存)；目标掌握：看懂链路、理解机制、设计与表达。

## 本课重点

场景：用户请求验证码，完成登录，30 分钟内持续访问多个需要登录的接口，然后主动退出。

- token 是随机会话索引，用户信息与失效控制保存在 Redis，当前方案不使用 JWT 验签。
- 全局中间件恢复可选身份并续期，路由级中间件决定是否必须登录。
- 会话在最后一次成功恢复身份后 30 分钟失效，退出通过删除 Redis 会话主动撤销。

## 主链与职责

```text
发送验证码
  -> 校验手机号
  -> Redis String login:code:{phone}，TTL 2 分钟

登录
  -> Redis 校验验证码
  -> MySQL 按手机号查找或创建用户
  -> 生成 UUID token
  -> Redis Hash login:token:{token} 保存精简用户，TTL 30 分钟
  -> 前端 sessionStorage 保存 token

后续请求
  -> Axios 附加 Authorization: Bearer ...
  -> 全局 RefreshToken 查询 Redis Hash
  -> 用户写入 request context，刷新 30 分钟 TTL
  -> 路由级 Login 检查 context 中是否有用户
```

它是“随机 token + 服务端会话”，不是 JWT。token 只充当不可预测的查找键，用户状态和失效控制在 Redis。

## 源码阅读

1. [`frontend/src/views/login/index.vue`](../../../frontend/src/views/login/index.vue)：验证码发送、倒计时和登录后获取用户。
2. [`frontend/src/stores/user.js`](../../../frontend/src/stores/user.js)：token 放在 `sessionStorage`，刷新页面后重新调用 `/user/me`。
3. [`internal/service/user.go`](../../../internal/service/user.go)：`SendCode`、`Login`、`Logout` 的 Redis/MySQL 顺序。
4. [`internal/middleware/refresh_token.go`](../../../internal/middleware/refresh_token.go)：全局会话恢复与滑动续期。
5. [`internal/middleware/login.go`](../../../internal/middleware/login.go)：路由级强制鉴权。
6. [`internal/base/user_holder/user_holder.go`](../../../internal/base/user_holder/user_holder.go)：类型安全的私有 context key。
7. [`internal/base/constants/redis.go`](../../../internal/base/constants/redis.go)：验证码 2 分钟、会话 30 分钟。
8. [`internal/server/http.go`](../../../internal/server/http.go)：全局恢复与局部拦截的挂载范围。

## 机制与边界

### 为什么使用 Redis 会话

- 多个 Go 实例共享同一会话数据，不依赖负载均衡粘性。
- 删除 Redis key 即可主动失效，服务端可以控制 TTL。
- Hash 只保存精简用户，避免每次请求查询 MySQL，也不把手机号等非必要字段传播到业务层。

代价是每个携带 token 的请求多一次 Redis 读取和一次续期写入；Redis 成为登录链路的关键依赖。

### 为什么拆成两个中间件

全局恢复中间件对“游客可访问但登录后需要个性化”的接口也有用，例如博客详情可以返回当前用户是否点赞。路由级登录中间件只做“当前路径必须有用户”的判定。拆分后，公开路由不会因缺 token 被拦截，但仍能识别可选身份。

### 滑动续期语义

每次成功恢复用户后把 TTL 重置为 30 分钟，所以这是“距最后活跃 30 分钟失效”，不是“登录后固定 30 分钟失效”。若业务还要求最长会话时长，需要额外保存登录时间或绝对过期时间。

### 已确认风险

- 验证码 key 不存在时，Redis 返回 `redis.Nil`，当前登录逻辑把它当普通内部错误，最终常映射为 500，而不是“验证码错误”。
- 登录成功后没有删除验证码，同一验证码在 2 分钟 TTL 内可重复换取多个 token。
- 前端“密码登录”仍调用同一接口，但后端 `code` 是必填且不校验 `password`，所以该页面与后端契约不匹配。
- 刷新 TTL 失败时，中间件直接 `return`，既没有 `ctx.Next()` 也没有明确错误响应，可能留下空响应。
- token 存放在 `sessionStorage`，同源脚本可以读取；发生 XSS 时凭据可被窃取。改成 HttpOnly Cookie 会改变 CSRF 与跨域策略，需要整体权衡。
- 会话 Hash 是登录时快照。用户改昵称或头像后，已有 token 中的精简用户不会自动更新。

### Bitmap 签到是同一身份链路的扩展

`Sign` 从 request context 取当前用户，把一个月映射成一个 key、一天映射成一位；`SignCount` 取当月前 N 位并从低位连续统计 1。它说明身份恢复完成后，service 可以只依赖标准 context，不需要再次解析 HTTP header。

## 需要懂什么

- **共享会话把登录状态从进程内移到 Redis。** 请求到达任一 Go 实例都可恢复同一用户，不需要粘性路由；代价是携带 token 的请求依赖 Redis 读取与续期写入。
- **身份恢复和强制鉴权解决不同问题。** 公开接口也可能需要识别点赞态等个性化信息，所以先恢复身份，再按路由决定是否拒绝游客。
- **滑动过期不是固定生命周期。** 每次恢复用户重置 TTL，持续活跃可持续续期；如需最长会话时长，必须另存绝对过期信息，当前没有该机制。
- **客户端存储策略影响安全边界。** sessionStorage 可被同源脚本读取；HttpOnly Cookie 能减少脚本直接读取，却会改变 CSRF 与跨域策略，不能只替换存储位置。

## 需要记住什么

- 验证码：Redis String，TTL 2 分钟；会话：Redis Hash，TTL 30 分钟。
- 会话主链：验证码校验 → 查找/创建用户 → 随机 token → Hash 精简用户 → 请求 context → 路由鉴权。
- token 不是用户数据；本地有字符串不代表 Redis 会话仍存在。
- 当前验证码成功使用后未删除，Hash 是用户资料快照；原子消费验证码、资料同步和最长时长均需另行设计。
