//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go-dianping/internal/base/cache_client"
	"go-dianping/internal/base/redis_worker"
	"go-dianping/internal/handler"
	"go-dianping/internal/server"
	"go-dianping/internal/service"
	"go-dianping/pkg/app"
	"go-dianping/pkg/log"
	"go-dianping/pkg/server/http"
)

var cacheClientSet = wire.NewSet(
	// 【启动阶段：缓存依赖】
	// 把创建店铺缓存客户端的函数交给 Wire 管理。
	cache_client.NewCacheClientForShop,
)

var redisWorkerSet = wire.NewSet(
	// 【启动阶段：Redis Worker 依赖】
	redis_worker.NewRedisWorker,
)

var serviceSet = wire.NewSet(
	// 【启动阶段：基础设施和业务 Service 依赖】
	// Wire 会先创建 DB/Redis，再把它们传给 Service 构造函数。
	service.NewDB,
	service.NewQuery,
	service.NewRedis,
	service.NewRedSync,

	service.NewService,
	service.NewBlogService,
	service.NewSeckillVoucherService,
	service.NewShopService,
	service.NewShopTypeService,
	service.NewUserService,
	service.NewVoucherService,
	service.NewVoucherOrderService,
)

var handlerSet = wire.NewSet(
	// 【启动阶段：Handler 依赖】
	// Handler 依赖 Service；因此 Service 必须先被创建。
	handler.NewHandler,
	handler.NewBlogHandler,
	handler.NewShopHandler,
	handler.NewShopTypeHandler,
	handler.NewUploadHandler,
	handler.NewUserHandler,
	handler.NewVoucherHandler,
	handler.NewVoucherOrderHandler,
)

var serverSet = wire.NewSet(
	// 【启动阶段：HTTP Server 依赖】
	// NewHTTPServer 依赖多个 Handler，Wire 会把它们注入进去。
	server.NewHTTPServer,
)

// build App
func newApp(
	httpServer *http.Server,
) *app.App {
	// 【启动阶段：应用容器】
	// 把已经配置好的 HTTP Server 放入 App。
	// App.Run 随后负责启动 Server，并等待退出信号。
	return app.NewApp(
		app.WithServer(httpServer),
		app.WithName("go-dianping"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	// 【启动阶段：依赖注入入口】
	// Wire 根据下面这些 Provider Set 生成 wire_gen.go。
	// 生成后的代码会按依赖关系依次创建 Redis、DB、Service、Handler 和 Server。
	panic(wire.Build(
		serverSet,
		cacheClientSet,
		redisWorkerSet,
		serviceSet,
		handlerSet,
		newApp,
	))
}
