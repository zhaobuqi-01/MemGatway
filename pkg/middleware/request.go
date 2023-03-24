package middleware

import (
	"bytes"
	"fmt"
	"gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

func RequestOutLog(c *gin.Context) {
	// 获取请求开始时间
	start := time.Now()

	// 处理请求
	c.Next()

	// 获取请求结束时间和执行时间
	end := time.Now()
	latency := end.Sub(start)

	// 记录日志
	path := c.Request.URL.Path
	clientIP := c.ClientIP()
	method := c.Request.Method
	statusCode := c.Writer.Status()
	userAgent := c.Request.UserAgent()
	requestID := c.Writer.Header().Get("X-Request-Id")
	fmt.Printf("Out path:%s | clientIP:%s | method:%s |userAgent:%s | requestID:%s\n", path, clientIP, method, userAgent, requestID)

	logger.Info(fmt.Sprintf("%s _com_request_out [%s] %s | %3d | %13v | %15s | %-7s %s",
		path,
		end.Format("2006-01-02T15:04:05.000Z0700"),
		clientIP,
		statusCode,
		latency,
		method,
		userAgent,
	), zap.String("request_id", requestID))
}

func RequestInLog(c *gin.Context) {
	// 获取请求 body
	var requestBody []byte
	if c.Request.Body != nil {
		requestBody, _ = ioutil.ReadAll(c.Request.Body)
	}

	// 重新生成 request body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

	// 记录日志
	path := c.Request.URL.Path
	clientIP := c.ClientIP()
	method := c.Request.Method
	userAgent := c.Request.UserAgent()
	requestID := c.Request.Header.Get("X-Request-Id")
	fmt.Printf("In path:%s | clientIP:%s  | method:%s |   \n userAgent:%s  |  requestID:%s\n", path, clientIP, method, userAgent, requestID)

	logger.Info(fmt.Sprintf("%s _com_request_in [%s] %s | %-7s %s | %s",
		path,
		time.Now().Format("2006-01-02T15:04:05.000Z0700"),
		clientIP,
		method,
		userAgent,
	), zap.String("request_id", requestID))
	logger.Debug("request body", zap.ByteString("body", requestBody))
}

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理请求前记录请求信息
		RequestInLog(c)

		// 处理请求
		c.Next()

		// 处理请求后记录响应信息
		RequestOutLog(c)
	}
}
