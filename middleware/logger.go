package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 读取请求体（用于记录入参）
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 重新设置请求体，因为读取后会被消耗
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 创建自定义的 ResponseWriter 来捕获响应体
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 结束时间
		duration := time.Since(start)

		// 获取响应体
		responseBody := blw.body.String()

		// 获取错误信息
		errors := make([]string, 0)
		for _, err := range c.Errors {
			errors = append(errors, err.Error())
		}

		// 构建完整的路径
		if query != "" {
			path = path + "?" + query
		}

		// 记录日志
		log.Printf("[GIN] %s | %3d | %13v | %15s | %-7s %s",
			time.Now().Format("2006/01/02 - 15:04:05"),
			c.Writer.Status(),
			duration,
			c.ClientIP(),
			c.Request.Method,
			path,
		)

		// 记录详细的信息
		log.Printf("📝 请求详情:")
		log.Printf("   🎯 方法: %s, 路径: %s", c.Request.Method, path)
		log.Printf("   📨 请求头: %v", c.Request.Header)

		// 记录请求体（入参）
		if len(requestBody) > 0 {
			log.Printf("   📥 请求体: %s", string(requestBody))
		} else {
			log.Printf("   📥 请求体: 空")
		}

		// 记录查询参数
		if len(c.Request.URL.Query()) > 0 {
			log.Printf("   🔍 查询参数: %v", c.Request.URL.Query())
		}

		// 记录响应信息
		log.Printf("   📤 响应状态: %d", c.Writer.Status())
		if len(responseBody) > 0 && len(responseBody) < 1000 { // 限制响应体长度
			log.Printf("   📦 响应体: %s", responseBody)
		} else if len(responseBody) >= 1000 {
			log.Printf("   📦 响应体: %s... (截断)", responseBody[:1000])
		}

		// 记录错误信息
		if len(errors) > 0 {
			log.Printf("   ❌ 错误信息: %v", errors)
		}

		// 如果是错误状态码，记录更详细的信息
		if c.Writer.Status() >= 400 {
			log.Printf("   ⚠️  错误请求详情:")
			log.Printf("     - 状态码: %d", c.Writer.Status())
			log.Printf("     - 错误数量: %d", len(errors))
			log.Printf("     - 处理时间: %v", duration)
		}

		log.Println("────────────────────────────────────────")
	}
}

// 自定义 ResponseWriter 来捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
