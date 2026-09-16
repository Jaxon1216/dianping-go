package main

import (
	"context"
	"flag"
	"fmt"

	"go-dianping/cmd/server/wire"
	"go-dianping/pkg/config"
	"go-dianping/pkg/log"

	"go.uber.org/zap"
)

// @title           Go Dian Ping
// @version         1.0.0
// @description     golang 实现的黑马点评
// @contact.name   llmons
// @contact.url    https://github.com/llmons
// @contact.email  llmons@foxmail.com
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
// @host      localhost:8080
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @externalDocs.description  GitHub
// @externalDocs.url          https://github.com/llmons/go-dianping
func main() {
	// 【启动阶段 1：读取启动参数】
	// flag.String 定义一个命令行参数：
	// - 参数名是 conf
	// - 默认值是 config/local.yml
	// - 最后一段文字是参数说明
	// 运行示例：go run ./cmd/server -conf config/prod.yml
	//
	// flag.String 返回的是 *string，也就是“指向字符串的指针”。
	// 因此这里的 envConf 不是配置文件内容，而是保存配置文件路径的指针。
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")

	// 【启动阶段 2：解析命令行参数】
	// flag.Parse 会读取程序启动时传入的参数，并把结果写回 envConf 指向的字符串。
	flag.Parse()

	// *envConf 表示取出指针指向的实际字符串，例如 config/local.yml。
	// NewConfig 读取这个 YAML 文件，创建后续代码使用的 Viper 配置对象。
	conf := config.NewConfig(*envConf)

	// 【启动阶段 3：创建日志组件】
	// 日志配置也来自 conf。之后启动过程和请求处理过程都使用 logger 记录日志。
	logger := log.NewLog(conf)

	// 【启动阶段 4：组装依赖】
	// wire.NewWire 会创建 Redis、MySQL、Service、Handler 和 HTTP Server。
	// app 是组装完成、可以运行的应用；cleanup 是退出时执行的清理函数；err 表示组装是否失败。
	app, cleanup, err := wire.NewWire(conf, logger)

	// defer 表示：main 函数结束前执行 cleanup。
	// 注意：它只在成功调用 NewWire 后才有实际清理意义；当前项目 cleanup 由 Wire 生成。
	defer cleanup()

	// 如果初始化 Redis、MySQL 或其他依赖失败，停止启动并打印错误。
	if err != nil {
		panic(err)
	}

	// 【启动阶段 5：打印启动信息】
	// 这里的 fmt.Sprintf 只负责拼接日志字符串，不会启动服务。
	logger.Info("server start", zap.String("host", fmt.Sprintf("http://%s:%d", conf.GetString("http.host"), conf.GetInt("http.port"))))
	logger.Info("docs addr", zap.String("addr", fmt.Sprintf("http://%s:%d/swagger/index.html", conf.GetString("http.host"), conf.GetInt("http.port"))))

	// 【启动阶段 6：进入应用运行阶段】
	// app.Run 会启动已经组装好的 HTTP Server。
	// 从这里开始，程序会等待客户端请求；请求到来后才进入 Router/Middleware/Handler。
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
