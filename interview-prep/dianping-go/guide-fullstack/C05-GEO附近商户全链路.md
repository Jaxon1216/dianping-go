# C05 GEO 附近商户全链路

> 课程性质：核心
> 返回：[主导学](../导学.md)；前置：[C02](C02-前后端HTTP契约.md)、分页、经纬度基础；下一课：[C06 关注与点赞](C06-关注与点赞的数据一致性.md)
> 对应：B3；[Q4](../面经.md#q4-附近商户为什么先查-redis-geo-再回到-mysql)；目标掌握：看懂链路、理解机制、设计与表达。

## 本课重点

场景：用户打开地图或商户列表，客户端携带当前坐标查询某一类型 5 公里内的商户。

- GEO 保存商户 ID 与坐标，先筛空间候选，再由 MySQL 补全详情。
- SQL IN 不保证输入顺序；当前代码显式恢复 Redis 返回顺序并关联距离。
- 坐标约定是 x=经度、y=纬度；当前初始化测试、距离排序和定位后刷新候选集均有边界。

## 主链与职责

```text
浏览器位置 / 默认坐标
  -> GET /shop/of/type?typeId=...&current=...&x=经度&y=纬度
  -> Redis GEOSEARCH shop:geo:{typeId}
       中心点 + 5000m + Count(pageEnd) + distance
  -> 截取当前页，得到商户 ID 与距离
  -> MySQL WHERE id IN (...)
  -> FIELD(id, ...) 恢复 Redis 返回顺序
  -> 把 distance 填入非持久化字段
  -> Vue 列表/地图渲染

Redis 索引为空
  -> 回退到按类型查询 MySQL，不返回距离
```

Redis 负责“空间候选集”，MySQL 负责“业务详情事实”。这是索引与事实源分离，而不是把整张商户表搬进 GEO。

## 源码阅读

1. [`frontend/src/views/shop-list/index.vue`](../../../frontend/src/views/shop-list/index.vue)：列表页传入默认坐标、普通页码和排序条件。
2. [`frontend/src/views/map/index.vue`](../../../frontend/src/views/map/index.vue)：地图先按默认位置加载，再请求浏览器定位，并在前端重新计算已取回商户的距离。
3. [`frontend/src/api/shop.js`](../../../frontend/src/api/shop.js)：查询参数如何进入 Axios。
4. [`internal/handler/shop.go`](../../../internal/handler/shop.go)：`x/y` 的绑定与文档注释。
5. [`internal/service/shop.go`](../../../internal/service/shop.go)：GEO 查询、分页截取、ID 顺序恢复与 MySQL 回退。
6. [`internal/model/shop.gen.go`](../../../internal/model/shop.gen.go)：数据库中 `x` 是经度、`y` 是纬度，`Distance` 是 `gorm:"-"` 的响应字段。
7. [`test/app/app_test.go`](../../../test/app/app_test.go)：`TestLoadShopData` 是当前唯一 GEO 预热入口，且需要先修正坐标与错误处理。

## 机制与边界

### 坐标契约必须统一

数据库和业务代码的实际约定是：

```text
x = longitude，经度，约 120
y = latitude，纬度，约 30
```

service 也按 `Longitude: x, Latitude: y` 调用 Redis。handler 的 Swagger 注释却把 x 写成纬度、y 写成经度，属于文档错误。地图与列表实际传参遵守“x 经度、y 纬度”。

### 为什么 MySQL 查询还要恢复顺序

SQL 的 `IN` 只限制集合，不承诺按输入 ID 排序。代码用 GORM Gen 的 `FIELD(ids...)` 生成数据库排序表达式，使返回列表与 Redis 候选顺序一致，再按 ID 填充对应距离。这个步骤是“索引结果和详情结果合并”的关键。

### 当前 GEO 索引初始化存在确定问题

`TestLoadShopData` 把 `shop.Y` 写成 Longitude、`shop.X` 写成 Latitude，正好交换了经纬度。纬度 120 超出合法范围，Redis 写入应失败，但测试没有检查 `GeoAdd` 的 error。测试里的分组初始化还错误地使用了 `shop.ID` 作为 key，虽然后续 append 到类型 key 仍可能产生数据，但会多出无用分组。

因此，仓库展示了 GEO 查询机制，但不能把“索引已可靠初始化”当成已验证事实。

### 排序与分页边界

- service 没有显式设置 GEO `Sort: "ASC"`。若要求严格按距离排序，应显式指定并通过 Redis 结果验证，不能依赖默认行为。
- 当 `sortBy` 非空时直接走 MySQL 评论数/评分排序，即使提供坐标也不走 GEO。这是明确的产品优先级选择。
- GEO 为空时回退 MySQL，可以提高可用性，但返回结果不再限制 5 公里且距离为 0；客户端需要知道这是降级语义。
- 页面分页先向 Redis 请求前 `current * pageSize` 条，再在应用层丢弃前页，深分页成本随页码增长；附近列表通常更适合游标或限制最大页数。

### 地图页面的端到端边界

地图在浏览器真实定位完成前就用默认坐标向后端加载每类第一页。定位成功后只重新计算和排序已加载的商户，没有重新请求真实位置附近的候选集。因此“地图显示距我最近”目前只对默认位置附近取回的有限集合成立。

## 需要懂什么

- **空间索引与详情存储分工。** GEO 适合筛选 ID 与距离，MySQL 保存经常变化的业务字段；稳定 ID 把两者关联起来，无需将完整 JSON 写进 member。
- **集合过滤不等于顺序保持。** IN 只限定哪些行符合条件，必须用 FIELD 或应用层按 ID 重组，才能保持空间索引返回的顺序。当前 GEO 未显式指定 ASC，恢复顺序不等于保证严格距离升序。
- **降级会改变结果语义。** 当前空候选集回退普通类型查询，结果不再保证位于 5 公里内，也没有真实距离；这与 Redis 故障一定可降级是不同主张。
- **前端重新算距离不能补回未查询的商户。** 定位成功后仅对已加载候选重算距离，不能保证新位置附近商户已进入结果集。

## 需要记住什么

- 主链：坐标/类型 → GEOSEARCH 候选 ID 和距离 → MySQL 批量详情 → 恢复候选顺序 → 返回列表。
- x=longitude，y=latitude；查询半径为 5000 米。
- 当前深分页先取 current × pageSize 条再截取，成本随页码增长。
- 索引预热代码交换了经纬度且忽略写入错误；当前不能把索引初始化或地图真实定位链路说成已闭环。
