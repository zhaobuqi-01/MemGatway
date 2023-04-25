package middleware

import (
	"bytes"
	"gateway/pkg/log"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestOutLog(c *gin.Context, startTime time.Time) {
	// after request
	// 记录响应信息
	path := c.Request.URL.Path
	clientIP := c.ClientIP()
	method := c.Request.Method
	statusCode := c.Writer.Status()
	userAgent := c.Request.UserAgent()
	requestID := c.Writer.Header().Get("X-Request-Id")
	responseHeaders := c.Writer.Header()
	responseBody := c.GetString("response")

	// 计算请求执行时间
	elapsed := time.Since(startTime)

	log.Info("Request Info",
		zap.String("path", path),
		zap.String("client_ip", clientIP),
		zap.String("method", method),
		zap.Int("status_code", statusCode),
		zap.String("user_agent", userAgent),
		zap.String("request_id", requestID),
		zap.Any("request_headers", responseHeaders),
		zap.Any("request_body", responseBody),
		zap.Duration("elapsed", elapsed),
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
	path := c.Request.URL.Path
	clientIP := c.ClientIP()
	method := c.Request.Method
	userAgent := c.Request.UserAgent()
	requestID := c.Request.Header.Get("X-Request-Id")
	requestHeaders := c.Request.Header
	requestBodyStr := string(requestBody)
	requestForm := c.Request.PostForm

	log.Info("Response Info",
		zap.String("path", path),
		zap.String("client_ip", clientIP),
		zap.String("method", method),
		zap.String("user_agent", userAgent),
		zap.String("request_id", requestID),
		zap.Any("response_headers", requestHeaders),
		zap.String("response_body", requestBodyStr),
		zap.Any("request_form", requestForm),
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

		// 处理请求后记录响应信息和请求执行时间
		RequestOutLog(c, startTime)
	}
}
