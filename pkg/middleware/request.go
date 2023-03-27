package middleware

import (
	"bytes"
	"gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"time"
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

	// 记录日志
	logger.Infof("{\"path\": \"%s\", \"client_ip\": \"%s\", \"method\": \"%s\", \"status_code\": %d, \"user_agent\": \"%s\", \"request_id\": \"%s\", "+
		"\"response_headers\": %s, \"response_body\": %q, \"elapsed\": %f}",
		path,
		clientIP,
		method,
		statusCode,
		userAgent,
		requestID,
		responseHeaders,
		responseBody,
		elapsed.Seconds(),
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

	logger.Infof("{\"path\": \"%s\", \"client_ip\": \"%s\", \"method\": \"%s\", \"user_agent\": \"%s\", \"request_id\": \"%s\", "+
		"\"request_headers\": %s, \"request_body\": %q, \"request_form\": %s}",
		path,
		clientIP,
		method,
		userAgent,
		requestID,
		requestHeaders,
		requestBodyStr,
		requestForm,
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
