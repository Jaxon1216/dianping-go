package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// 这个net/http 是啥呀，为什么在他后面加 // 会让本文件灰色字体一起改动？
	// 【answer】net/http 是 Go 标准库的 HTTP 包；这里使用 http.StatusNoContent 返回 204。// 只会注释本行，不会改变后续代码的语义，灰色通常是 IDE 的注释或未使用提示。
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.Header("Access-Control-Allow-Methods", c.GetHeader("Access-Control-Request-Method"))
			c.Header("Access-Control-Allow-Headers", c.GetHeader("Access-Control-Request-Headers"))
			c.Header("Access-Control-Max-Age", "7200")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
