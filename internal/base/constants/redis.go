package constants

import "time"

// RedisKey 和 RedisTTL 是为 Redis 常量预留的语义类型。
// 当前文件中的 key 常量仍由 Go 推断为 string，TTL 常量实际使用 time.Duration，
// 因此这两个自定义类型目前没有参与业务代码。
type RedisKey string
type RedisTTL int

const (
	// 登录验证码使用 Redis String。
	// 最终 key 形如 login:code:13800138000，value 是验证码。
	RedisLoginCodeKey = "login:code:"
	RedisLoginCodeTTL = time.Minute * 2

	// 登录会话使用 Redis Hash。
	// 最终 key 形如 login:token:<token>，Hash 保存精简后的用户信息。
	RedisLoginUserKey = "login:token:"
	RedisLoginUserTTL = time.Minute * 30

	// 店铺详情缓存使用 Redis String，value 是序列化后的店铺数据。
	// 最终 key 形如 cache:shop:1。
	RedisCacheShopKey = "cache:shop:"
	RedisCacheShopTTL = time.Minute * 30

	// 店铺类型列表使用一个 Redis List，因此这里是完整 key，不需要再拼接 ID。
	RedisCacheShopTypeKey = "cache:shop-type"
	RedisCacheShopTypeTTL = time.Minute * 30

	// 数据库中也不存在目标数据时，短暂缓存空字符串，防止缓存穿透。
	RedisCacheNullTTL = time.Minute * 2

	// 店铺缓存重建锁的 key 前缀，最终 key 形如 redis_lock:shop:1。
	RedisLockShopKey = "redis_lock:shop:"

	// 秒杀库存使用 Redis String，key 后拼接优惠券 ID，value 是剩余库存。
	RedisSeckillStockKey = "seckill:stock:"

	// 笔记点赞使用 Sorted Set：member 是用户 ID，score 是点赞时间。
	RedisBlogLikeKey = "blog:liked:"
	// 关注流收件箱使用 Sorted Set：member 是笔记 ID，score 是发布时间。
	RedisBlogFeedKey = "feed:"

	// 店铺位置使用 Redis GEO，key 后拼接店铺类型 ID。
	RedisShopGeoKey = "shop:geo:"

	// 用户签到使用 Bitmap。
	// 最终 key 形如 user:sign:<userID>:<yyyyMM>，第 day-1 位表示当日是否签到。
	RedisUserSignKey = "user:sign:"
)
