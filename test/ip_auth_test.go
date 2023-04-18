package test

import (
	"gateway/pkg/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestIPAuthMiddleware(t *testing.T) {
	// 设置允许的 IP 列表
	viper.Set("config.server.allow_ip", []string{"127.0.0.1"})

	router := gin.New()
	router.Use(middleware.IPAuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	// 创建一个请求 127.0.0.1，它在允许的 IP 列表中
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"message":"success"}`, w.Body.String())

	// 创建一个不在允许的 IP 列表中的请求
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 预期错误代码
	expectedErrMsg := "192.168.1.1, not in iplist"
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), expectedErrMsg)
}
