// test/middleware_test.go

package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gateway/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func TestRequestLog(t *testing.T) {
	// 创建一个gin引擎
	router := gin.New()

	// 添加中间件
	router.Use(middleware.RequestLog())

	// 注册GET路由
	router.GET("/test", func(c *gin.Context) {
		// 设置响应头
		c.Writer.Header().Set("X-Request-Id", "12345")

		// 设置响应内容
		c.String(http.StatusOK, "OK")
	})

	// 注册POST路由
	router.POST("/test", func(c *gin.Context) {
		// 获取请求body
		bodyBytes := []byte("username=test&password=123456")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// 设置请求头
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request.Header.Set("User-Agent", "Test-Agent")

		// 设置响应头
		c.Writer.Header().Set("X-Request-Id", "12345")
		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.Writer.Header().Set("X-Custom-Header", "CustomValue")

		// 设置响应内容
		c.String(http.StatusOK, "OK")
	})

	// 发送GET请求
	req1, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	// 发送POST请求
	req2, err := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer([]byte("username=test&password=123456")))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req2.Header.Set("User-Agent", "Test-Agent")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// 打印响应结果
	fmt.Printf("GET response: %s\n", w1.Body.String())
	fmt.Printf("POST response: %s\n", w2.Body.String())
}
