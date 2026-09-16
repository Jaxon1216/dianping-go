package server

// 不太懂什么意思，这里是把本文件export成 名叫server的 package么？
// 【answer】不是 export；它表示当前文件属于 server 包。Go 中标识符首字母大写才表示可以被其他包使用，例如 NewHTTPServer。

import (
	"go-dianping/docs"
	"go-dianping/internal/handler"
	"go-dianping/internal/middleware"
	"go-dianping/pkg/log"
	"go-dianping/pkg/server/http"
	netHttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 这里的依赖，一部分GitHub开头的是包依赖（也就是类似前端package里写好，npm install后 随处可import的？）
// 【answer】基本可以这样类比：第三方依赖版本由 go.mod 管理，但当前 Go 文件仍需要显式 import 后才能使用。go-dianping 开头的是项目内部包。

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	rdb *redis.Client,
	blogHandler *handler.BlogHandler,
	followHandler *handler.FollowHandler,
	shopHandler *handler.ShopHandler,
	shopTypeHandler *handler.ShopTypeHandler,
	uploadHandler *handler.UploadHandler,
	userHandler *handler.UserHandler,
	voucherHandler *handler.VoucherHandler,
	voucherOrderHandler *handler.VoucherOrderHandler,
	// 入参都是地址是么？*标识指针类型？ * 左侧空一格的是形参别名，那*右边的是什么？
	// 【answer】这些参数的类型确实都是指针类型；左边是参数名，右边是类型，* 是类型的一部分，例如 logger 的类型是 *log.Logger。
) *http.Server {
	// 【启动阶段：创建并配置 HTTP/Gin 服务】
	// 这个函数运行在程序启动时，不是每个请求到来时运行。
	// 它负责创建 Gin Engine、注册全局中间件和注册所有路由。
	// conf viper是干啥的呀，看不懂这个函数在做什么
	// 【answer】conf 是 Viper 配置对象，用 GetString/GetInt 读取 config/local.yml 或 config/prod.yml 中的配置。
	if conf.GetString("env") == "prod" {
		// 这里的灰色的 key：怎么复制不了？只复制出了一个 "env"
		// 【answer】"env" 是传给 Viper 的普通字符串配置键，对应 YAML 中的 env: local 或 env: prod。
		gin.SetMode(gin.ReleaseMode)
		// 这个方法需要记住吗？
		// 【answer】当前只需记住作用：生产环境切换到 Gin 发布模式，减少调试输出；API 细节以后需要时再查。
	}
	// 初始化实例吗？点进去发现看不懂呀怎么办
	// 【answer】对，gin.Default() 创建 Gin 路由引擎，点进去发现会默认挂载 Logger 和 Recovery 中间件。
	g := gin.Default()

	// 全文件搜了一下，没搜到这个属性，不知道哪来的，也不知道干啥的
	// 【answer】这是 g 的类型 *gin.Engine 上定义的字段；允许请求路径末尾多一个 / 时自动重定向到已注册的路径。
	g.RedirectTrailingSlash = true

	// 不懂为什么传这些参数
	// 【answer】g 是 Gin 引擎，logger 是日志对象，WithServerHost/Port 是配置自定义 HTTP Server 监听地址和端口的选项。
	s := http.NewServer(
		g,
		// 为什么这里engine是黑色，并且复制不了？
		// 【answer】engine 是 NewServer 函数定义中的参数名；调用函数时只传入 g 这个表达式，不需要写参数名。
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	// swagger doc
	docs.SwaggerInfo.BasePath = "/v1" // 这一行是在干啥呀？
	// 【answer】docs.SwaggerInfo 是生成的 Swagger 配置；这行设置接口文档声明的基础路径为 /v1。
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	s.Use(
		// import 进来的跨域中间件
		// 【answer】import 只是引入包；这里实际调用函数得到中间件，并通过 Use 注册为全局中间件。
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		middleware.RefreshToken(rdb),
	)

	// 【启动阶段：全局 Middleware 注册完成】
	// 这里先保存中间件执行规则；真正的请求到来后，Gin 才按顺序执行它们。
	// 以上 Use 注册的是全局中间件：所有后续注册的路由都会先经过它们。
	// 下面的 Group 只负责给路由拼接公共路径；是否需要登录，要看该 Group
	// 是否额外调用 Use(middleware.Login())。

	// ========== router ==========
	// 看不太懂这个get的封装，看起来没有return啥东西？
	// 【answer】s.GET 是注册路由：收到 GET / 时执行下面的匿名函数。响应由 ctx.String 写出，调用处不需要使用它的返回值。
	s.GET("/", func(ctx *gin.Context) {
		ctx.String(netHttp.StatusOK, "OK")
	})
	{
		blogRouter := s.Group("/blog")
		// 下面这个花括号有啥说法么？
		// 【answer】这是 Go 的普通代码块，用于按业务分组并限制局部变量作用域，不是 Gin 专用语法，通常也可以省略。
		{
			blogRouter.POST("", blogHandler.SaveBlog)
			blogRouter.PUT("/like/:id", blogHandler.LikeBlog)
			blogRouter.GET("/of/me", blogHandler.QueryMyBlog)
			blogRouter.GET("/hot", blogHandler.QueryHotBlog)
			blogRouter.GET("/:id", blogHandler.QueryById)
			blogRouter.GET("/of/follow", blogHandler.QueryBlogOfFollow)
		}

		followRouter := s.Group("/follow")
		{
			followRouter.PUT("/:id/:isFollow", followHandler.Follow)
			followRouter.GET("/or/not/:id", followHandler.IsFollow)
			followRouter.GET("/common/:id", followHandler.FollowCommons)
		}

		shopRouter := s.Group("/shop")
		{
			shopRouter.GET("/:id", shopHandler.QueryShopById)
			shopRouter.PUT("", shopHandler.UpdateShop)
		}

		shopTypeRouter := s.Group("/shop-type")
		{
			shopTypeRouter.GET("/list", shopTypeHandler.QueryTypeList)
		}

		uploadRouter := s.Group("/upload")
		{
			uploadRouter.POST("/blog", uploadHandler.UploadImage)
			uploadRouter.GET("/blog/delete", uploadHandler.DeleteBlogImg)
		}

		userRouter := s.Group("/user")
		{
			// 【启动阶段：Router 注册】
			// Group 返回一个路由分组对象，不会处理请求，也不会返回 HTTP 响应。
			// 这里的公共前缀是 /user。
			// userRouter.Group("/") 的返回值仍然是一个路由分组对象，
			// 这里只是把共同前缀 /user 继续保留下来。
			noAuthRouter := userRouter.Group("/")
			{
				// 不带 Login 中间件：登录前也可以访问发送验证码和登录接口。
				// 这两行是在注册“方法 + 路径 + Handler”的对应关系：
				// POST /user/code  -> userHandler.SendCode
				// POST /user/login -> userHandler.Login
				noAuthRouter.POST("/code", userHandler.SendCode)
				noAuthRouter.POST("/login", userHandler.Login)
			}

			// Use(middleware.Login()) 给这个分组追加局部中间件。
			// 因此 /user/me 请求会先检查当前请求上下文里有没有用户。
			authRouter := userRouter.Group("/").Use(middleware.Login())
			{
				authRouter.GET("/me", userHandler.Me)
			}
		}

		voucherRouter := s.Group("/voucher")
		{
			voucherRouter.POST("/seckill", voucherHandler.AddSeckillVoucher)
			voucherRouter.POST("", voucherHandler.AddVoucher)
			voucherRouter.GET("/list/:shopId", voucherHandler.QueryVoucherOfShop)
		}

		voucherOrderRouter := s.Group("/voucher-order").Use(middleware.Login())
		{
			voucherOrderRouter.POST("/seckill/:id", voucherOrderHandler.SeckillVoucher)
		}
	}

	return s
}
