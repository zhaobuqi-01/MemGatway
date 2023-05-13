package middleware

import (
	"bytes"
	"gateway/pkg/log"

	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestOutLog(c *gin.Context, responseTime time.Duration) {
	// after request
	// 记录响应信息
	uri := c.Request.RequestURI
	ip := c.ClientIP()
	method := c.Request.Method
	statusCode := c.Writer.Status()
	args := c.Request.PostForm
	response := c.GetString("response")

	log.Info("Response Info",
		zap.String("uri", uri),
		zap.String("client_ip", ip),
		zap.String("method", method),
		zap.Int("status_code", statusCode),
		zap.Any("args", args),
		zap.Any("response", response),
		zap.Duration("response_time", responseTime),
		zap.String("trace_id", c.GetString("TraceID")),
	)
}

func RequestInLog(c *gin.Context) {
	// 获取请求 body
	var requestBody []byte
	if c.Request.Body != nil {
		requestBody, _ = ioutil.ReadAll(c.Request.Body)
	}

	// 重新生成 request body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

	// 记录请求信息
	uri := c.Request.RequestURI
	clientIP := c.ClientIP()
	method := c.Request.Method
	requestBodyStr := string(requestBody)

	log.Info("Request Info",
		zap.String("uri", uri),
		zap.String("client_ip", clientIP),
		zap.String("method", method),
		zap.String("reques_body", requestBodyStr),
		zap.String("trace_id", c.GetString("TraceID")),
	)
}

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理请求前记录请求信息和请求进入时间
		startTime := time.Now()
		RequestInLog(c)

		// 处理请求
		c.Next()

		responseTime := time.Since(startTime)
		// 处理请求后记录响应信息和请求执行时间
		RequestOutLog(c, responseTime)
	}
}
