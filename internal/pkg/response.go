package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseCode 错误码类型，1000以下为通用码，1000以上为用户自定义码
type ResponseCode int

// Response 响应结构体
type Response struct {
	ErrorCode ResponseCode `json:"errno"`
	ErrorMsg  string       `json:"errmsg"`
	Data      interface{}  `json:"data"`
	TraceID   interface{}  `json:"trace_id"`
	Stack     interface{}  `json:"stack"`
}

// ResponseError 错误响应函数，code 为错误码，err 为错误信息
func ResponseError(c *gin.Context, code ResponseCode, err error) {
	// 获取 TraceID
	traceID, _ := c.Get("TraceID")

	// 设置响应头
	c.Header("Content-Type", "application/json; charset=utf-8")

	// 构造响应体
	c.JSON(http.StatusInternalServerError, Response{
		ErrorCode: code,
		ErrorMsg:  err.Error(),
		TraceID:   traceID,
	})
}

// ResponseSuccess 成功响应函数，data 为响应数据
func ResponseSuccess(c *gin.Context, msg string, data any) {
	traceID, _ := c.Get("TraceID")

	// 设置响应头
	c.Header("Content-Type", "application/json; charset=utf-8")
	if msg == "" {
		msg = "success"
	}
	// 构造响应体
	c.JSON(http.StatusOK, Response{
		ErrorCode: SuccessCode,
		ErrorMsg:  msg,
		Data:      data,
		TraceID:   traceID,
	})
}
